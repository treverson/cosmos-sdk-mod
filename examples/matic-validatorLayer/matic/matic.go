package matic

import (
	bapp "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	dbm "github.com/tendermint/tmlibs/db"
	cmn "github.com/tendermint/tmlibs/common"
	"github.com/tendermint/tmlibs/log"
	"encoding/json"
)

const (
	appName = "maticApp"
)


//// Extended ABCI application
//type maticApp struct {
//	*bam.BaseApp
//	cdc *wire.Codec
//
//	// keys to access the substores
//	keyMain     *sdk.KVStoreKey
//	keyAccount  *sdk.KVStoreKey
//	keyIBC      *sdk.KVStoreKey
//	keyStake    *sdk.KVStoreKey
//	keySlashing *sdk.KVStoreKey
//
//	// Manage getting and setting accounts
//	accountMapper       auth.AccountMapper
//	feeCollectionKeeper auth.FeeCollectionKeeper
//	coinKeeper          bank.Keeper
//	ibcMapper           ibc.Mapper
//	stakeKeeper         stake.Keeper
//	slashingKeeper      slashing.Keeper
//}
func NewApp1(logger log.Logger, db dbm.DB) *bapp.BaseApp {

	cdc := wire.NewCodec()

	// Create the base application object.
	app := bapp.NewBaseApp(appName, cdc, logger, db)

	// Create a key for accessing the account store.
	keyAccount := sdk.NewKVStoreKey("acc")

	// Determine how transactions are decoded.
	app.SetTxDecoder(txDecoder)

	// Register message routes.
	// Note the handler gets access to the account store.
	app.Router().
		AddRoute("send", handleMsgSend(keyAccount))

	// Mount stores and load the latest state.
	app.MountStoresIAVL(keyAccount)
	err := app.LoadLatestVersion(keyAccount)
	if err != nil {
		cmn.Exit(err.Error())
	}
	return app
}
var _ sdk.Msg = MsgSend{}

// MsgSend to send coins from Input to Output
type MsgSend struct {
	From   sdk.Address `json:"from"`
	To     sdk.Address `json:"to"`
	Amount sdk.Coins      `json:"amount"`
}

// NewMsgSend
func NewMsgSend(from, to sdk.Address, amt sdk.Coins) MsgSend {
	return MsgSend{from, to, amt}
}

// Implements Msg.
func (msg MsgSend) Type() string { return "send" }

// Implements Msg. Ensure the addresses are good and the
// amount is positive.
func (msg MsgSend) ValidateBasic() sdk.Error {
	if len(msg.From) == 0 {
		return sdk.ErrInvalidAddress("From address is empty")
	}
	if len(msg.To) == 0 {
		return sdk.ErrInvalidAddress("To address is empty")
	}
	if !msg.Amount.IsPositive() {
		return sdk.ErrInvalidCoins("Amount is not positive")
	}
	return nil
}
// Implements Msg. JSON encode the message.
func (msg MsgSend) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(bz)
}

// Implements Msg. Return the signer.
func (msg MsgSend) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

