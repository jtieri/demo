package coinz_test

import (
	"testing"

	keepertest "github.com/jtieri/demo/testutil/keeper"
	"github.com/jtieri/demo/testutil/nullify"
	"github.com/jtieri/demo/testutil/sample"
	coinz "github.com/jtieri/demo/x/coinz/module"
	"github.com/jtieri/demo/x/coinz/types"
	"github.com/stretchr/testify/require"
)

// TestGenesis asserts that the expected genesis values are populated via InitGenesis and are exported properly via ExportGenesis.
func TestGenesis(t *testing.T) {
	adminAddress := sample.AccAddress()

	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		Admin:  &types.Admin{Address: adminAddress},

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.CoinzKeeper(t)
	coinz.InitGenesis(ctx, k, genesisState)
	got := coinz.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.NotNil(t, got.Admin)
	require.Equal(t, adminAddress, got.Admin.Address)

	// this line is used by starport scaffolding # genesis/test/assert
}
