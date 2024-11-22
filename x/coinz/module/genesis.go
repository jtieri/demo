package coinz

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/jtieri/demo/x/coinz/keeper"
	"github.com/jtieri/demo/x/coinz/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}

	if genState.Admin != nil {
		k.SetAdmin(ctx, *genState.Admin)
	}

	if genState.Asset != nil {
		k.SetAssetMetadata(ctx, *genState.Asset)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	admin, found := k.GetAdmin(ctx)
	if found {
		genesis.Admin = &admin
	}

	asset, found := k.GetAssetMetadata(ctx)
	if found {
		genesis.Asset = &asset
	}

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
