package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/jtieri/demo/x/coinz/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PauseState handles QueryGetPauseStateRequests and attempts to retrieve the PauseState from the KV store.
// If successful, it returns a response containing the PauseState.
func (k Keeper) PauseState(goCtx context.Context, req *types.QueryGetPauseStateRequest) (*types.QueryGetPauseStateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetPauseState(ctx)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetPauseStateResponse{Paused: val.Paused}, nil
}
