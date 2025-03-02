package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"soltracker/config"
	tokenService "soltracker/feature/token/service"
	"soltracker/feature/transaction"
)

func GetTransaction(signature string) transaction.TransactionResponse {

	rpcURL := config.GetRPCURL()

	// Constructing the request payload
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "getTransaction",
		"params":  []interface{}{signature, map[string]int{"maxSupportedTransactionVersion": 0}},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Error marshalling payload: %v", err)
	}

	// Sending the request
	resp, err := http.Post(rpcURL, "application/json", bytes.NewBuffer(payloadBytes))

	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	defer resp.Body.Close()

	// Reading the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	var txResponse transaction.TransactionResponse
	err = json.Unmarshal(body, &txResponse)
	if err != nil {
		log.Fatalf("Error unmarshalling response: %v", err)
	}

	return txResponse

}

func GetTransactionLogs(transactionResponse transaction.TransactionResponse) []string {

	logs := transactionResponse.Result.Meta.LogMessages

	return logs

}

func GetWalletTransfer(transactionResponse transaction.TransactionResponse, walletAddress string) []transaction.TransactionTransfer {
	// list to save all transfers
	var transfers []transaction.TransactionTransfer

	// get account keys from transaction
	accountKeys := transactionResponse.Result.Transaction.Message.AccountKeys

	// save pre and post balance
	preBalances := transaction.ConvertBalances(transactionResponse.Result.Meta.PreBalances)
	postBalances := transaction.ConvertBalances(transactionResponse.Result.Meta.PostBalances)

	// Iterasi melalui setiap instruksi dalam transaksi
	for _, instruction := range transactionResponse.Result.Transaction.Message.Instructions {
		// instruction with two accounts (sender and receiver)
		if len(instruction.Accounts) == 2 {
			fromIndex := instruction.Accounts[0] // sender
			toIndex := instruction.Accounts[1]   // receiver

			// get address of the account involved
			from := accountKeys[fromIndex]
			to := accountKeys[toIndex]

			// check if wallet address is involved as sender or receiver
			if from == walletAddress || to == walletAddress {
				// get pre and post balance for sender and receiver
				preFromBalance := preBalances[fromIndex]
				preToBalance := preBalances[toIndex]
				postFromBalance := postBalances[fromIndex]
				postToBalance := postBalances[toIndex]

				// check who is sending and who is receiving
				var amount float64

				// if sender balance is decreased (pre is greater than post)
				if preFromBalance > postFromBalance {
					// sender
					amount = float64(preFromBalance-postFromBalance) / 1000000000 // convert lamports to SOL
					transfer := transaction.TransactionTransfer{
						Signature: transactionResponse.Result.Transaction.Signatures[0],
						From:      from,
						To:        to,
						Amount:    fmt.Sprintf("%.8f", amount), // format dengan 8 desimal
					}
					transfers = append(transfers, transfer)
				} else if preToBalance > postToBalance {
					// receiver
					amount = float64(postToBalance-preToBalance) / 1000000000 // konversi lamports ke SOL
					transfer := transaction.TransactionTransfer{
						Signature: transactionResponse.Result.Transaction.Signatures[0],
						From:      from,
						To:        to,
						Amount:    fmt.Sprintf("%.8f", amount),
					}
					transfers = append(transfers, transfer)
				}
			}
		}
	}

	// if there is no transfer, return nil
	if len(transfers) == 0 {
		return nil
	}

	return transfers
}

func GetSPLTokenTransaction(transactionResponse transaction.TransactionResponse, walletAddress string) transaction.SPLTokenTransfer {
	accountKeys := transactionResponse.Result.Transaction.Message.AccountKeys
	preBalances := transaction.ConvertBalances(transactionResponse.Result.Meta.PreBalances)
	postBalances := transaction.ConvertBalances(transactionResponse.Result.Meta.PostBalances)
	preTokenBalances := transactionResponse.Result.Meta.PreTokenBalances
	postTokenBalances := transactionResponse.Result.Meta.PostTokenBalances
	signature := transactionResponse.Result.Transaction.Signatures[0]

	// find index of wallet in accountKeys
	var walletIndex int = -1
	for i, key := range accountKeys {
		if key == walletAddress {
			walletIndex = i
			break
		}
	}

	// if wallet is not found in accountKeys, return nil
	if walletIndex == -1 {
		log.Println("Wallet address not found in account keys")
		return transaction.SPLTokenTransfer{}
	}

	// loop to find SPL transaction related to wallet
	for _, postToken := range postTokenBalances {
		// check if accountIndex of postToken is the same as walletIndex
		if postToken.Owner == walletAddress {
			// find preTokenBalance that is the same as postToken.AccountIndex
			var preAmount string = "0"
			for _, preToken := range preTokenBalances {
				if preToken.AccountIndex == postToken.AccountIndex {
					preAmount = preToken.UITokenAmount.Amount
					break
				}
			}

			// convert amount from string to int
			postAmountInt, _ := strconv.ParseInt(postToken.UITokenAmount.Amount, 10, 64)
			preAmountInt, _ := strconv.ParseInt(preAmount, 10, 64)

			// calculate token transferred (always positive)
			tokenTransferred := int64(math.Abs(float64(postAmountInt - preAmountInt)))

			// get token info based on mint
			tokenInfoResponse := tokenService.GetTokenInfo(postToken.Mint)
			tokenName := tokenService.GetTokenName(tokenInfoResponse)

			// calculate SOL that is subtracted from preBalance and postBalance of wallet (always positive)
			solChange := float64(preBalances[walletIndex]-postBalances[walletIndex]) / 1e9
			solAmount := fmt.Sprintf("%.9f", math.Abs(solChange))

			// convert tokenTransferred to string with the correct format
			tokenAmount := fmt.Sprintf("%.6f", float64(tokenTransferred)/float64(1e6)) // according to decimals 6

			// return SPL transaction result
			return transaction.SPLTokenTransfer{
				Signature:     signature,
				WalletAddress: walletAddress,
				TokenName:     tokenName,
				SolAmount:     solAmount,
				TokenAmount:   tokenAmount,
			}
		}
	}

	// if there is no SPL transaction, return nil
	return transaction.SPLTokenTransfer{}
}

func IsWalletTransfer(accountKeys []string) bool {

	sistemProgram := "11111111111111111111111111111111"

	for _, accountKey := range accountKeys {
		if strings.Contains(accountKey, sistemProgram) {
			return true
		}
	}

	return false
}
