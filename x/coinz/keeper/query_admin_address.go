package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/jtieri/demo/x/coinz/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AdminAddress handles QueryAdminAddressRequests and attempts to retrieve the Admin from the KV store.
// If successful, it returns a response containing the Admins address.
// This should always return a valid response since the Admin MUST be initialized at genesis.
func (k Keeper) AdminAddress(goCtx context.Context, req *types.QueryAdminAddressRequest) (*types.QueryAdminAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	admin, found := k.GetAdmin(ctx)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryAdminAddressResponse{
		Address: admin.Address,
	}, nil
}
