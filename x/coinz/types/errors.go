package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/coinz module sentinel errors
var (
	ErrInvalidSigner     = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrInvalidAdmin      = sdkerrors.Register(ModuleName, 1102, "invalid admin")
	ErrCannotUpdateAdmin = sdkerrors.Register(ModuleName, 1103, "cannot update admin")
)
