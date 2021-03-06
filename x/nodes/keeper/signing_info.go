package keeper

import (
	"github.com/pokt-network/pocket-core/x/nodes/types"
	sdk "github.com/pokt-network/posmint/types"
)

// get signing information for the validator by address
func (k Keeper) GetValidatorSigningInfo(ctx sdk.Context, addr sdk.Address) (info types.ValidatorSigningInfo, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetValidatorSigningInfoKey(addr))
	if bz == nil {
		found = false
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &info)
	found = true
	return
}

// set signing information for the validator by address
func (k Keeper) SetValidatorSigningInfo(ctx sdk.Context, addr sdk.Address, info types.ValidatorSigningInfo) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(info)
	store.Set(types.GetValidatorSigningInfoKey(addr), bz)
}

func (k Keeper) IterateAndExecuteOverValSigningInfo(ctx sdk.Context,
	handler func(addr sdk.Address, info types.ValidatorSigningInfo) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ValidatorSigningInfoKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		address := types.GetValidatorSigningInfoAddress(iter.Key())
		var info types.ValidatorSigningInfo
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &info)
		if handler(address, info) {
			break
		}
	}
}

func (k Keeper) valMissedAt(ctx sdk.Context, addr sdk.Address, index int64) (missed bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetValMissedBlockKey(addr, index))
	if bz == nil { // lazy: treat empty key as not missed
		missed = false
		return
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &missed)
	return
}

func (k Keeper) SetValidatorMissedAt(ctx sdk.Context, addr sdk.Address, index int64, missed bool) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(missed)
	store.Set(types.GetValMissedBlockKey(addr, index), bz)
}

func (k Keeper) clearValidatorMissed(ctx sdk.Context, addr sdk.Address) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetValMissedBlockPrefixKey(addr))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}

// Stored by *validator* address (not operator address)
func (k Keeper) IterateAndExecuteOverMissedArray(ctx sdk.Context,
	address sdk.Address, handler func(index int64, missed bool) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	index := int64(0)
	// Array may be sparse
	for ; index < k.SignedBlocksWindow(ctx); index++ {
		var missed bool
		bz := store.Get(types.GetValMissedBlockKey(address, index))
		if bz == nil {
			continue
		}
		k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &missed)
		if handler(index, missed) {
			break
		}
	}
}
