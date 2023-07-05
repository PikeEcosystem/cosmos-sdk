package keeper

import (
	sdk "github.com/PikeEcosystem/cosmos-sdk/types"
	sdkerrors "github.com/PikeEcosystem/cosmos-sdk/types/errors"
	"github.com/PikeEcosystem/cosmos-sdk/x/collection"
)

func (k Keeper) GetParams(ctx sdk.Context) collection.Params {
	store := ctx.KVStore(k.storeKey)
	key := paramsKey
	bz := store.Get(key)
	if bz == nil {
		panic(sdkerrors.ErrNotFound.Wrap("params does not exist"))
	}

	var params collection.Params
	k.cdc.MustUnmarshal(bz, &params)

	return params
}

func (k Keeper) SetParams(ctx sdk.Context, params collection.Params) {
	store := ctx.KVStore(k.storeKey)
	key := paramsKey

	bz, err := params.Marshal()
	if err != nil {
		panic(err)
	}
	store.Set(key, bz)
}
