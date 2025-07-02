package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"time"
)

type Block struct {
	Hash         []byte
	Transactions []*Transaction
	PrevHash     []byte
	Nonce        int
	Timestamp    int64
	BlockHeight  int
}

func GetBlock(blockHeight int) (*Block, error) {
	data, err := GetData(fmt.Sprintf("%d", blockHeight))
	if err != nil {
		return nil, err
	}
	block := Deserialize(data)
	return block, nil
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.Hash())
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

func (b *Block) SolvePuzzle(done <-chan bool) *Block {
	pow := NewProof(b)
	nonce, hash := pow.Run(done)

	b.Hash = hash[:]
	b.Nonce = nonce

	return b
}

func MineBlock(txs []*Transaction, prevHash []byte, prevHashHeight int, done <-chan bool) *Block {
	block := &Block{
		Hash:         []byte{},
		Transactions: txs,
		PrevHash:     prevHash,
		Nonce:        0,
		Timestamp:    time.Now().Unix(),
		BlockHeight:  prevHashHeight + 1,
	}

	return block.SolvePuzzle(done)
}

func Genesis(coinbase *Transaction) *Block {
	var done chan bool = nil
	return MineBlock([]*Transaction{coinbase}, []byte{}, 0, done)
}

func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)

	Handle(err)

	return res.Bytes()
}

func (b *Block) SerializeJson() []byte {

	transactions := []struct {
		TxID    string
		Fee     int
		Inputs  []any
		Outputs []any
	}{}

	for _, tx := range b.Transactions {
		txInput := []struct {
			ID        string
			Out       int
			Signature string
			PubKey    string
		}{}

		for _, in := range tx.Inputs {
			txInput = append(txInput, struct {
				ID        string
				Out       int
				Signature string
				PubKey    string
			}{
				ID:        fmt.Sprintf("%x", in.ID),
				Out:       in.Out,
				Signature: fmt.Sprintf("%x", in.Signature),
				PubKey:    fmt.Sprintf("%x", in.PubKey),
			})
		}

		txOutput := []struct {
			Value      int
			PubKeyHash string
		}{}

		for _, out := range tx.Outputs {
			txOutput = append(txOutput, struct {
				Value      int
				PubKeyHash string
			}{
				Value:      out.Value,
				PubKeyHash: fmt.Sprintf("%x", out.PubKeyHash),
			})
		}

		// Convert txInput to []any
		txInputAny := make([]any, len(txInput))
		for i, v := range txInput {
			txInputAny[i] = v
		}

		// Convert txOutput to []any
		txOutputAny := make([]any, len(txOutput))
		for i, v := range txOutput {
			txOutputAny[i] = v
		}

		transactions = append(transactions, struct {
			TxID    string
			Fee     int
			Inputs  []any
			Outputs []any
		}{
			TxID:    fmt.Sprintf("%x", tx.TxID),
			Fee:     tx.Fee,
			Inputs:  txInputAny,
			Outputs: txOutputAny,
		})
	}

	bj := struct {
		Hash         string
		Transactions any
		PrevHash     string
		Nonce        int
		Timestamp    int64
		BlockHeight  int
	}{
		fmt.Sprintf("%x", b.Hash),
		transactions,
		fmt.Sprintf("%x", b.PrevHash),
		b.Nonce,
		b.Timestamp,
		b.BlockHeight,
	}

	jsonData, err := json.Marshal(bj)
	Handle(err)
	return jsonData
}

func Deserialize(data []byte) *Block {
	if data == nil {
		return nil
	}
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)
	if err != nil {
		fmt.Println("Error decoding block:", err)
	}

	return &block
}
