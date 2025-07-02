package main

import (
	"fa-blockchain/server/wallet/controllers"
	"fa-blockchain/src/blockchain"
	"fa-blockchain/src/grpcclient"
	"fa-blockchain/src/utils"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func Routes() http.Handler {
	walletController := controllers.NewWalletController()

	mux := http.NewServeMux()
	// Wallet
	mux.HandleFunc("POST /api/wallet", walletController.CreateWallet)
	mux.HandleFunc("POST /api/send-amount", walletController.SendAmount)

	return mux
}

func main() {
	utils.SetEnv()
	blockchain.InitDB()
	utils.LogErrors(os.Getenv("ERROR_LOG_FILE"))
	log.Println("Starting wallet server... ", time.Now())

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = ":8080"
	}

	CoinbaseTxAddress := "1AGHuAjXbsdRnaKwXcWNnbWzpz4dD4AxFF"
	chain := blockchain.InitBlockChain(CoinbaseTxAddress)
	fmt.Printf("Current Block Hash: %x, Current Block Height: %d \n", chain.LastHash, chain.LastBlockHeight)

	go transLogger()

	go grpcclient.VerifyAddBlock(chain, func(bc *blockchain.BlockChain) {
		// This callback is intentionally left empty because no additional action is required
		// after VerifyAddBlock completes in this context.
	})

	fmt.Printf("Starting wallet server on %s\n", port)
	if err := http.ListenAndServe(port, Routes()); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

func transLogger() {
	stream, err := grpcclient.NewSubscription("trans_logger")
	if err != nil {
		panic("Failed to subscribe")
	}

	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Fatalf("error receiving message: %v", err)
			return
		}
		fmt.Printf("verify_block: SenderId: %s, EventId: %s\n", msg.GetSenderId(), msg.GetEventId())
	}
}
