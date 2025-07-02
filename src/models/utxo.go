package models

import (
	"fmt"
)

type UTXO struct {
	BlockHeight int        `json:"block_height"`
	TxId        string     `json:"tx_id"`
	OutIndex    int        `json:"out_index"`
	Amount      int        `json:"amount"`
	PubKeyHash  string     `json:"pub_key_hash"`
	Status      UTXOStatus `json:"status"`
}

// Enum for UTXO status
type UTXOStatus string

const (
	Unspent UTXOStatus = "Unspent"
	Spent   UTXOStatus = "Spent"
	Locked  UTXOStatus = "Locked"
	Invalid UTXOStatus = "Invalid"
)

func (utxo *UTXO) Insert() error {
	db := Conn()
	// sql := "INSERT INTO utxos (block_height, tx_id, out_index, amount, pub_key_hash, status) VALUES ($1, $2, $3, $4, $5, $6)"
	sql := "INSERT INTO utxos (block_height, tx_id, out_index, amount, pub_key_hash, status) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(sql, utxo.BlockHeight, utxo.TxId, utxo.OutIndex, utxo.Amount, utxo.PubKeyHash, utxo.Status)
	return err
}

func UpdateUTXOStatus(txId string, outIndex int, status UTXOStatus) error {
	db := Conn()
	// sql := "UPDATE utxos SET status = $1 WHERE tx_id = $2 AND out_index = $3"
	sql := "UPDATE utxos SET status = ? WHERE tx_id = ? AND out_index = ?"
	_, err := db.Exec(sql, status, txId, outIndex)
	return err
}

// GetUTXOByPubKeyHash retrieves UTXOs by public key hash
func GetUTXOByPubKeyHash(pubKeyHash string) ([]UTXO, error) {
	db := Conn()
	// sql := "SELECT block_height, tx_id, out_index, amount FROM utxos WHERE pub_key_hash = $1 AND status = $2"
	sql := "SELECT block_height, tx_id, out_index, amount FROM utxos WHERE pub_key_hash = ? AND status = ?"
	rows, err := db.Query(sql, pubKeyHash, Unspent)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var utxos []UTXO
	for rows.Next() {
		var utxo UTXO
		if err := rows.Scan(&utxo.BlockHeight, &utxo.TxId, &utxo.OutIndex, &utxo.Amount); err != nil {
			return nil, err
		}
		utxos = append(utxos, utxo)
	}
	return utxos, nil
}

// Update blockHeight by txId
// func UpdateBlockHeightByTxId(txIdByte []byte, blockHeight int) error {
// 	txId := fmt.Sprintf("%x", txIdByte)
// 	db := Conn()
// 	// WIP
// 	_, err := db.Exec("UPDATE utxos SET block_height = $1, status = $2 WHERE tx_id = $3",
// 		blockHeight, Unspent, txId)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	return err
// }

func GetBlockHeightByTxId(txIdByte []byte) (int, error) {
	db := Conn()
	// sql := "SELECT block_height FROM utxos WHERE tx_id = $1 AND (status = $2 OR status = $3)"
	sql := "SELECT block_height FROM utxos WHERE tx_id = ? AND (status = ? OR status = ?)"
	txId := fmt.Sprintf("%x", txIdByte)
	row := db.QueryRow(sql, txId, Unspent, Locked)
	blockHeight := 0
	err := row.Scan(&blockHeight)
	return blockHeight, err
}
