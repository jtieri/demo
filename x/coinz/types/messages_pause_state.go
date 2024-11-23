package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdatePauseState{}

func NewMsgUpdatePauseState(from string, paused bool) *MsgUpdatePauseState {
	return &MsgUpdatePauseState{
		From:   from,
		Paused: paused,
	}
}

func (msg *MsgUpdatePauseState) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.From)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid from address (%s)", err)
	}
	return nil
}
