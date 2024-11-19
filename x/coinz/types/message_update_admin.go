package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdateAdmin{}

// NewMsgUpdateAdmin initialized a new MsgUpdateAdmin with the specified addresses.
func NewMsgUpdateAdmin(from string, address string) *MsgUpdateAdmin {
	return &MsgUpdateAdmin{
		From:    from,
		Address: address,
	}
}

// ValidateBasic performs basic sanity checks on the values of MsgUpdateAdmin.
func (msg *MsgUpdateAdmin) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid from address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address (%s)", err)
	}

	return nil
}
