package blockchain

import (
	"fa-blockchain/src/models"
	"fmt"
	"log"
)

func FindSpendableOutputs(pubKeyHash []byte, amount int) (int, []models.UTXO, error) {
	accAmount := 0
	unspentOuts := []models.UTXO{}

	utxos, err := models.GetUTXOByPubKeyHash(fmt.Sprintf("%x", pubKeyHash))
	if err != nil {
		log.Println(err)
		return 0, unspentOuts, err
	}

	for _, utxo := range utxos {
		accAmount += utxo.Amount
		unspentOuts = append(unspentOuts, utxo)

		if amount <= accAmount {
			break
		}
	}

	return accAmount, unspentOuts, nil
}
