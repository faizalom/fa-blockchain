package main

import (
	"bytes"
	"encoding/gob"
	"fa-blockchain/src/blockchain"
	"fa-blockchain/src/grpcclient"
	"fa-blockchain/src/mempool"
	"fa-blockchain/src/utils"
	"fmt"
	"log"
	"os"
	"slices"
	"time"
)

var (
	// LastHash        = []byte("0000333d75f6b3929d5b91727fa629da3a6b865a58012881069fa8341c204c9f")
	// LastBlockHeight = 1
	ch = make(chan bool) // Channel to send mined blocks
)
var isMinerFree bool

func main() {
	utils.SetEnv()
	blockchain.InitDB()

	utils.LogErrors(os.Getenv("ERROR_LOG_FILE"))
	log.Println("Starting miner service... ", time.Now())
	fmt.Println("Starting miner service")

	CoinbaseTxAddress := "1AGHuAjXbsdRnaKwXcWNnbWzpz4dD4AxFF"
	chain := blockchain.InitBlockChain(CoinbaseTxAddress)
	fmt.Printf("Current Block Hash: %x, Current Block Height: %d \n", chain.LastHash, chain.LastBlockHeight)

	isMinerFree = true
	go grpcclient.VerifyAddBlock(chain, func(bc *blockchain.BlockChain) {
		isMinerFree = true
		mineBlock(bc)
	}, ch)

	// Subscribe to the topic and receive a stream
	stream, err := grpcclient.NewSubscription("create_block")
	if err != nil {
		panic("Failed to subscribe")
	}

	// Handle incoming messages in a separate goroutine
	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Fatalf("error receiving message: %v", err)
		}
		fmt.Printf("new_mining: SenderId: %s, EventId: %s\n", msg.GetSenderId(), msg.GetEventId())

		// Deserialize the message to get the transaction
		tx, err := Deserialize(msg.GetBinary())
		if err != nil {
			log.Printf("error decoding base64 data: %v \n", err)
			continue
		}
		mempool.AddToMempool(&tx)
		// Trigger another function after adding to mempool

		// 	block := blockchain.Deserialize(msg.GetBinary())
		// 	blockchain.GetBlockProof(block)
		// 	grpcclient.PublishVerifyNewBlock(block)

		go mineBlock(chain)
	}
}

func mineBlock(chain *blockchain.BlockChain) {
	if isMinerFree {
		// This function would contain the logic to mine a block
		// It would typically involve creating a new block, adding transactions from the mempool,
		// and then running the proof of work algorithm to find a valid hash.
		// After mining, it would publish the new block to the network.

		mempoolTrans := mempool.GetMempool()
		for i, tx := range mempoolTrans {
			verify := blockchain.VerifyTransaction(tx)
			if !verify {
				mempoolTrans = slices.Delete(mempoolTrans, i, i+1)
				// mempoolTrans = append(mempoolTrans[:i], mempoolTrans[i+1:]...)
			}
			tx.LockTransactionToUXTO()
		}
		if len(mempoolTrans) == 0 {
			log.Println("No transactions in mempool to mine")
			return
		}

		isMinerFree = false
		fmt.Printf("****************** New Mining Started with %d trans ******************\n", len(mempoolTrans))
		block := blockchain.MineBlock(mempoolTrans, chain.LastHash, chain.LastBlockHeight, ch)
		if block.Hash == nil {
			log.Println("Block mining failed, other miner might have mined it first")
			return
		}
		log.Printf("Block mined successfully: %x, Height: %d", block.Hash, block.BlockHeight)
		grpcclient.PublishVerifyNewBlock(block)
		isMinerFree = true
		chain.SaveBlock(block)
	}
}

func Deserialize(data []byte) (blockchain.Transaction, error) {
	var trans blockchain.Transaction

	if data == nil {
		return trans, fmt.Errorf("no data to deserialize")
	}

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&trans)
	if err != nil {
		return trans, fmt.Errorf("error decoding transaction: %v", err)
	}

	return trans, nil
}
