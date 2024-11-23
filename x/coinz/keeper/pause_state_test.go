package keeper_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/jtieri/demo/testutil/keeper"
	"github.com/jtieri/demo/testutil/nullify"
	"github.com/jtieri/demo/x/coinz/keeper"
	"github.com/jtieri/demo/x/coinz/types"
)

func createTestPauseState(keeper keeper.Keeper, ctx context.Context) types.PauseState {
	item := types.PauseState{Paused: true}
	keeper.SetPauseState(ctx, item)
	return item
}

func TestPauseStateGet(t *testing.T) {
	keeper, ctx := keepertest.CoinzKeeper(t)
	item := createTestPauseState(keeper, ctx)
	rst, found := keeper.GetPauseState(ctx)
	require.True(t, found)
	require.Equal(t,
		nullify.Fill(&item),
		nullify.Fill(&rst),
	)
}
