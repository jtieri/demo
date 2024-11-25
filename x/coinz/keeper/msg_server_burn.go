package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/jtieri/demo/x/coinz/types"
)

// Burn attempts to burn assets from the msg senders account by transferring them to the module account and burning them.
// If the specified denom does not equal the modules denom then the msg fails to execute.
func (k msgServer) Burn(goCtx context.Context, msg *types.MsgBurn) (*types.MsgBurnResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the module is paused, if so burning cannot occur until it is not paused.
	pauseState, found := k.GetPauseState(ctx)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrPauseStateNotFound, "pause state may not be initialized")
	}

	if pauseState.Paused {
		return nil, sdkerrors.Wrap(types.ErrCannotBurn, "coinz module is currently paused")
	}

	// Check the AssetMetadata denom in state and validate it is equal to the msgs specified asset denom.
	assetDenom, found := k.GetAssetDenom(ctx)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrAssetMetadataNotFound, "asset metadata may not be initialized")
	}

	if assetDenom != msg.Amount.Denom {
		return nil, sdkerrors.Wrapf(types.ErrInvalidAssetDenom, "expected(%s), got(%s)", assetDenom, msg.Amount.Denom)
	}

	// Check that the user has the appropriate balance to perform the burn.
	addr, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return nil, err
	}

	bal := k.bankKeeper.SpendableCoin(ctx, addr, msg.Amount.Denom)
	if bal.IsLT(msg.Amount) {
		return nil, sdkerrors.Wrapf(errors.ErrInsufficientFunds, "spendable amount %s is smaller than burn amount %s", bal.String(), msg.Amount.String())
	}

	// Move funds back to the module account and burn them.
	burnAmount := sdk.NewCoins(msg.Amount)
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, burnAmount)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to send coins from account to module")
	}

	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, burnAmount)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to burn coins")
	}

	return &types.MsgBurnResponse{}, ctx.EventManager().EmitTypedEvent(msg)
}
