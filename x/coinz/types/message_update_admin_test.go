package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/jtieri/demo/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateAdmin_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateAdmin
		err  error
	}{
		{
			name: "invalid from address",
			msg: MsgUpdateAdmin{
				From:    "invalid_address",
				Address: sample.AccAddress(),
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid addresses",
			msg: MsgUpdateAdmin{
				From:    sample.AccAddress(),
				Address: sample.AccAddress(),
			},
		}, {
			name: "invalid address",
			msg: MsgUpdateAdmin{
				From:    sample.AccAddress(),
				Address: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
