package teams

import (
	"errors"
	"fmt"

	"golang.org/x/net/context"

	"github.com/keybase/client/go/libkb"
	"github.com/keybase/client/go/protocol/keybase1"
)

func membersUIDsToUsernames(ctx context.Context, g *libkb.GlobalContext, m keybase1.TeamMembers) (keybase1.TeamMembersUsernames, error) {
	var ret keybase1.TeamMembersUsernames
	var err error
	ret.Owners, err = userVersionsToUsernames(ctx, g, m.Owners)
	if err != nil {
		return ret, err
	}
	ret.Admins, err = userVersionsToUsernames(ctx, g, m.Admins)
	if err != nil {
		return ret, err
	}
	ret.Writers, err = userVersionsToUsernames(ctx, g, m.Writers)
	if err != nil {
		return ret, err
	}
	ret.Readers, err = userVersionsToUsernames(ctx, g, m.Readers)
	if err != nil {
		return ret, err
	}
	return ret, nil
}

func Members(ctx context.Context, g *libkb.GlobalContext, name string) (keybase1.TeamMembersUsernames, error) {
	t, err := Get(ctx, g, name)
	if err != nil {
		return keybase1.TeamMembersUsernames{}, err
	}
	members, err := t.Members()
	if err != nil {
		return keybase1.TeamMembersUsernames{}, err
	}
	return membersUIDsToUsernames(ctx, g, members)
}

func usernameToUID(g *libkb.GlobalContext, username string) (uid keybase1.UID, err error) {
	rres := g.Resolver.Resolve(username)
	if err = rres.GetError(); err != nil {
		return uid, err
	}
	return rres.GetUID(), nil
}

func uidToUsername(ctx context.Context, g *libkb.GlobalContext, uid keybase1.UID) (libkb.NormalizedUsername, error) {
	return g.GetUPAKLoader().LookupUsername(ctx, uid)
}

func userVersionsToUsernames(ctx context.Context, g *libkb.GlobalContext, uvs []keybase1.UserVersion) (ret []string, err error) {
	for _, uv := range uvs {
		un, err := uidToUsername(ctx, g, uv.Uid)
		if err != nil {
			return nil, err
		}
		ret = append(ret, string(un))
	}
	return ret, nil
}

func SetRoleOwner(ctx context.Context, g *libkb.GlobalContext, teamname, username string) error {
	uid, err := usernameToUID(g, username)
	if err != nil {
		return err
	}
	return ChangeRoles(ctx, g, teamname, keybase1.TeamChangeReq{Owners: []keybase1.UID{uid}})
}

func SetRoleAdmin(ctx context.Context, g *libkb.GlobalContext, teamname, username string) error {
	uid, err := usernameToUID(g, username)
	if err != nil {
		return err
	}
	return ChangeRoles(ctx, g, teamname, keybase1.TeamChangeReq{Admins: []keybase1.UID{uid}})
}

func SetRoleWriter(ctx context.Context, g *libkb.GlobalContext, teamname, username string) error {
	uid, err := usernameToUID(g, username)
	if err != nil {
		return err
	}
	return ChangeRoles(ctx, g, teamname, keybase1.TeamChangeReq{Writers: []keybase1.UID{uid}})
}

func SetRoleReader(ctx context.Context, g *libkb.GlobalContext, teamname, username string) error {
	uid, err := usernameToUID(g, username)
	if err != nil {
		return err
	}
	return ChangeRoles(ctx, g, teamname, keybase1.TeamChangeReq{Readers: []keybase1.UID{uid}})
}

func AddMember(ctx context.Context, g *libkb.GlobalContext, teamname, username string, role keybase1.TeamRole) error {
	t, err := Get(ctx, g, teamname)
	if err != nil {
		return err
	}
	uid, err := usernameToUID(g, username)
	if err != nil {
		return err
	}
	if t.IsMember(ctx, uid) {
		return fmt.Errorf("user %q is already a member of team %q", username, teamname)
	}
	req, err := reqFromRole(uid, role)
	if err != nil {
		return err
	}

	return t.ChangeMembership(ctx, req)
}

func EditMember(ctx context.Context, g *libkb.GlobalContext, teamname, username string, role keybase1.TeamRole) error {
	t, err := Get(ctx, g, teamname)
	if err != nil {
		return err
	}
	uid, err := usernameToUID(g, username)
	if err != nil {
		return err
	}
	if !t.IsMember(ctx, uid) {
		return fmt.Errorf("user %q is not a member of team %q", username, teamname)
	}
	existingRole, err := t.MemberRole(ctx, uid)
	if err != nil {
		return err
	}
	if existingRole == role {
		return fmt.Errorf("user %q in team %q already has the role %s", username, teamname, role)
	}

	req, err := reqFromRole(uid, role)
	if err != nil {
		return err
	}

	return t.ChangeMembership(ctx, req)
}

func MemberRole(ctx context.Context, g *libkb.GlobalContext, teamname, username string) (keybase1.TeamRole, error) {
	t, err := Get(ctx, g, teamname)
	if err != nil {
		return keybase1.TeamRole_NONE, err
	}
	uid, err := usernameToUID(g, username)
	if err != nil {
		return keybase1.TeamRole_NONE, err
	}
	return t.MemberRole(ctx, uid)
}

func RemoveMember(ctx context.Context, g *libkb.GlobalContext, teamname, username string) error {
	t, err := Get(ctx, g, teamname)
	if err != nil {
		return err
	}
	uid, err := usernameToUID(g, username)
	if err != nil {
		return err
	}
	if !t.IsMember(ctx, uid) {
		return libkb.NotFoundError{Msg: fmt.Sprintf("user %q is not a member of team %q", username, teamname)}
	}
	req := keybase1.TeamChangeReq{None: []keybase1.UID{uid}}
	return t.ChangeMembership(ctx, req)
}

func ChangeRoles(ctx context.Context, g *libkb.GlobalContext, teamname string, req keybase1.TeamChangeReq) error {
	t, err := Get(ctx, g, teamname)
	if err != nil {
		return err
	}
	return t.ChangeMembership(ctx, req)
}

func loadUserVersionByUsername(ctx context.Context, g *libkb.GlobalContext, username string) (keybase1.UserVersion, error) {
	res := g.Resolver.ResolveWithBody(username)
	if res.GetError() != nil {
		return keybase1.UserVersion{}, res.GetError()
	}
	return loadUserVersionByUID(ctx, g, res.GetUID())
}

func loadUserVersionByUID(ctx context.Context, g *libkb.GlobalContext, uid keybase1.UID) (keybase1.UserVersion, error) {
	arg := libkb.NewLoadUserByUIDArg(ctx, g, uid)
	upak, _, err := g.GetUPAKLoader().Load(arg)
	if err != nil {
		return keybase1.UserVersion{}, err
	}

	return NewUserVersion(upak.Base.Uid, upak.Base.EldestSeqno), nil
}

func reqFromRole(uid keybase1.UID, role keybase1.TeamRole) (keybase1.TeamChangeReq, error) {

	var req keybase1.TeamChangeReq
	list := []keybase1.UID{uid}
	switch role {
	case keybase1.TeamRole_OWNER:
		req.Owners = list
	case keybase1.TeamRole_ADMIN:
		req.Admins = list
	case keybase1.TeamRole_WRITER:
		req.Writers = list
	case keybase1.TeamRole_READER:
		req.Readers = list
	default:
		return keybase1.TeamChangeReq{}, errors.New("invalid team role")
	}

	return req, nil
}
