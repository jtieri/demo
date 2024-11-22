package keeper

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/jtieri/demo/x/coinz/types"
)

// GetAssetDenom attempts to retrieve the assets denom string from the KV store.
// If the AssetMetadata is not initialized it returns an empty string and false,
// otherwise it returns the asset denom string and true.
func (k Keeper) GetAssetDenom(ctx sdk.Context) (string, bool) {
	asset, found := k.GetAssetMetadata(ctx)
	if !found {
		return "", false
	}

	return asset.Asset.Denom, true
}

// GetAssetInitialSupply attempts to retrieve the assets initial supply from the KV store.
// If the AssetMetadata is not initialized it returns an empty struct and false,
// otherwise it returns the assets initial supply and true.
func (k Keeper) GetAssetInitialSupply(ctx sdk.Context) (math.Int, bool) {
	asset, found := k.GetAssetMetadata(ctx)
	if !found {
		return math.Int{}, false
	}

	return asset.Asset.Amount, true
}

// GetAssetMetadata attempts to retrieve the AssetMetadata from the KV store.
// If the AssetMetadata is not initialized it returns an empty struct and false,
// otherwise it returns the AssetMetadata and true.
func (k Keeper) GetAssetMetadata(ctx sdk.Context) (types.AssetMetadata, bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	assetBz := store.Get(types.KeyPrefix(types.AssetMetadataKey))
	if assetBz == nil {
		return types.AssetMetadata{}, false
	}

	assetMetadata := types.AssetMetadata{}
	k.cdc.MustUnmarshal(assetBz, &assetMetadata)

	return assetMetadata, true
}

// SetAssetMetadata attempts to set the specified AssetMetadata in the KV store.
// The method will panic if it fails to marshal the AssetMetadata or fails to set the AssetMetadata in the KV store.
func (k Keeper) SetAssetMetadata(ctx sdk.Context, assetMetadata types.AssetMetadata) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := k.cdc.MustMarshal(&assetMetadata)
	store.Set(types.KeyPrefix(types.AssetMetadataKey), bz)
}
