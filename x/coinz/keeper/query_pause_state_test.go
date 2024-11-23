package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/jtieri/demo/testutil/keeper"
	"github.com/jtieri/demo/testutil/nullify"
	"github.com/jtieri/demo/x/coinz/types"
)

func TestPauseStateQuery(t *testing.T) {
	keeper, ctx := keepertest.CoinzKeeper(t)
	item := createTestPauseState(keeper, ctx)
	tests := []struct {
		desc     string
		request  *types.QueryGetPauseStateRequest
		response *types.QueryGetPauseStateResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetPauseStateRequest{},
			response: &types.QueryGetPauseStateResponse{Paused: item.Paused},
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.PauseState(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}
