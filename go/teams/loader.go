package teams

import (
	"encoding/json"
	"errors"
	"fmt"

	"golang.org/x/net/context"

	"github.com/keybase/client/go/libkb"
	"github.com/keybase/client/go/protocol/keybase1"
)

// TODO: big problem: If you are added to a subteam then you are allowed to see previously hidden
// parent team sigchain links (the subteam links).

type LoadTeamFreshness int

const (
	LoadTeamFreshnessRANCID LoadTeamFreshness = 0
	LoadTeamFreshnessAGED   LoadTeamFreshness = 1
	LoadTeamFreshnessFRESH  LoadTeamFreshness = 2
)

// Load a Team from the TeamLoader.
// Can be called from inside the teams.
func Load(ctx context.Context, g *libkb.GlobalContext, lArg libkb.LoadTeamArg) (*Team, error) {
	teamData, err := g.GetTeamLoader().Load(ctx, lArg)
	if err != nil {
		return nil, err
	}
	return &Team{
		Contextified: libkb.NewContextified(g),
		TeamData:     teamData,
	}, nil
}

// Loader of keybase1.TeamData objects. Handles caching.
// Because there is one of this global object and it is attached to G,
// its Load interface must return a keybase1.TeamData not a teams.Team.
// To load a teams.Team use the package-level function Load.
// Threadsafe.
type TeamLoader struct {
	libkb.Contextified
	storage *Storage
	// Single-flight locks pre team ID.
	locktab libkb.LockTable
}

func NewTeamLoader(g *libkb.GlobalContext, storage *Storage) *TeamLoader {
	return &TeamLoader{
		Contextified: libkb.NewContextified(g),
		storage:      storage,
	}
}

// NewTeamLoaderAndInstall creates a new loader and installs it into G.
func NewTeamLoaderAndInstall(g *libkb.GlobalContext) *TeamLoader {
	st := NewStorage(g)
	l := NewTeamLoader(g, st)
	g.SetTeamLoader(l)
	return l
}

func (l *TeamLoader) Load(ctx context.Context, lArg libkb.LoadTeamArg) (res *keybase1.TeamData, err error) {
	me, err := l.getMe(ctx)
	if err != nil {
		return nil, err
	}
	return l.load(ctx, me, lArg)
}

type loadInfoT struct {
	hitCache          bool
	loadedFromServer  bool
	loadedFromScratch bool
}

func (l *TeamLoader) getMe(ctx context.Context) (res keybase1.UserVersion, err error) {
	return loadUserVersionByUID(ctx, l.G(), l.G().Env.GetUID())
}

func (l *TeamLoader) load(ctx context.Context, me keybase1.UserVersion, lArg libkb.LoadTeamArg) (res *keybase1.TeamData, err error) {
	res, info, err := l.loadInner(ctx, me, lArg)
	if err != nil {
		l.G().Log.CDebugf(ctx, "TeamLoader#Load -> err:%v info:%+v", err, info)
	} else {
		l.G().Log.CDebugf(ctx, "TeamLoader#Load -> info:%+v", info)
	}
	return res, err
}

func (l *TeamLoader) loadInner(ctx context.Context, me keybase1.UserVersion, lArg libkb.LoadTeamArg) (res *keybase1.TeamData, info loadInfoT, err error) {
	// GIANT TODO: check for role change
	// GIANT TODO: load recursively to load subteams

	err = lArg.Check()
	if err != nil {
		return nil, info, err
	}

	// extract team id
	teamID := lArg.ID
	if len(lArg.ID) == 0 {
		teamName, err := TeamNameFromString(lArg.Name)
		if err != nil {
			return nil, info, err
		}
		if teamName.IsSubTeam() {
			// To resolve a subteam name to a team ID do this:
			// Split the name and find the root team name.
			// Load the root team.
			// Load the L2 subteam by id gleaned from the root team chain.
			// Continue until the LN subteam is reached.
			return nil, info, fmt.Errorf("TODO: support loading subteams by name")
		}
		teamID = teamName.ToTeamID()
	}
	if len(teamID) == 0 {
		return nil, info, fmt.Errorf("team loader fault: empty team id")
	}

	lock := l.locktab.AcquireOnName(ctx, l.G(), teamID.String())
	defer lock.Release(ctx)

	// pull from cache
	var cacheResult *keybase1.TeamData
	if !lArg.ForceFullReload {
		cacheResult = l.storage.Get(ctx, teamID)
	}
	info.hitCache = (cacheResult != nil)

	// calculate freshness
	freshness := LoadTeamFreshnessRANCID
	if cacheResult != nil {
		// TODO the cache is not always fresh
		freshness = LoadTeamFreshnessFRESH

		if lArg.KeyGeneration != 0 {
			cKey, err := TeamSigChainState{inner: cacheResult.Chain}.GetLatestPerTeamKey()
			if (err != nil) || cKey.Gen < lArg.KeyGeneration {
				freshness = LoadTeamFreshnessRANCID
			}
		}
	}

	// TODO StaleOK is wrong around here
	cacheTooOld := (freshness != LoadTeamFreshnessFRESH)

	// pull from server
	if cacheResult == nil || cacheTooOld || lArg.ForceSync || lArg.ForceFullReload {
		if lArg.NoNetwork {
			return nil, info, fmt.Errorf("cannot load from server with no-network set")
		}
		if cacheResult == nil || lArg.ForceFullReload {
			res, err = l.loadFromServerFromScratch(ctx, me, teamID)
			if err != nil {
				return nil, info, err
			}
			info.loadedFromScratch = true
		} else {
			res, err = l.loadFromServerWithCached(ctx, me, cacheResult)
			if err != nil {
				return nil, info, err
			}
		}
		info.loadedFromServer = true
	}

	// TODO check freshness and load increment from server
	// TODO check key generation

	if res == nil {
		return nil, info, fmt.Errorf("team loader fault: no result")
	}
	resID := TeamSigChainState{inner: res.Chain}.GetID()
	if !resID.Equal(teamID) {
		return nil, info, fmt.Errorf("retrieved team with wrong id %v != %v", resID, teamID)
	}

	// Put the result even if it hasn't changed to bump the freshness
	l.storage.Put(ctx, *res)

	if len(lArg.Name) > 0 {
		retrievedName := TeamSigChainState{inner: res.Chain}.GetName()
		if lArg.Name != retrievedName {
			return nil, info, fmt.Errorf("retrieved team with wrong name %v != %v", retrievedName, lArg.Name)
		}
	}

	return res, info, nil
}

// Load a team from the server with no cached data.
func (l *TeamLoader) loadFromServerFromScratch(ctx context.Context, me keybase1.UserVersion, teamID keybase1.TeamID) (*keybase1.TeamData, error) {
	sArg := libkb.NewRetryAPIArg("team/get")
	sArg.NetContext = ctx
	sArg.SessionType = libkb.APISessionTypeREQUIRED
	sArg.Args = libkb.HTTPArgs{
		// "name": libkb.S{Val: string(lArg.Name)},
		"id": libkb.S{Val: teamID.String()},
		// TODO used cached last seqno 0 (in a non from-scratch function)
		"low": libkb.I{Val: 0},
	}
	var rt rawTeam
	if err := l.G().API.GetDecode(sArg, &rt); err != nil {
		l.G().Log.CDebugf(ctx, "TeamLoader fetch error: %v", err)
		return nil, err
	}

	links, err := l.parseChainLinks(ctx, &rt)
	if err != nil {
		return nil, err
	}

	player, err := l.newPlayer(ctx, me, links)
	if err != nil {
		return nil, err
	}

	state, err := player.GetState()
	if err != nil {
		return nil, err
	}

	// TODO (non-critical) validate reader key masks

	res := keybase1.TeamData{
		Chain:           state.inner,
		PerTeamKeySeeds: nil,
		ReaderKeyMasks:  rt.ReaderKeyMasks,
	}

	seed, err := l.openBox(ctx, rt.Box, state)
	if err != nil {
		return nil, err
	}
	res.PerTeamKeySeeds = append(res.PerTeamKeySeeds, *seed)

	// TODO receive prevs (and sort seeds list)

	return &res, nil
}

// Load a team from the server with cached data
func (l *TeamLoader) loadFromServerWithCached(ctx context.Context, me keybase1.UserVersion, cached *keybase1.TeamData) (*keybase1.TeamData, error) {
	if cached == nil {
		return nil, fmt.Errorf("team loader fault: no cached object")
	}
	cachedChain := TeamSigChainState{inner: cached.Chain}

	sArg := libkb.NewRetryAPIArg("team/get")
	sArg.NetContext = ctx
	sArg.SessionType = libkb.APISessionTypeREQUIRED
	sArg.Args = libkb.HTTPArgs{
		"id":  libkb.S{Val: cachedChain.GetID().String()},
		"low": libkb.I{Val: int(cachedChain.GetLatestSeqno())},
	}
	var rt rawTeam
	if err := l.G().API.GetDecode(sArg, &rt); err != nil {
		l.G().Log.CDebugf(ctx, "TeamLoader fetch error: %v", err)
		return nil, err
	}

	links, err := l.parseChainLinks(ctx, &rt)
	if err != nil {
		return nil, err
	}

	player, err := l.newPlayerWithState(ctx, me, cachedChain, links)
	if err != nil {
		return nil, err
	}

	state, err := player.GetState()
	if err != nil {
		return nil, err
	}

	// TODO receive prevs (and sort seeds list)

	res, err := l.merge(ctx, cached, state, rt.ReaderKeyMasks, rt.Box)
	if err != nil {
		return nil, err
	}

	nNewLinks := state.GetLatestSeqno() - cachedChain.GetLatestSeqno()
	if nNewLinks > 0 {
		l.G().Log.CDebugf(ctx, "TeamLoader loaded %v new links: seqno %v -> %v",
			nNewLinks,
			cachedChain.GetLatestSeqno(), state.GetLatestSeqno())
	}

	return res, nil
}

func (l *TeamLoader) merge(ctx context.Context, cached *keybase1.TeamData, newState TeamSigChainState, serverReaderKeyMasks []keybase1.ReaderKeyMask, box TeamBox) (*keybase1.TeamData, error) {

	// TODO (non-critical) validate reader key masks
	// TODO dedup and maybe sort reader key masks

	seed, err := l.openBox(ctx, box, newState)
	if err != nil {
		return nil, err
	}

	return &keybase1.TeamData{
		Chain:           newState.inner,
		PerTeamKeySeeds: append(cached.PerTeamKeySeeds, *seed),
		ReaderKeyMasks:  append(cached.ReaderKeyMasks, serverReaderKeyMasks...),
	}, nil
}

func (l *TeamLoader) parseChainLinks(ctx context.Context, rawTeam *rawTeam) ([]SCChainLink, error) {
	var links []SCChainLink
	for _, raw := range rawTeam.Chain {
		link, err := ParseTeamChainLink(string(raw))
		if err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	return links, nil
}

func (l *TeamLoader) newPlayer(ctx context.Context, me keybase1.UserVersion, links []SCChainLink) (*TeamSigChainPlayer, error) {
	// TODO determine whether this player is for a subteam
	isSubteam := false

	f := newFinder(l.G())
	player := NewTeamSigChainPlayer(l.G(), f, me, isSubteam)
	if len(links) > 0 {
		if err := player.AddChainLinks(ctx, links); err != nil {
			return nil, err
		}
	}
	return player, nil
}

func (l *TeamLoader) newPlayerWithState(ctx context.Context, me keybase1.UserVersion, state TeamSigChainState, links []SCChainLink) (*TeamSigChainPlayer, error) {
	f := newFinder(l.G())
	player := NewTeamSigChainPlayerWithState(l.G(), f, me, state)
	if len(links) > 0 {
		if err := player.AddChainLinks(ctx, links); err != nil {
			return nil, err
		}
	}
	return player, nil
}

func (l *TeamLoader) openBox(ctx context.Context, box TeamBox, chain TeamSigChainState) (*keybase1.PerTeamKeySeedItem, error) {
	userEncKey, err := l.perUserEncryptionKeyForBox(ctx, box)
	if err != nil {
		return nil, err
	}

	secret, err := box.Open(userEncKey)
	if err != nil {
		return nil, err
	}

	keyManager, err := NewTeamKeyManagerWithSecret(l.G(), secret, box.Generation)
	if err != nil {
		return nil, err
	}

	signingKey, err := keyManager.SigningKey()
	if err != nil {
		return nil, err
	}
	encryptionKey, err := keyManager.EncryptionKey()
	if err != nil {
		return nil, err
	}

	teamKey, err := chain.GetPerTeamKeyAtGeneration(int(box.Generation))
	if err != nil {
		return nil, err
	}

	if !teamKey.SigKID.SecureEqual(signingKey.GetKID()) {
		return nil, errors.New("derived signing key did not match key in team chain")
	}

	if !teamKey.EncKID.SecureEqual(encryptionKey.GetKID()) {
		return nil, errors.New("derived encryption key did not match key in team chain")
	}

	// TODO: check that t.Box.SenderKID is a known device DH key for the
	// user that signed the link.
	// See CORE-5399

	seed, err := libkb.MakeByte32Soft(secret)
	if err != nil {
		return nil, fmt.Errorf("invalid seed: %v", err)
	}

	record := keybase1.PerTeamKeySeedItem{
		Seed:       seed,
		Generation: int(box.Generation),
		Seqno:      teamKey.Seqno,
	}

	return &record, nil
}

func (l *TeamLoader) perUserEncryptionKeyForBox(ctx context.Context, box TeamBox) (*libkb.NaclDHKeyPair, error) {
	kr, err := l.G().GetPerUserKeyring()
	if err != nil {
		return nil, err
	}
	// XXX this seems to be necessary:
	if err := kr.Sync(ctx); err != nil {
		return nil, err
	}
	encKey, err := kr.GetEncryptionKeyBySeqno(ctx, box.PerUserKeySeqno)
	if err != nil {
		return nil, err
	}

	return encKey, nil
}

type rawTeam struct {
	ID             keybase1.TeamID          `json:"id"`
	Name           keybase1.TeamNameParts   `json:"name"`
	Status         libkb.AppStatus          `json:"status"`
	Chain          []json.RawMessage        `json:"chain"`
	Box            TeamBox                  `json:"box"`
	ReaderKeyMasks []keybase1.ReaderKeyMask `json:"reader_key_masks"`
}

func (r *rawTeam) GetAppStatus() *libkb.AppStatus {
	return &r.Status
}
