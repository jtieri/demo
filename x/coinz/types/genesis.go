package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
		Asset:  nil,
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate

	// If the admin address is set in genesis.json AFTER the validators run gentx then a panic will occur due to
	// AppModuleBasic.ValidateGenesis being called during gentx, due to this we allow the Admin to be nil.
	// The admin MUST be explicitly set in genesis.json before starting the chain.
	if gs.Admin != nil {
		if _, err := sdk.AccAddressFromBech32(gs.Admin.Address); err != nil {
			return sdkerrors.Wrapf(ErrInvalidAdmin, "admin address is invalid got(%s)", gs.Admin.Address)
		}
	}

	if gs.Asset != nil {
		if !gs.Asset.Asset.IsValid() {
			return sdkerrors.Wrapf(ErrInvalidAssetMetadata, "check denom and initial supply amount got(%s)", gs.Asset.Asset)
		}
	}

	return gs.Params.Validate()
}
