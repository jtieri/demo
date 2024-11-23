package keeper_test

import (
	"testing"

	"github.com/jtieri/demo/testutil/sample"
	"github.com/stretchr/testify/require"

	keepertest "github.com/jtieri/demo/testutil/keeper"
	"github.com/jtieri/demo/x/coinz/keeper"
	"github.com/jtieri/demo/x/coinz/types"
)

func TestPauseStateMsgServerUpdate(t *testing.T) {
	admin := types.Admin{Address: sample.AccAddress()}
	pauseState := types.PauseState{Paused: true}

	tests := []struct {
		desc    string
		request *types.MsgUpdatePauseState
		err     error
	}{
		{
			desc: "Completed",
			request: &types.MsgUpdatePauseState{
				From:   admin.Address,
				Paused: false,
			},
		},
		{
			desc: "Unauthorized",
			request: &types.MsgUpdatePauseState{
				From:   "B",
				Paused: true,
			},
			err: types.ErrCannotUpdatePauseState,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			k, ctx := keepertest.CoinzKeeper(t)
			srv := keeper.NewMsgServerImpl(k)

			k.SetAdmin(ctx, admin)
			k.SetPauseState(ctx, pauseState)

			_, err := srv.UpdatePauseState(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)

				ps, found := k.GetPauseState(ctx)
				require.True(t, found)
				require.Equal(t, ps.Paused, tc.request.Paused)
			}
		})
	}
}
