package sideBlock


import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

type Keeper struct {
	storeKey     sdk.StoreKey
	cdc          *wire.Codec
	//validatorSet sdk.ValidatorSet

	// codespace
	codespace sdk.CodespaceType
}

func NewKeeper(cdc *wire.Codec, key sdk.StoreKey, codespace sdk.CodespaceType) Keeper {
	keeper := Keeper{
		storeKey:   key,
		cdc:        cdc,
		codespace:  codespace,
	}
	return keeper
}

// will get you the block details provided with the block hash , cool !
func (k Keeper) getBlock(ctx sdk.Context, blockhash string) []byte {
	store := ctx.KVStore(k.storeKey)
	return store.Get([]byte(blockhash))
}
// will store the block hash under the block hash for now , later will store the block struct with key as block hash
func (k Keeper) addBlock(ctx sdk.Context,blockHash string)  {
	logger := ctx.Logger().With("module", "x/sideBlock")
	store := ctx.KVStore(k.storeKey)
	// TODO replace the second param with block struct and first will remain block hash
	// we are using block hash as the key here !
	store.Set([]byte(blockHash),[]byte(blockHash))
	logger.Info("oh okay so the logs so work, ctx is %s", ctx)
}