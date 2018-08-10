package sideBlock


import (
sdk "github.com/cosmos/cosmos-sdk/types"
"github.com/cosmos/cosmos-sdk/wire"
)

var cdc = wire.NewCodec()

// name to identify transaction types
const MsgType = "sideBlock"

// verify interface at compile time
var _ sdk.Msg = &MsgSideBlock{}

// MsgUnrevoke - struct for unrevoking revoked validator
type MsgSideBlock struct {
	// TODO variable as we dont know who will call this
	VariableAddress sdk.Address `json:"address"` // address of the validator owner
}

func NewMsgSideBlock(variableAddr sdk.Address) MsgSideBlock {
	return MsgSideBlock{
		VariableAddress: variableAddr,
	}
}

//nolint
func (msg MsgSideBlock) Type() string              { return MsgType }
func (msg MsgSideBlock) GetSigners() []sdk.Address { return []sdk.Address{msg.VariableAddress} }

// get the bytes for the message signer to sign on
func (msg MsgSideBlock) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(struct {
		VariableAddr string `json:"address"`
	}{
		VariableAddr: sdk.MustBech32ifyVal(msg.VariableAddress),
	})
	if err != nil {
		panic(err)
	}
	return b
}

// quick validity check
func (msg MsgSideBlock) ValidateBasic() sdk.Error {
	if msg.VariableAddress == nil {
		//TODO create error and return respective error here, right now it will allow nil
		//return ErrBadValidatorAddr(DefaultCodespace)
		return nil
	}
	return nil
}

