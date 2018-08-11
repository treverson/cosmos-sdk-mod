package sideBlock

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"bytes"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		// NOTE msg already has validate basic run
		switch msg := msg.(type) {
		case MsgSideBlock:
			return handleMsgSideBlock(ctx, msg, k)
		default:
			return sdk.ErrTxDecode("Invalid message in side block module ").Result()
		}
	}
}

func handleMsgSideBlock(ctx sdk.Context, msg MsgSideBlock, k Keeper) sdk.Result {
	// TODO make web3 call using block hash , validate data here and return true in tags if true
	blockHash := k.getBlock(ctx,msg.blockHash)
	if bytes.Equal(blockHash,blockHash){
		return sdk.Result{
			// TODO return block data here
		}
	} else {
		// TODO return error here 
		return sdk.Result{}
	}

}
