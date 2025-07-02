package grpcclient

import (
	"fa-blockchain/src/blockchain"
	"fmt"
	"log"
	"os"
)

func VerifyAddBlock(chain *blockchain.BlockChain, mineBlock func(*blockchain.BlockChain), ch ...chan<- bool) {
	// This function will verify the new block received from the miner
	// It will check if the block is valid, if it has the correct hash, and if it contains valid transactions
	// If the block is valid, it will be added to the blockchain
	// If the block is invalid, it will be rejected

	// Subscribe to the topic and receive a stream,
	stream, err := NewSubscription("verify_block")
	if err != nil {
		panic("Failed to subscribe")
	}

	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Fatalf("error receiving message: %v", err)
		}
		fmt.Printf("verify_block: SenderId: %s, EventId: %s\n", msg.GetSenderId(), msg.GetEventId())

		if os.Getenv("DEVICE_ID") == msg.GetSenderId() {
			mineBlock(chain)

			log.Println("Skipping verification for block sent by this miner")
			continue // Skip verification if the block is sent by this miner
		}

		block := blockchain.Deserialize(msg.GetBinary())
		if block == nil {
			log.Println("Received an empty block, skipping verification")
			continue
		}

		if block.BlockHeight <= chain.LastBlockHeight {
			// log.Printf("Received block with height %d, but last block height is %d, skipping verification\n", block.BlockHeight, chain.LastBlockHeight)
			continue // Skip if the block height is not greater than the last block height
		} else {
			log.Printf("Received block with height: %d, winner: %s\n", block.BlockHeight, msg.GetSenderId())
		}

		status := blockchain.NewProof(block).Validate()
		if !status {
			log.Printf("Block is invalid: %t\n", status)
			continue // Skip to the next block if this one is invalid
		}

		for _, c := range ch {
			c <- true
		}
		chain.SaveBlock(block)
		mineBlock(chain)
	}
}
