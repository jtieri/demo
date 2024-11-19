package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/jtieri/demo/x/coinz/types"
)

// GetAdmin attempts to retrieve the Admin from the KV store.
// If the Admin is not initialized it returns an empty struct and false,
// otherwise it returns the Admin and true.
func (k Keeper) GetAdmin(ctx context.Context) (types.Admin, bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	adminBz := store.Get(types.KeyPrefix(types.AdminKey))
	if adminBz == nil {
		return types.Admin{}, false
	}

	admin := types.Admin{}
	k.cdc.MustUnmarshal(adminBz, &admin)

	return admin, true
}

// SetAdmin attempts to set the specified Admin in the KV store.
// The method will panic if it fails to marshal the Admin or fails to set the Admin in the KV store.
func (k Keeper) SetAdmin(ctx context.Context, admin types.Admin) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := k.cdc.MustMarshal(&admin)
	store.Set(types.KeyPrefix(types.AdminKey), bz)
}
