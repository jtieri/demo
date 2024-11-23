package coinz_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/jtieri/demo/testutil/keeper"
	"github.com/jtieri/demo/testutil/sample"
	coinz "github.com/jtieri/demo/x/coinz/module"
	"github.com/jtieri/demo/x/coinz/types"
	"github.com/stretchr/testify/require"
)

// TestGenesis asserts that the expected genesis values are populated via InitGenesis and are exported properly via ExportGenesis.
func TestGenesis(t *testing.T) {
	adminAddress := sample.AccAddress()

	assetMetadata := types.AssetMetadata{
		Asset: sdk.NewCoin("myToken", math.NewInt(10_000)),
	}

	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		Admin:  &types.Admin{Address: adminAddress},
		Asset:  &assetMetadata,
		Pause:  &types.PauseState{Paused: true},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.CoinzKeeper(t)
	coinz.InitGenesis(ctx, k, genesisState)
	got := coinz.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	//nullify.Fill(&genesisState)
	//nullify.Fill(got)

	require.NotNil(t, got.Admin)
	require.Equal(t, adminAddress, got.Admin.Address)

	require.NotNil(t, got.Asset)
	require.True(t, got.Asset.Asset.IsValid())
	require.Equal(t, assetMetadata.Asset.Denom, got.Asset.Asset.Denom)
	require.True(t, assetMetadata.Asset.Amount.Equal(got.Asset.Asset.Amount))

	require.Equal(t, genesisState.Pause, got.Pause)
	// this line is used by starport scaffolding # genesis/test/assert
}
