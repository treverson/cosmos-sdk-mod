package sideBlock

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/abci/types"
)
//TODO implement begin and end blocker if needed , we might need endBlocker to push data to main chain
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, sk Keeper) (tags sdk.Tags) {
	return
}

