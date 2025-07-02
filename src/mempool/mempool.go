package mempool

import (
	"fa-blockchain/src/blockchain"
	"fmt"
	"sync"
)

var mempool []*blockchain.Transaction
var mutex sync.Mutex

// Add transaction to mempool
func AddToMempool(tx *blockchain.Transaction) {
	mutex.Lock()
	mempool = append(mempool, tx)
	mutex.Unlock()
	fmt.Printf("Transaction added to memPool: %x: %d\n", tx.TxID, len(mempool))
}

// Get sorted mempool transactions by fee priority
func GetMempool() []*blockchain.Transaction {
	mutex.Lock()
	defer mutex.Unlock()

	defer func() {
		mempool = nil
	}()

	// Sort transactions by fee (highest first)
	// (Implement sorting logic here)
	return mempool
}

// // Get sorted mempool transactions by fee priority
// func GetMempool() []blockchain.Transaction {
// 	mutex.Lock()
// 	defer mutex.Unlock()

// 	// Sort transactions by fee (highest first)
// 	// (Implement sorting logic here)
// 	return mempool
// }

// func main() {
// 	// Simulating incoming transactions
// 	AddToMempool(Transaction{"tx123", 100, "A", "B"})
// 	AddToMempool(Transaction{"tx456", 50, "C", "D"})

// 	fmt.Println("Mempool transactions:", GetMempool())
// }
