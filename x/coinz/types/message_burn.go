package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgBurn{}

func NewMsgBurn(from string, amount sdk.Coin) *MsgBurn {
	return &MsgBurn{
		From:   from,
		Amount: amount,
	}
}

func (msg *MsgBurn) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid from address (%s)", err)
	}

	if !msg.Amount.IsValid() {
		return errorsmod.Wrap(ErrInvalidAmount, "token specified for burn is invalid, check denom and amount")
	}

	if msg.Amount.Amount.LTE(math.ZeroInt()) {
		return errorsmod.Wrap(ErrInvalidAmount, "amount specified for burn must be greater than zero")
	}

	return nil
}
