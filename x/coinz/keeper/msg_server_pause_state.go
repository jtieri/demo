package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/jtieri/demo/x/coinz/types"
)

// UpdatePauseState attempts to update the coinz modules PauseState, which represents if certain state changes can take place.
// If the message sender address does not equal the current Admin address then the msg fails to execute.
func (k msgServer) UpdatePauseState(goCtx context.Context, msg *types.MsgUpdatePauseState) (*types.MsgUpdatePauseStateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate MsgUpdatePauseState is being sent from the admin account.
	admin, found := k.GetAdmin(ctx)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrInvalidAdmin, "the admin is not initialized")
	}

	if msg.From != admin.Address {
		return nil, sdkerrors.Wrapf(types.ErrCannotUpdatePauseState, "address %s is not authorized to pause/unpause", msg.From)
	}

	// Check if the PauseState exists.
	pauseState, found := k.GetPauseState(ctx)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrPauseStateNotFound, "pause state may not be initialized")
	}

	// No need to update the PauseState if the new proposed value is the same as the one already in the KV store.
	if pauseState.Paused == msg.Paused {
		return nil, sdkerrors.Wrapf(types.ErrPauseStateNotUpdated, "pause state is already set to %v", msg.Paused)
	}

	// Update the PauseState.
	updatedState := types.PauseState{
		Paused: msg.Paused,
	}

	k.SetPauseState(ctx, updatedState)

	return &types.MsgUpdatePauseStateResponse{}, ctx.EventManager().EmitTypedEvent(msg)
}
