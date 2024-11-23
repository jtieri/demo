package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"cosmossdk.io/math"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/testreporter"
	"github.com/stretchr/testify/require"
)

// TestPause asserts that initializing the PauseState at genesis works correctly,
// as well as updating and querying the PauseState.
func TestPause(t *testing.T) {
	var (
		ctx  = context.Background()
		rep  = testreporter.NewNopReporter()
		eRep = rep.RelayerExecReporter(t)
	)

	chains, err := startDemoChain(ctx, t, eRep, []*interchaintest.ChainSpec{&chainASpec})
	require.NoError(t, err)
	require.Len(t, chains, 1)

	chain := chains[0]
	chainA := chain.(*cosmos.CosmosChain)

	// Create the Admin on chain.
	adminWallet, err := interchaintest.GetAndFundTestUserWithMnemonic(ctx, "admin", defaultAdminSeed, math.NewInt(10_000), chainA)
	require.NoError(t, err)
	require.Equal(t, defaultAdminAddr, adminWallet.FormattedAddress(), fmt.Sprintf("expected(%s) got(%s)", defaultAdminAddr, adminWallet.FormattedAddress()))

	// Create a user account on chain.
	wallets := interchaintest.GetAndFundTestUsers(t, ctx, "user", math.NewInt(10_000), chainA)
	userWallet := wallets[0]
	userAddress := userWallet.FormattedAddress()
	require.NotEqual(t, "", userAddress)

	// Query the initial PauseState.
	valNode := chainA.GetNode()

	pauseState, err := queryPauseState(ctx, t, valNode)
	require.NoError(t, err)
	require.True(t, pauseState)

	// Attempt to update the PauseState from a user account.
	err = updatePauseState(ctx, t, userWallet.KeyName(), false, valNode)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "cannot update pause state"))

	// Assert that the PauseState was not updated.
	pauseState, err = queryPauseState(ctx, t, valNode)
	require.NoError(t, err)
	require.True(t, pauseState)

	// Attempt to update the PauseState from the admin account.
	err = updatePauseState(ctx, t, adminWallet.KeyName(), false, valNode)
	require.NoError(t, err)

	// Assert that the PauseState was updated.
	pauseState, err = queryPauseState(ctx, t, valNode)
	require.NoError(t, err)
	require.False(t, pauseState)

	// Attempt to update the PauseState with a redundant value.
	err = updatePauseState(ctx, t, adminWallet.KeyName(), false, valNode)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "pause state is already set to"))
}

// updatePauseState is a helper method for executing the update pause state tx via demod tx update-pause-state.
// If the tx is executed successfully it returns nil, otherwise it returns an error.
func updatePauseState(ctx context.Context, t *testing.T, keyName string, pause bool, chain *cosmos.ChainNode) error {
	t.Helper()

	updateState := "false"

	if pause {
		updateState = "true"
	}

	cmd := []string{
		"coinz", "update-pause-state", updateState,
	}

	_, err := chain.ExecTx(ctx, keyName, cmd...)
	if err != nil {
		return err
	}

	return nil
}

// PauseState is used to unmarshal the json response from the show pause state query.
type PausedState struct {
	Paused bool `json:"paused"`
}

// queryPauseState queries the current PauseState of the coinz module.
func queryPauseState(ctx context.Context, t *testing.T, chain *cosmos.ChainNode) (bool, error) {
	t.Helper()

	cmd := []string{
		"coinz", "show-pause-state",
	}

	stdout, _, err := chain.ExecQuery(ctx, cmd...)
	if err != nil {
		return false, fmt.Errorf("failed to query pause state: %v", err)
	}

	ps := &PausedState{}
	err = json.Unmarshal(stdout, ps)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal pause state: %v", err)
	}

	return ps.Paused, nil
}
