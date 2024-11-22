package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgMint{}

func NewMsgMint(from string, address string, amount sdk.Coin) *MsgMint {
	return &MsgMint{
		From:    from,
		Address: address,
		Amount:  amount,
	}
}

// ValidateBasic performs basic sanity checks on the msg fields.
func (msg *MsgMint) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid from address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Address)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid destination address (%s)", err)
	}

	if !msg.Amount.IsValid() {
		return errorsmod.Wrap(ErrInvalidAmount, "token specified for mint is invalid, check denom and amount")
	}

	if msg.Amount.Amount.LTE(math.ZeroInt()) {
		return errorsmod.Wrap(ErrInvalidAmount, "amount specified for mint must be greater than zero")
	}

	return nil
}
