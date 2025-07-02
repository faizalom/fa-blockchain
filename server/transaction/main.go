package main

import (
	"bytes"
	"encoding/gob"
	"fa-blockchain/src/grpcclient"
	"fa-blockchain/src/utils"
	"fmt"
	"log"
	"os"
	"time"

	"fa-blockchain/src/blockchain"
)

func main() {
	utils.SetEnv()
	blockchain.InitDB()

	utils.LogErrors(os.Getenv("ERROR_LOG_FILE"))
	log.Println("Starting transaction service...", time.Now())
	CoinbaseTxAddress := "1AGHuAjXbsdRnaKwXcWNnbWzpz4dD4AxFF"
	chain := blockchain.InitBlockChain(CoinbaseTxAddress)
	fmt.Printf("Current Block Hash: %x, Current Block Height: %d \n", chain.LastHash, chain.LastBlockHeight)

	go grpcclient.VerifyAddBlock(chain, func(bc *blockchain.BlockChain) {
		// This callback is intentionally left empty because no additional action is required
		// after VerifyAddBlock completes in this context.
	})

	// Subscribe to the topic and receive a stream
	stream, err := grpcclient.NewSubscription("verify_trans")
	if err != nil {
		panic("Failed to subscribe")
	}

	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Fatalf("error receiving message: %v", err)
		}
		log.Printf("event ID: %s, sender Id: %s \n", msg.GetEventId(), msg.GetSenderId())

		tx, err := Deserialize(msg.GetBinary())
		if err != nil {
			log.Printf("error deserializing transaction: %v for event ID: %s", err, msg.GetEventId())
			continue
		}

		verify := blockchain.VerifyTransaction(tx)
		if !verify {
			log.Printf("Transaction verification failed for event ID: %s", msg.GetEventId())
			continue // Skip processing this transaction if verification fails
			// WIP
		}

		b, err := utils.Serialize(tx)
		if err != nil {
			log.Printf("error serializing transaction: %v for event ID: %s", err, msg.GetEventId())
			continue // Skip processing this transaction if serialization fails		}
			// WIP
		}

		_, err = grpcclient.PublishCreateBlock(b, msg.GetEventId())
		if err != nil {
			log.Printf("error publishing create block: %v for event ID: %s", err, msg.GetEventId())
			continue // Skip processing this transaction if publishing fails
			// WIP
		}
	}
}

func Deserialize(data []byte) (*blockchain.Transaction, error) {
	if data == nil {
		return nil, fmt.Errorf("no data to deserialize")
	}
	var trans *blockchain.Transaction

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&trans)
	if err != nil {
		return nil, fmt.Errorf("error decoding transaction: %v", err)
	}

	return trans, nil
}
