package blockchain

import (
	"bytes"
	"encoding/gob"
	"fa-blockchain/src/wallet"
)

type TxOutput struct {
	Value      int
	PubKeyHash []byte
}

// Remove if not needed
type TxOutputs struct {
	Outputs []TxOutput
}

type TxInput struct {
	ID        []byte
	Out       int
	Signature []byte
	PubKey    []byte
}

func (out *TxOutput) SetPubHash(address string) {
	pubKeyHash, _ := wallet.Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
}

func NewTXOutput(value int, address string) *TxOutput {
	txo := &TxOutput{
		Value:      value,
		PubKeyHash: nil,
	}
	txo.SetPubHash(address)

	return txo
}

func (outs TxOutputs) Serialize() []byte {
	var buffer bytes.Buffer

	encode := gob.NewEncoder(&buffer)
	err := encode.Encode(outs)
	Handle(err)

	return buffer.Bytes()
}
