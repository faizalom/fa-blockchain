package utils

import "errors"

var ErrWalletNotFound = errors.New("wallet not found")
var ErrInvalidPublicKey = errors.New("invalid public key")
var ErrInvalidPrivateKey = errors.New("invalid private key")
var ErrNotEnoughFunds = errors.New("not enough funds")
var ErrInvalidPrevTransaction = errors.New("invalid previous transaction")
var ErrTransactionNotFound = errors.New("transaction does not exist") // WIP

// Broker ERR
var ErrMQBrokerUnavailable = errors.New("message broker is down")
var ErrTopicNotFound = errors.New("topic not found")
var ErrCreateTransTopicNotFound = ErrTopicNotFound

var ErrString map[error]string

func init() {
	ErrString = make(map[error]string)
	ErrString[ErrTopicNotFound] = "Topic not found"
	ErrString[ErrCreateTransTopicNotFound] = "Transaction server is down"
	ErrString[ErrMQBrokerUnavailable] = "Message broker is down"
}

func GetErrString(err error) string {
	str, ok := ErrString[err]
	if !ok {
		return err.Error()
	}

	return str
}

// var ErrInvalidTransaction = errors.New("invalid transaction")
// var ErrInsufficientFunds = errors.New("insufficient funds")
// var ErrInvalidAddress = errors.New("invalid address")
// var ErrInvalidSignature = errors.New("invalid signature")
// var ErrTransactionNotFound = errors.New("transaction not found")
// var ErrInvalidInput = errors.New("invalid input")
// var ErrInvalidOutput = errors.New("invalid output")
// var ErrInvalidBlock = errors.New("invalid block")
// var ErrInvalidChain = errors.New("invalid chain")
// var ErrInvalidGenesisBlock = errors.New("invalid genesis block")
// var ErrInvalidBlockHash = errors.New("invalid block hash")
// var ErrInvalidBlockHeight = errors.New("invalid block height")
// var ErrInvalidBlockTime = errors.New("invalid block time")
// var ErrInvalidBlockNonce = errors.New("invalid block nonce")
// var ErrInvalidBlockTransactions = errors.New("invalid block transactions")
// var ErrInvalidBlockMerkleRoot = errors.New("invalid block merkle root")
// var ErrInvalidBlockDifficulty = errors.New("invalid block difficulty")
// var ErrInvalidBlockVersion = errors.New("invalid block version")
// var ErrInvalidBlockPreviousHash = errors.New("invalid block previous hash")
// var ErrInvalidBlockReward = errors.New("invalid block reward")
// var ErrInvalidBlockFees = errors.New("invalid block fees")
// var ErrInvalidBlockMiner = errors.New("invalid block miner")
// var ErrInvalidBlockSignature = errors.New("invalid block signature")
// var ErrInvalidBlockPublicKey = errors.New("invalid block public key")
// var ErrInvalidBlockPrivateKey = errors.New("invalid block private key")
// var ErrInvalidBlockAddress = errors.New("invalid block address")
// var ErrInvalidBlockData = errors.New("invalid block data")
// var ErrInvalidBlockDataHash = errors.New("invalid block data hash")
// var ErrInvalidBlockDataSignature = errors.New("invalid block data signature")
// var ErrInvalidBlockDataPublicKey = errors.New("invalid block data public key")
// var ErrInvalidBlockDataPrivateKey = errors.New("invalid block data private key")
// var ErrInvalidBlockDataAddress = errors.New("invalid block data address")
// var ErrInvalidBlockDataSignatureHash = errors.New("invalid block data signature hash")
// var ErrInvalidBlockDataSignaturePublicKey = errors.New("invalid block data signature public key")
// var ErrInvalidBlockDataSignaturePrivateKey = errors.New("invalid block data signature private key")
// var ErrInvalidBlockDataSignatureAddress = errors.New("invalid block data signature address")
// var ErrInvalidBlockDataSignatureData = errors.New("invalid block data signature data")
// var ErrInvalidBlockDataSignatureDataHash = errors.New("invalid block data signature data hash")
// var ErrInvalidBlockDataSignatureDataPublicKey = errors.New("invalid block data signature data public key")
// var ErrInvalidBlockDataSignatureDataPrivateKey = errors.New("invalid block data signature data private key")
