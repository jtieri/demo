package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/jtieri/demo/x/coinz/types"
)

// UpdateAdmin attempts to update the current Admin address to the specified address.
// If the message sender address does not equal the current Admin address then the msg fails to execute.
func (k msgServer) UpdateAdmin(goCtx context.Context, msg *types.MsgUpdateAdmin) (*types.MsgUpdateAdminResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	admin, found := k.GetAdmin(ctx)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrInvalidAdmin, "the admin is not initialized")
	}

	if msg.From != admin.Address {
		return nil, sdkerrors.Wrapf(types.ErrCannotUpdateAdmin, "adress %s is not authorized to update the admin", msg.From)
	}

	admin.Address = msg.Address

	k.SetAdmin(ctx, admin)

	return &types.MsgUpdateAdminResponse{}, ctx.EventManager().EmitTypedEvent(msg)
}
