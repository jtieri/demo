package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

// this line is used by starport scaffolding # genesis/types/import

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
		Admin:  nil,
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate

	if gs.Admin == nil {
		sdkerrors.Wrapf(ErrInvalidAdmin, "admin cannot be nil")
	}

	if _, err := sdk.AccAddressFromBech32(gs.Admin.Address); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "admin address is invalid got(%s)", gs.Admin.Address)
	}

	return gs.Params.Validate()
}
