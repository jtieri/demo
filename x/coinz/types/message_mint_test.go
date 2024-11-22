package types

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/jtieri/demo/testutil/sample"
	"github.com/stretchr/testify/require"
)

// TestMsgMint_ValidateBasic asserts that the ValidateBasic method on MsgMint properly validates all msg fields.
func TestMsgMint_ValidateBasic(t *testing.T) {
	validAmount := sdk.Coin{
		Denom:  "myToken",
		Amount: math.NewInt(1234),
	}

	invalidDenom := sdk.Coin{
		Denom:  "",
		Amount: math.NewInt(1234),
	}

	negativeAmount := sdk.Coin{
		Denom:  "myToken",
		Amount: math.NewInt(-1234),
	}

	zeroAmount := sdk.Coin{
		Denom:  "myToken",
		Amount: math.NewInt(0),
	}

	tests := []struct {
		name string
		msg  MsgMint
		err  error
	}{
		{
			name: "invalid from address",
			msg: MsgMint{
				From:    "invalid_address",
				Address: sample.AccAddress(),
				Amount:  validAmount,
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "invalid receiver address",
			msg: MsgMint{
				From:    sample.AccAddress(),
				Address: "invalid_address",
				Amount:  validAmount,
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "invalid asset denom",
			msg: MsgMint{
				From:    sample.AccAddress(),
				Address: sample.AccAddress(),
				Amount:  invalidDenom,
			},
			err: ErrInvalidAmount,
		}, {
			name: "invalid asset amount - zero amount",
			msg: MsgMint{
				From:    sample.AccAddress(),
				Address: sample.AccAddress(),
				Amount:  zeroAmount,
			},
			err: ErrInvalidAmount,
		}, {
			name: "invalid asset amount - negative amount",
			msg: MsgMint{
				From:    sample.AccAddress(),
				Address: sample.AccAddress(),
				Amount:  negativeAmount,
			},
			err: ErrInvalidAmount,
		},
		{
			name: "All values valid",
			msg: MsgMint{
				From:    sample.AccAddress(),
				Address: sample.AccAddress(),
				Amount:  validAmount,
			},
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
