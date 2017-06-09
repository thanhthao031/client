package teams

import (
	"context"
	"testing"

	"github.com/keybase/client/go/kbtest"
	"github.com/keybase/client/go/libkb"
	"github.com/keybase/client/go/protocol/keybase1"
	"github.com/stretchr/testify/require"
)

// Create n TestContexts with logged in users
func setupNTests(t *testing.T, n int) ([]*kbtest.FakeUser, []*libkb.TestContext) {
	require.True(t, n > 0, "must create at least 1 tc")
	var fus []*kbtest.FakeUser
	var tcs []*libkb.TestContext
	for i := 0; i < n; i++ {
		tc := SetupTest(t, "team", 1)
		tcs = append(tcs, &tc)
		fu, err := kbtest.CreateAndSignupFakeUser("team", tc.G)
		require.NoError(t, err)
		fus = append(fus, fu)
	}
	return fus, tcs
}

// Test that tests with two Gs don't share the same database.
func TestTestsDontAlias(t *testing.T) {
	tc1 := SetupTest(t, "team", 1)
	_, err := kbtest.CreateAndSignupFakeUser("team", tc1.G)
	require.NoError(t, err)

	teamName := createTeam(tc1)

	_, err = Load(context.TODO(), tc1.G, libkb.LoadTeamArg{
		Name:      teamName,
		ForceSync: true,
	})
	require.NoError(t, err)

	tc2 := SetupTest(t, "team", 1)
	_, err = kbtest.CreateAndSignupFakeUser("team", tc2.G)

	_, err = Load(context.TODO(), tc2.G, libkb.LoadTeamArg{
		Name:      teamName,
		NoNetwork: true,
	})
	require.Error(t, err, "load from cache should fail")
	require.Equal(t, err.Error(), "cannot load from server with no-network set")
}

func TestLoadAfterPromotion(t *testing.T) {

	// TODO unskip
	// TODO unskip
	// TODO unskip
	// TODO unskip
	// TODO unskip
	// TODO unskip
	// TODO unskip
	// TODO unskip
	// TODO unskip
	// TODO unskip
	// TODO unskip
	// TODO unskip
	// TODO unskip
	t.Skip()

	fus, tcs := setupNTests(t, 2)

	t.Logf("u0 creates a team")
	teamName := createTeam(*tcs[0])

	t.Logf("u0 makes u1 a reader")
	err := SetRoleReader(context.TODO(), tcs[0].G, teamName, fus[1].Username)
	require.NoError(t, err)

	t.Logf("u0 adds a subteam")
	err = CreateSubteam(context.TODO(), tcs[0].G, "smilers", TeamName(teamName))
	require.NoError(t, err)

	t.Logf("u1 loads the team")
	u1uv, err := loadUserVersionByUID(context.TODO(), tcs[1].G, tcs[1].G.Env.GetUID())
	require.NoError(t, err)
	team, err := Load(context.TODO(), tcs[1].G, libkb.LoadTeamArg{
		Name:      teamName,
		ForceSync: true,
	})
	require.NoError(t, err)
	role, err := team.GetSigChainState().GetUserRole(u1uv)
	require.NoError(t, err)
	require.Equal(t, keybase1.TeamRole_READER, role)

	t.Logf("u0 makes u1 an admin")
	err = SetRoleAdmin(context.TODO(), tcs[0].G, teamName, fus[1].Username)
	require.NoError(t, err)

	t.Logf("u1 loads the team again")
	// u1 must discard their cache because now they have the potential to see
	// v2 links that might have been hidden before (subteam links).
	// This design may be totally broken.
	team, err = Load(context.TODO(), tcs[1].G, libkb.LoadTeamArg{
		Name:      teamName,
		ForceSync: true,
	})
	require.NoError(t, err)
	role, err = team.GetSigChainState().GetUserRole(u1uv)
	require.NoError(t, err)
	require.Equal(t, keybase1.TeamRole_ADMIN, role)
}
