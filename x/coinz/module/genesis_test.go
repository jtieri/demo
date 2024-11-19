package coinz_test

import (
	"testing"

	keepertest "github.com/jtieri/demo/testutil/keeper"
	"github.com/jtieri/demo/testutil/nullify"
	coinz "github.com/jtieri/demo/x/coinz/module"
	"github.com/jtieri/demo/x/coinz/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.CoinzKeeper(t)
	coinz.InitGenesis(ctx, k, genesisState)
	got := coinz.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
