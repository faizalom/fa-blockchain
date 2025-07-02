package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fa-blockchain/src/models"
	"fa-blockchain/src/utils"
	"fmt"
	"log"
)

type BlockChain struct {
	LastHash        []byte
	LastBlockHeight int
}

type BlockChainIterator struct {
	CurrentBlockHeight int
}

func (chain *BlockChain) SaveBlock(block *Block) *Block {
	for _, tx := range block.Transactions {
		if VerifyTransaction(tx) != true {
			// WIP IMP
			log.Panic("Invalid Transaction")
		}
	}

	for _, tx := range block.Transactions {
		tx.AddTransactionToUXTO(block.BlockHeight)
		// models.UpdateBlockHeightByTxId(tx.TxID, block.BlockHeight)
	}

	chain.LastHash = block.Hash
	chain.LastBlockHeight = block.BlockHeight
	SaveToFile(block.Serialize(), fmt.Sprintf("%d", block.BlockHeight))
	SaveToFile(block.SerializeJson(), fmt.Sprintf("%d.json", block.BlockHeight))
	SaveToFile(IntToBytes(block.BlockHeight), "lh")

	return block
}

func IntToBytes(num int) []byte {
	return []byte(fmt.Sprintf("%d", num))
}

func InitBlockChain(address string) *BlockChain {
	chain := &BlockChain{}
	b := GetLastHash()
	if b != nil {
		chain.LastHash = b.Hash
		chain.LastBlockHeight = b.BlockHeight
	}
	if chain.LastBlockHeight == 0 {
		cbtx := CoinbaseTx(address)
		new := Genesis(cbtx)
		cbtx.AddTransactionToUXTO(1)
		// models.UpdateBlockHeightByTxId(new.Transactions[0].TxID, new.BlockHeight)

		SaveToFile(new.Serialize(), fmt.Sprintf("%d", new.BlockHeight))
		SaveToFile(new.SerializeJson(), fmt.Sprintf("%d.json", new.BlockHeight))
		SaveToFile(IntToBytes(new.BlockHeight), "lh")
		chain.LastHash = new.Hash
		chain.LastBlockHeight = new.BlockHeight
	}
	return chain
}

func GetLastHash() *Block {
	b, err := GetData("lh")
	if err != nil || b == nil {
		return nil
	}
	var height int
	fmt.Sscanf(string(b), "%d", &height)
	block, _ := GetBlock(height)
	return block
}

func (tx *Transaction) SignTransaction(utxos []models.UTXO, privKey ecdsa.PrivateKey) (*Transaction, error) {
	prevTXs := make(map[string]Transaction)

	for _, utxo := range utxos {
		block, err := GetBlock(utxo.BlockHeight)
		if err != nil {
			return nil, err
		}

		txId, _ := hex.DecodeString(utxo.TxId)
		prevTX, err := block.FindTransaction(txId)
		if err != nil {
			return nil, err
		}
		prevTXs[utxo.TxId] = prevTX
	}

	return tx, tx.Sign(privKey, prevTXs)
}

func VerifyTransaction(tx *Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}
	prevTXs := make(map[string]Transaction)

	for _, in := range tx.Inputs {
		blockheight, err := models.GetBlockHeightByTxId(in.ID)
		if err != nil {
			log.Println("Error getting block height by transaction ID:", err)
			return false
		}
		if blockheight == 0 {
			log.Println("Transaction already spent", in.ID)
			return false
		}

		bc, _ := GetBlock(blockheight)
		if bc == nil {
			log.Println("Block not found for transaction ID:", in.ID)
			return false
		}

		prevTX, err := bc.FindTransaction(in.ID)
		if err != nil {
			log.Println("Error finding previous transaction:", err)
			return false
		}

		prevTXs[hex.EncodeToString(prevTX.TxID)] = prevTX
	}

	return tx.Verify(prevTXs)
}

func (block *Block) FindTransaction(txID []byte) (Transaction, error) {
	for _, tx := range block.Transactions {
		if bytes.Compare(tx.TxID, txID) == 0 {
			return *tx, nil
		}
	}

	return Transaction{}, utils.ErrTransactionNotFound
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{chain.LastBlockHeight}
}

func (iter *BlockChainIterator) Next() *Block {
	block, _ := GetBlock(iter.CurrentBlockHeight)
	iter.CurrentBlockHeight = block.BlockHeight - 1
	return block
}
