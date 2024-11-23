package keeper

import (
	"context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/jtieri/demo/x/coinz/types"
)

// Mint attempts to mint new coins in the coinz module account and then send them to the specified receiver address.
// If the message sender address does not equal the current Admin address then the msg fails to execute.
func (k msgServer) Mint(goCtx context.Context, msg *types.MsgMint) (*types.MsgMintResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the module is paused, if so minting cannot occur until it is not paused.
	pauseState, found := k.GetPauseState(ctx)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrPauseStateNotFound, "pause state may not be initialized")
	}

	if pauseState.Paused {
		return nil, sdkerrors.Wrap(types.ErrCannotMint, "coinz module is currently paused")
	}

	// Validate the MsgMint is being sent from the admin account.
	admin, found := k.GetAdmin(ctx)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrInvalidAdmin, "the admin is not initialized")
	}

	if msg.From != admin.Address {
		return nil, sdkerrors.Wrapf(types.ErrCannotMint, "address %s is not authorized to mint", msg.From)
	}

	// Check the AssetMetadata denom in state and validate it is equal to the msgs specified asset denom.
	assetDenom, found := k.GetAssetDenom(ctx)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrAssetMetadataNotFound, "asset metadata may not be initialized")
	}

	if assetDenom != msg.Amount.Denom {
		return nil, sdkerrors.Wrapf(types.ErrInvalidAssetDenom, "expected(%s), got(%s)", assetDenom, msg.Amount.Denom)
	}

	// Attempt to mint the new assets and send to the receiver address.
	mintAmount := sdk.NewCoins(msg.Amount)
	err := k.bankKeeper.MintCoins(ctx, types.ModuleName, mintAmount)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to mint coins %s", msg.Amount.String())
	}

	receiver, err := sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return nil, err
	}

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, mintAmount)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to send coins to receiver %s", receiver.String())
	}

	return &types.MsgMintResponse{}, ctx.EventManager().EmitTypedEvent(msg)
}
