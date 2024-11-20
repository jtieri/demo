package e2e

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"cosmossdk.io/math"
	coinztypes "github.com/jtieri/demo/x/coinz/types"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/testreporter"
	"github.com/strangelove-ventures/interchaintest/v8/testutil"
	"github.com/stretchr/testify/require"
)

// TestUpdateAdmin asserts that initializing the Admin at genesis works correctly,
// as well as updating the admin via MsgUpdateAdmin.
func TestUpdateAdmin(t *testing.T) {
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

	// Query the admin address from the coinz module and assert it equals the expected admin address.
	valNode := chainA.GetNode()
	addr, err := queryAdminAddress(ctx, t, valNode)
	require.NoError(t, err)
	require.Equal(t, defaultAdminAddr, addr, fmt.Sprintf("expected(%s) got(%s)", defaultAdminAddr, addr))

	// Attempt to update the admin address from a non-admin account.
	err = updateAdminAddress(ctx, t, userWallet.KeyName(), userAddress, valNode)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), coinztypes.ErrCannotUpdateAdmin.Error()), "error", err.Error())

	// Query the admin address again and assert that it hasn't changed.
	notUpdatedAddr, err := queryAdminAddress(ctx, t, valNode)
	require.NoError(t, err)
	require.Equal(t, defaultAdminAddr, notUpdatedAddr, fmt.Sprintf("expected(%s) got(%s)", defaultAdminAddr, notUpdatedAddr))

	// Attempt to update the admin address from the admin account.
	err = updateAdminAddress(ctx, t, adminWallet.KeyName(), userAddress, valNode)
	require.NoError(t, err)

	// Wait a few blocks for admin to be updated.
	err = testutil.WaitForBlocks(ctx, 5, chainA)
	require.NoError(t, err)

	// Query the admin address again and assert that it has changed to the expected address.
	updatedAddr, err := queryAdminAddress(ctx, t, valNode)
	require.NoError(t, err)
	require.Equal(t, userAddress, updatedAddr, fmt.Sprintf("expected(%s) got(%s)", userAddress, updatedAddr))

	// Attempt to change the admin address back to the original admin address from the new admin address.
	err = updateAdminAddress(ctx, t, userWallet.KeyName(), defaultAdminAddr, valNode)
	require.NoError(t, err)

	// Assert that the admin address equals the original admin address again.
	addr, err = queryAdminAddress(ctx, t, valNode)
	require.NoError(t, err)
	require.Equal(t, defaultAdminAddr, addr, fmt.Sprintf("expected(%s) got(%s)", defaultAdminAddr, addr))

	// Attempt to update the admin address to an empty string.
	emptyAddress := ""
	err = updateAdminAddress(ctx, t, adminWallet.KeyName(), emptyAddress, valNode)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "empty address string is not allowed"), "error", err.Error())

	// Assert that the admin address was not updated.
	addr, err = queryAdminAddress(ctx, t, valNode)
	require.NoError(t, err)
	require.Equal(t, defaultAdminAddr, addr, fmt.Sprintf("expected(%s) got(%s)", defaultAdminAddr, addr))

	// Attempt to update the admin address to an invalid address.
	invalidAddress := "invalid-address"
	err = updateAdminAddress(ctx, t, adminWallet.KeyName(), invalidAddress, valNode)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "decoding bech32 failed"), "error", err.Error())

	// Assert that the admin address was not updated.
	addr, err = queryAdminAddress(ctx, t, valNode)
	require.NoError(t, err)
	require.Equal(t, defaultAdminAddr, addr, fmt.Sprintf("expected(%s) got(%s)", defaultAdminAddr, addr))

	return
}

// queryAdminAddress is a helper method for executing the admin-address query via demod q admin-address.
// If the query is executed successfully it returns the admin address, otherwise it returns an empty string and an error.
func queryAdminAddress(ctx context.Context, t *testing.T, chain *cosmos.ChainNode) (string, error) {
	t.Helper()

	cmd := []string{
		"coinz", "admin-address",
	}

	stdout, stderr, err := chain.ExecQuery(ctx, cmd...)
	if err != nil {
		return "", fmt.Errorf("query failed: %v, stdout: %s, stderr: %s", err, stdout, stderr)
	}

	// Output of the query looks like: { "address": "cosmos1jtvahy53nvaftwq6r04dmnalv6re2av0xgp705" }
	// So we split at the character " and retrieve only the address string from the output.
	parts := strings.Split(string(stdout), "\"")

	return parts[3], nil
}

// updateAdminAddress is a helper method for executing the update-admin tx via demod tx update-admin.
// If the tx is executed successfully it returns nil, otherwise it returns an error.
func updateAdminAddress(
	ctx context.Context,
	t *testing.T,
	keyName string,
	adminAddr string,
	chain *cosmos.ChainNode,
) error {
	t.Helper()

	cmd := []string{
		"coinz", "update-admin", adminAddr,
	}

	_, err := chain.ExecTx(ctx, keyName, cmd...)
	if err != nil {
		return fmt.Errorf("update admin address failed: %v", err)
	}

	return nil
}
