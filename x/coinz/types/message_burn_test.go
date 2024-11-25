package types

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/jtieri/demo/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgBurn_ValidateBasic(t *testing.T) {
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
		msg  MsgBurn
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgBurn{
				From:   "invalid_address",
				Amount: validAmount,
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "all values valid",
			msg: MsgBurn{
				From:   sample.AccAddress(),
				Amount: validAmount,
			},
		}, {
			name: "invalid asset denom",
			msg: MsgBurn{
				From:   sample.AccAddress(),
				Amount: invalidDenom,
			},
			err: ErrInvalidAmount,
		}, {
			name: "invalid asset amount - zero amount",
			msg: MsgBurn{
				From:   sample.AccAddress(),
				Amount: zeroAmount,
			},
			err: ErrInvalidAmount,
		}, {
			name: "invalid asset amount - negative amount",
			msg: MsgBurn{
				From:   sample.AccAddress(),
				Amount: negativeAmount,
			},
			err: ErrInvalidAmount,
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
