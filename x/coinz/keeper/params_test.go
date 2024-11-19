package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/jtieri/demo/testutil/keeper"
	"github.com/jtieri/demo/x/coinz/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.CoinzKeeper(t)
	params := types.DefaultParams()

	require.NoError(t, k.SetParams(ctx, params))
	require.EqualValues(t, params, k.GetParams(ctx))
}
