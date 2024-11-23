package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/testreporter"
	"github.com/stretchr/testify/require"
)

// TestMint asserts that initializing the Asset at genesis works correctly,
// as well as the minting of new tokens via the coinz module.
func TestMint(t *testing.T) {
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

	// Attempt to mint assets from the admin account to the user account.
	// This should fail because the module is paused.
	amount := int64(10_000)
	mintAmount := sdk.NewCoin(coinzDenom, math.NewInt(amount))

	err = mintAssets(ctx, t, userWallet.FormattedAddress(), mintAmount, adminWallet.KeyName(), valNode)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "coinz module is currently paused"))

	// Assert that the assets were not minted.
	totalSupply, err = queryTotalSupply(ctx, t, coinzDenom, valNode)
	require.NoError(t, err)
	require.Equal(t, initialSupply, totalSupply)

	// Assert that the user did not receive any funds.
	balance, err := chainA.GetBalance(ctx, userWallet.FormattedAddress(), coinzDenom)
	require.NoError(t, err)
	require.True(t, balance.Equal(math.ZeroInt()))

	// Attempt to update the PauseState from the admin account.
	err = updatePauseState(ctx, t, adminWallet.KeyName(), false, valNode)
	require.NoError(t, err)

	// Assert that the PauseState was updated.
	pauseState, err := queryPauseState(ctx, t, valNode)
	require.NoError(t, err)
	require.False(t, pauseState)

	// Attempt to mint assets from the admin account to the user account.
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

	// Attempt to mint assets from a user account.
	err = mintAssets(ctx, t, userWallet.FormattedAddress(), mintAmount, userWallet.KeyName(), valNode)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "failed to mint assets"))

	// Assert that the assets were not minted and total supply did not change.
	totalSupply, err = queryTotalSupply(ctx, t, coinzDenom, valNode)
	require.NoError(t, err)
	require.Equal(t, mintAmount.Amount.Int64(), int64(totalSupply))

	// Attempt to mint assets for the wrong denom.
	invalidAmount := sdk.NewCoin("wrongDenom", math.NewInt(10_000))

	err = mintAssets(ctx, t, userWallet.FormattedAddress(), invalidAmount, adminWallet.KeyName(), valNode)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "failed to mint assets"))

	// Assert that the assets were not minted.
	totalSupply, err = queryTotalSupply(ctx, t, invalidAmount.Denom, valNode)
	require.NoError(t, err)
	require.Equal(t, 0, totalSupply)

	// Attempt to mint assets with a zero value for the amount.
	invalidAmount = sdk.NewCoin(coinzDenom, math.NewInt(0))

	err = mintAssets(ctx, t, userWallet.FormattedAddress(), invalidAmount, adminWallet.KeyName(), valNode)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "failed to mint assets"))

	// Assert that the assets were not minted.
	totalSupply, err = queryTotalSupply(ctx, t, invalidAmount.Denom, valNode)
	require.NoError(t, err)
	require.Equal(t, mintAmount.Amount.Int64(), int64(totalSupply))

	// Attempt to mint assets with a negative value for the amount.
	invalidAmount = sdk.Coin{
		Denom:  coinzDenom,
		Amount: math.NewInt(-1234),
	}

	err = mintAssets(ctx, t, userWallet.FormattedAddress(), invalidAmount, adminWallet.KeyName(), valNode)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "failed to mint assets"))

	// Assert that the assets were not minted.
	totalSupply, err = queryTotalSupply(ctx, t, invalidAmount.Denom, valNode)
	require.NoError(t, err)
	require.Equal(t, mintAmount.Amount.Int64(), int64(totalSupply))
}

// mintAssets is a helper method for executing the mint tx via demod tx mint.
// If the tx is executed successfully it returns nil, otherwise it returns an error.
func mintAssets(
	ctx context.Context,
	t *testing.T,
	toAddress string,
	amount sdk.Coin,
	keyName string,
	chain *cosmos.ChainNode,
) error {
	t.Helper()

	cmd := []string{
		"coinz", "mint", toAddress, amount.String(),
	}

	_, err := chain.ExecTx(ctx, keyName, cmd...)
	if err != nil {
		return fmt.Errorf("failed to mint assets: %v", err)
	}

	return nil
}

// Coin is used to unmarshal the json response from the total supply query.
type Coin struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

// Response is used to unmarshal the json response from the total supply query.
type Response struct {
	Amount Coin `json:"amount"`
}

// queryTotalSupply queries the total supply of a specified denom via the bank module.
func queryTotalSupply(
	ctx context.Context,
	t *testing.T,
	denom string,
	chain *cosmos.ChainNode,
) (int, error) {
	t.Helper()

	cmd := []string{
		"bank", "total-supply-of", denom,
	}

	stdout, _, err := chain.ExecQuery(ctx, cmd...)
	if err != nil {
		return 0, fmt.Errorf("failed to query total supply of %s: %v", denom, err)
	}

	var resp Response
	err = json.Unmarshal(stdout, &resp)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal total supply of %s: %v", denom, err)
	}

	return strconv.Atoi(resp.Amount.Amount)
}
