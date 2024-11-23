package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/jtieri/demo/x/coinz/types"
)

// GetPauseState attempts to retrieve the PauseState from the KV store.
// If the PauseState is not initialized it returns an empty struct and false,
// otherwise it returns the PauseState and true.
func (k Keeper) GetPauseState(ctx context.Context) (types.PauseState, bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	bz := store.Get(types.KeyPrefix(types.PauseStateKey))
	if bz == nil {
		return types.PauseState{}, false
	}

	pauseState := types.PauseState{}
	k.cdc.MustUnmarshal(bz, &pauseState)

	return pauseState, true
}

// SetPauseState attempts to set the specified PauseState in the KV store.
// The method will panic if it fails to marshal the PauseState or fails to set the PauseState in the KV store.
func (k Keeper) SetPauseState(ctx context.Context, pauseState types.PauseState) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := k.cdc.MustMarshal(&pauseState)
	store.Set(types.KeyPrefix(types.PauseStateKey), bz)
}
