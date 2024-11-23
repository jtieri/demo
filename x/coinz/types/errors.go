package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/coinz module sentinel errors
var (
	ErrInvalidSigner          = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrInvalidAdmin           = sdkerrors.Register(ModuleName, 1102, "invalid admin")
	ErrCannotUpdateAdmin      = sdkerrors.Register(ModuleName, 1103, "cannot update admin")
	ErrCannotMint             = sdkerrors.Register(ModuleName, 1104, "cannot mint")
	ErrInvalidAmount          = sdkerrors.Register(ModuleName, 1105, "invalid amount")
	ErrInvalidAssetDenom      = sdkerrors.Register(ModuleName, 1106, "invalid asset denom")
	ErrAssetMetadataNotFound  = sdkerrors.Register(ModuleName, 1107, "asset metadata not found")
	ErrInvalidAssetMetadata   = sdkerrors.Register(ModuleName, 1108, "invalid asset metadata")
	ErrCannotUpdatePauseState = sdkerrors.Register(ModuleName, 1109, "cannot update pause state")
	ErrPauseStateNotFound     = sdkerrors.Register(ModuleName, 1110, "pause state not found")
	ErrPauseStateNotUpdated   = sdkerrors.Register(ModuleName, 1111, "pause state not updated")
)
