package sideBlock

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
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


	return sdk.Result{

	}
}
