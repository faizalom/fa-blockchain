package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fa-blockchain/src/models"
	"fa-blockchain/src/utils"
	"fa-blockchain/src/wallet"
	"math/big"

	"fmt"
	"log"
)

type Transaction struct {
	TxID    []byte
	Fee     int
	Inputs  []TxInput
	Outputs []TxOutput
}

func (tx *Transaction) Hash() []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.TxID = []byte{}

	hash = sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

func (tx Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

func CoinbaseTx(toAddress string) *Transaction {
	randData := make([]byte, 20)
	_, err := rand.Read(randData)
	if err != nil {
		log.Println(err)
	}
	data := fmt.Sprintf("%x", randData)

	txin := TxInput{
		ID:        []byte{},
		Out:       -1,
		Signature: nil,
		PubKey:    []byte(data),
	}
	txout := NewTXOutput(100, toAddress)

	tx := Transaction{
		Inputs:  []TxInput{txin},
		Outputs: []TxOutput{*txout}}
	tx.TxID = tx.Hash()

	return &tx
}

func NewTransaction(fromPrivateKey, toAddr string, amount int) (*Transaction, error) {
	w, err := wallet.GetWallet(fromPrivateKey)
	if err != nil {
		return nil, utils.ErrWalletNotFound
	}

	pubKeyHash := wallet.PublicKeyHash(w.PublicKey)

	accAmount, unspentOuts, err := FindSpendableOutputs(pubKeyHash, amount)
	if err != nil {
		return nil, err
	}
	if accAmount < amount {
		return nil, utils.ErrNotEnoughFunds
	}

	var txInputs []TxInput
	var txOutputs []TxOutput

	for _, out := range unspentOuts {
		// WIP
		txID, err := hex.DecodeString(out.TxId)
		if err != nil {
			log.Println(err)
		}
		txInput := TxInput{
			ID:        txID,
			Out:       out.OutIndex,
			Signature: nil,
			PubKey:    w.PublicKey,
		}
		txInputs = append(txInputs, txInput)
	}

	txOutputs = append(txOutputs, *NewTXOutput(amount, toAddr))

	if accAmount > amount {
		txOutputs = append(txOutputs, *NewTXOutput(accAmount-amount, w.Address()))
	}

	tx := &Transaction{
		Inputs:  txInputs,
		Outputs: txOutputs,
	}
	tx.TxID = tx.Hash()
	return tx.SignTransaction(unspentOuts, w.PrivateKey)
}

func (tx *Transaction) LockTransactionToUXTO() {
	for _, in := range tx.Inputs {
		models.UpdateUTXOStatus(fmt.Sprintf("%x", in.ID), in.Out, models.Locked)
	}
}

func (tx *Transaction) AddTransactionToUXTO(blockheight int) {
	// senderPubKeyHash := wallet.PublicKeyHash(tx.Inputs[0].PubKey)
	for _, in := range tx.Inputs {
		models.UpdateUTXOStatus(fmt.Sprintf("%x", in.ID), in.Out, models.Spent)
	}

	txId := tx.TxID
	for i, out := range tx.Outputs {
		// status := models.Locked
		// if bytes.Equal(senderPubKeyHash, out.PubKeyHash) {
		status := models.Unspent
		//}
		utxo := models.UTXO{
			TxId:        fmt.Sprintf("%x", txId),
			OutIndex:    i,
			Amount:      out.Value,
			PubKeyHash:  fmt.Sprintf("%x", out.PubKeyHash),
			BlockHeight: blockheight,
			Status:      status,
		}
		err := utxo.Insert()
		if err != nil {
			log.Println(err)
		}
	}
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) error {
	if tx.IsCoinbase() {
		return nil
	}

	for _, in := range tx.Inputs {
		if prevTXs[hex.EncodeToString(in.ID)].TxID == nil {
			return utils.ErrInvalidPrevTransaction
		}
	}

	txCopy := tx.TrimmedCopy()

	for inId, in := range txCopy.Inputs {
		prevTX := prevTXs[hex.EncodeToString(in.ID)]
		txCopy.Inputs[inId].Signature = nil
		txCopy.Inputs[inId].PubKey = prevTX.Outputs[in.Out].PubKeyHash
		txCopy.TxID = txCopy.Hash()
		txCopy.Inputs[inId].PubKey = nil

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.TxID)
		if err != nil {
			return err
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Inputs[inId].Signature = signature
	}
	return nil
}

func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	for _, in := range tx.Inputs {
		if prevTXs[hex.EncodeToString(in.ID)].TxID == nil {
			log.Panic("Previous transaction not correct")
		}
	}

	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inId, in := range tx.Inputs {
		prevTx := prevTXs[hex.EncodeToString(in.ID)]
		txCopy.Inputs[inId].Signature = nil
		txCopy.Inputs[inId].PubKey = prevTx.Outputs[in.Out].PubKeyHash
		txCopy.TxID = txCopy.Hash()
		txCopy.Inputs[inId].PubKey = nil

		r := big.Int{}
		s := big.Int{}

		sigLen := len(in.Signature)
		r.SetBytes(in.Signature[:(sigLen / 2)])
		s.SetBytes(in.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(in.PubKey)
		x.SetBytes(in.PubKey[:(keyLen / 2)])
		y.SetBytes(in.PubKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if ecdsa.Verify(&rawPubKey, txCopy.TxID, &r, &s) == false {
			return false
		}
	}

	return true
}

func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	for _, in := range tx.Inputs {
		inputs = append(inputs, TxInput{in.ID, in.Out, nil, nil})
	}

	for _, out := range tx.Outputs {
		outputs = append(outputs, TxOutput{out.Value, out.PubKeyHash})
	}

	txCopy := Transaction{
		TxID:    tx.TxID,
		Fee:     tx.Fee,
		Inputs:  inputs,
		Outputs: outputs,
	}

	return txCopy
}
