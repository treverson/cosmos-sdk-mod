package app

import (
	"encoding/json"

	abci "github.com/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
	cmn "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	"github.com/cosmos/cosmos-sdk/examples/basecoin/types"
	"github.com/cosmos/cosmos-sdk/x/sideBlock"
)

const (
	appName = "maticApp"
)

// Extended ABCI application
type maticApp struct {
	*bam.BaseApp
	cdc *wire.Codec

	// keys to access the substores
	keyMain     *sdk.KVStoreKey
	keyAccount  *sdk.KVStoreKey
	keyIBC      *sdk.KVStoreKey
	//keyStake    *sdk.KVStoreKey
	//keySlashing *sdk.KVStoreKey
	keySideBlock *sdk.KVStoreKey

	// Manage getting and setting accounts
	accountMapper       auth.AccountMapper
	feeCollectionKeeper auth.FeeCollectionKeeper
	coinKeeper          bank.Keeper
	ibcMapper           ibc.Mapper
	// removing staking for now
	//stakeKeeper         stake.Keeper
	//slashingKeeper      slashing.Keeper
	sideBlockKeeper		sideBlock.Keeper
}

func NewmaticApp(logger log.Logger, db dbm.DB) *maticApp {

	// Create app-level codec for txs and accounts.
	var cdc = MakeCodec()

	// Create your application object.
	var app = &maticApp{
		BaseApp:     bam.NewBaseApp(appName, cdc, logger, db),
		cdc:         cdc,
		keyMain:     sdk.NewKVStoreKey("main"),
		keyAccount:  sdk.NewKVStoreKey("acc"),
		keyIBC:      sdk.NewKVStoreKey("ibc"),
		keySideBlock:sdk.NewKVStoreKey("sideBlock"),
		//keyStake:    sdk.NewKVStoreKey("stake"),
		//keySlashing: sdk.NewKVStoreKey("slashing"),

	}

	// Define the accountMapper.
	app.accountMapper = auth.NewAccountMapper(
		cdc,
		app.keyAccount,      // target store
		&types.AppAccount{}, // prototype
	)

	// add accountMapper/handlers
	app.coinKeeper = bank.NewKeeper(app.accountMapper)
	app.ibcMapper = ibc.NewMapper(app.cdc, app.keyIBC, app.RegisterCodespace(ibc.DefaultCodespace))
	//app.stakeKeeper = stake.NewKeeper(app.cdc, app.keyStake, app.coinKeeper, app.RegisterCodespace(stake.DefaultCodespace))
	//app.slashingKeeper = slashing.NewKeeper(app.cdc, app.keySlashing, app.stakeKeeper, app.RegisterCodespace(slashing.DefaultCodespace))
	app.sideBlockKeeper = sideBlock.NewKeeper(app.cdc,app.keySideBlock,app.RegisterCodespace(sideBlock.DefaultCodespace))
	// register message routes
	app.Router().
		AddRoute("auth", auth.NewHandler(app.accountMapper)).
		AddRoute("bank", bank.NewHandler(app.coinKeeper)).
		AddRoute("ibc", ibc.NewHandler(app.ibcMapper, app.coinKeeper)).
		//AddRoute("stake", stake.NewHandler(app.stakeKeeper))
		AddRoute("sideBlock",sideBlock.NewHandler(app.sideBlockKeeper))

	// Initialize BaseApp.
	app.SetInitChainer(app.initChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)
	app.SetAnteHandler(auth.NewAnteHandler(app.accountMapper, app.feeCollectionKeeper))
	app.MountStoresIAVL(app.keyMain, app.keyAccount, app.keyIBC, app.keySideBlock)
	err := app.LoadLatestVersion(app.keyMain)
	if err != nil {
		cmn.Exit(err.Error())
	}
	return app
}


// Custom tx codec
func MakeCodec() *wire.Codec {
	var cdc = wire.NewCodec()
	wire.RegisterCrypto(cdc) // Register crypto.
	sdk.RegisterWire(cdc)    // Register Msgs
	bank.RegisterWire(cdc)
	//stake.RegisterWire(cdc)
	//slashing.RegisterWire(cdc)
	sideBlock.RegisterWire(cdc)
	ibc.RegisterWire(cdc)

	// register custom AppAccount
	cdc.RegisterInterface((*auth.Account)(nil), nil)
	//TODO Change this to respective matic form 
	cdc.RegisterConcrete(&types.AppAccount{}, "basecoin/Account", nil)
	return cdc
}

// application updates every end block
func (app *maticApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	// TODO check what was happening here
	//tags := slashing.BeginBlocker(ctx, req, app.slashingKeeper)

	//return abci.ResponseBeginBlock{
	//	Tags: tags.ToKVPairs(),
	//}
	return abci.ResponseBeginBlock{}
}

// application updates every end block
func (app *maticApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	// TODO fetch this from main chain
	//validatorUpdates := stake.EndBlocker(ctx, app.stakeKeeper)

	//return abci.ResponseEndBlock{
	//	ValidatorUpdates: validatorUpdates,
	//}
	return abci.ResponseEndBlock{}
}

// Custom logic for basecoin initialization
func (app *maticApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes

	genesisState := new(types.GenesisState)
	err := app.cdc.UnmarshalJSON(stateJSON, genesisState)
	if err != nil {
		panic(err) // TODO https://github.com/cosmos/cosmos-sdk/issues/468
		// return sdk.ErrGenesisParse("").TraceCause(err, "")
	}

	for _, gacc := range genesisState.Accounts {
		acc, err := gacc.ToAppAccount()
		if err != nil {
			panic(err) // TODO https://github.com/cosmos/cosmos-sdk/issues/468
			//	return sdk.ErrGenesisParse("").TraceCause(err, "")
		}
		acc.AccountNumber = app.accountMapper.GetNextAccountNumber(ctx)
		app.accountMapper.SetAccount(ctx, acc)
	}

	// load the initial stake information
	// TODO fetch from main chain (eth)
	//stake.InitGenesis(ctx, app.stakeKeeper, genesisState.StakeData)

	return abci.ResponseInitChain{}
}

// Custom logic for state export
func (app *maticApp) ExportAppStateAndValidators() (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	ctx := app.NewContext(true, abci.Header{})

	// iterate to get the accounts
	accounts := []*types.GenesisAccount{}
	appendAccount := func(acc auth.Account) (stop bool) {
		account := &types.GenesisAccount{
			Address: acc.GetAddress(),
			Coins:   acc.GetCoins(),
		}
		accounts = append(accounts, account)
		return false
	}
	app.accountMapper.IterateAccounts(ctx, appendAccount)

	genState := types.GenesisState{
		Accounts: accounts,
	}
	appState, err = wire.MarshalJSONIndent(app.cdc, genState)
	if err != nil {
		return nil, nil, err
	}
	//validators = stake.WriteValidators(ctx, app.stakeKeeper)
	return appState, validators, err
}
