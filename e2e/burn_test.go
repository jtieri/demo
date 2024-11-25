package e2e

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/testreporter"
	"github.com/stretchr/testify/require"
)

// TestBurn asserts that the behavior around burning assets via the coinz module works properly.
func TestBurn(t *testing.T) {
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

	// Query the initial total supply and assert it equals the default initial total supply.
	valNode := chainA.GetNode()
	totalSupply, err := queryTotalSupply(ctx, t, coinzDenom, valNode)
	require.NoError(t, err)
	require.Equal(t, initialSupply, totalSupply)

	// Attempt to burn assets from the user account. This should fail because the module is paused.
	userInitBal := math.NewInt(10_000)

	burnAmount := sdk.NewCoin(coinzDenom, userInitBal)
	err = burnAssets(ctx, t, burnAmount, userWallet.KeyName(), valNode)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "coinz module is currently paused"))

	// Assert that the assets were not burned and the user still has the funds.
	balance, err := chainA.GetBalance(ctx, userWallet.FormattedAddress(), coinzDenom)
	require.NoError(t, err)
	require.True(t, balance.Equal(math.ZeroInt()))

	// Assert that the total supply did not change.
	totalSupply, err = queryTotalSupply(ctx, t, coinzDenom, valNode)
	require.NoError(t, err)
	require.Equal(t, initialSupply, totalSupply)

	// Attempt to update the PauseState from the admin account.
	err = updatePauseState(ctx, t, adminWallet.KeyName(), false, valNode)
	require.NoError(t, err)

	// Assert that the PauseState was updated.
	pauseState, err := queryPauseState(ctx, t, valNode)
	require.NoError(t, err)
	require.False(t, pauseState)

	// Attempt to mint assets from the admin account to the user account.
	mintAmount := sdk.NewCoin(coinzDenom, userInitBal)

	err = mintAssets(ctx, t, userWallet.FormattedAddress(), mintAmount, adminWallet.KeyName(), valNode)
	require.NoError(t, err)

	// Assert that the assets were minted and total supply increased.
	totalSupply, err = queryTotalSupply(ctx, t, coinzDenom, valNode)
	require.NoError(t, err)
	require.Equal(t, mintAmount.Amount.Int64(), int64(totalSupply))

	// Assert that the user received the minted funds.
	balance, err = chainA.GetBalance(ctx, userWallet.FormattedAddress(), coinzDenom)
	require.NoError(t, err)
	require.True(t, balance.Equal(mintAmount.Amount))

	// Attempt to burn assets from the user account.
	err = burnAssets(ctx, t, burnAmount, userWallet.KeyName(), valNode)
	require.NoError(t, err)

	// Assert that the assets were burned.
	balance, err = chainA.GetBalance(ctx, userWallet.FormattedAddress(), coinzDenom)
	require.NoError(t, err)
	require.True(t, balance.Equal(userInitBal.Sub(burnAmount.Amount)))

	// Attempt to burn assets from the user account when the user balance is < the burn amount.
	updatedSupply, err := queryTotalSupply(ctx, t, coinzDenom, valNode)
	require.NoError(t, err)
	require.True(t, math.NewInt(int64(totalSupply)).Sub(burnAmount.Amount).Equal(math.NewInt(int64(updatedSupply))))
}

// burnAssets is a helper method for executing the burn tx via demod tx burn.
// If the tx is successful it returns nil, otherwise it returns an error.
func burnAssets(
	ctx context.Context,
	t *testing.T,
	amount sdk.Coin,
	keyName string,
	chain *cosmos.ChainNode,
) error {
	t.Helper()

	cmd := []string{
		"coinz", "burn", amount.String(),
	}

	_, err := chain.ExecTx(ctx, keyName, cmd...)
	if err != nil {
		return fmt.Errorf("failed to burn assets: %v", err)
	}

	return nil
}
