package controllers

import (
	"crypto/x509"
	"encoding/json"
	"fa-blockchain/src/blockchain"
	"fa-blockchain/src/grpcclient"
	"fa-blockchain/src/utils"
	"fa-blockchain/src/wallet"
	"fmt"
	"net/http"
)

type WalletController struct {
	// BaseController
	// assetService services.IAssetService
}

type RespW map[string]any

func NewWalletController() *WalletController {
	return &WalletController{}
}

// CreateWallet handles the creation of a new wallet
func (WalletController) CreateWallet(w http.ResponseWriter, r *http.Request) {
	wlt := wallet.CreateNewWallet()

	privateKey, err := x509.MarshalECPrivateKey(&wlt.PrivateKey)
	if err != nil {
		utils.Error(w, utils.Message("Failed to marshal private key"), http.StatusInternalServerError)
		return
	}

	// Return the wallet address and private key and public key
	walletResponse := map[string]string{
		"address":     wlt.Address(),
		"private_key": fmt.Sprintf("%x", privateKey),
		"public_key":  fmt.Sprintf("%x", wlt.PublicKey),
	}

	utils.Success(w, RespW{"wallet": walletResponse})
}

func (WalletController) SendAmount(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the transaction details
	var requestBody struct {
		FromAddress    string `json:"from_address"`
		ToAddress      string `json:"to_address"`
		Amount         int    `json:"amount"`
		FromPrivateKey string `json:"from_private_key"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		utils.Error(w, utils.Message("Invalid request body"), http.StatusBadRequest)
		return
	}

	from := requestBody.FromAddress
	to := requestBody.ToAddress
	amount := requestBody.Amount

	if from == "" || to == "" || amount <= 0 {
		utils.Error(w, utils.Message("Invalid transaction details"), http.StatusBadRequest)
		return
	}

	if from == to {
		utils.Error(w, utils.Message("From and To addresses cannot be the same"), http.StatusBadRequest)
		return
	}
	// Validate the amount
	if amount <= 0 {
		utils.Error(w, utils.Message("Invalid amount"), http.StatusBadRequest)
		return
	}

	// Validate the from address
	ok := wallet.ValidateAddress(from)
	if !ok {
		utils.Error(w, utils.Message("Invalid from address"), http.StatusBadRequest)
		return
	}

	// Validate the to address
	ok = wallet.ValidateAddress(to)
	if !ok {
		utils.Error(w, utils.Message("Invalid to address"), http.StatusBadRequest)
		return
	}

	wallet, err := wallet.GetWallet(requestBody.FromPrivateKey)
	if err != nil {
		utils.Error(w, utils.Message("Invalid private key"), http.StatusBadRequest)
		return
	}
	if wallet.Address() != from {
		utils.Error(w, utils.Message("From address does not match the private key"), http.StatusBadRequest)
		return
	}

	// Create a new transaction
	tx, err := blockchain.NewTransaction(requestBody.FromPrivateKey, to, amount)
	if err != nil {
		utils.Error(w, utils.Message(utils.GetErrString(err)), http.StatusBadRequest)
		return
	}

	b, err := utils.Serialize(tx)
	if err != nil {
		utils.Error(w, utils.Message("Failed to serialize transaction"), http.StatusInternalServerError)
		return
	}

	eventId, err := grpcclient.PublishVerifyTransaction(b)
	if err != nil {
		response := RespW{
			"event_id": eventId,
			"message":  utils.GetErrString(err),
			"status":   "failed",
		}

		utils.Error(w, response, http.StatusBadGateway)
		return
	}

	// Assuming the transaction is successfully created and sent to the gRPC server
	response := RespW{
		"event_id": eventId,
		"message":  "Your transaction has been initiated and is being processed",
		"status":   "success",
	}
	utils.Success(w, response)
}
