package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"soltracker/config"
	transactionService "soltracker/feature/transaction/service"

	"github.com/gagliardetto/solana-go"
	"github.com/gorilla/websocket"
)

const API_KEY = "0f803376-0189-4d72-95f6-a5f41cef157d"

const (
	wssURL = "wss://mainnet.helius-rpc.com/?api-key=" + API_KEY
	rpcURL = "https://mainnet.helius-rpc.com/?api-key=" + API_KEY
)

type LogSubscribePayload struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

type LogsNotification struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  struct {
		Result struct {
			Context struct {
				Slot uint64 `json:"slot"`
			} `json:"context"`
			Value struct {
				Signature string   `json:"signature"`
				Err       *string  `json:"err,omitempty"`
				Logs      []string `json:"logs"`
			} `json:"value"`
		} `json:"result"`
		Subscription int `json:"subscription"`
	} `json:"params"`
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please enter the Solana wallet address you'd like to track: ")
	walletAddress, _ := reader.ReadString('\n')
	walletAddress = strings.TrimSpace(walletAddress)

	if _, err := solana.PublicKeyFromBase58(walletAddress); err != nil {
		log.Fatal("Invalid Solana wallet address!")
	}

	fmt.Println("Wallet validated! - Monitoring Transactions...")

	conn, _, err := websocket.DefaultDialer.Dial(wssURL, nil)
	if err != nil {
		log.Fatalf("Failed to connect to Solana WebSocket: %v", err)
	}
	defer conn.Close()

	subscription := fmt.Sprintf(`{
		"jsonrpc": "2.0",
		"id": 1,
		"method": "logsSubscribe",
		"params": [
			{
				"mentions": [
					"%s"
				]
			},
			{
				"commitment": "finalized"
			}
		]
	}`, walletAddress)

	err = conn.WriteMessage(websocket.TextMessage, []byte(subscription))
	if err != nil {
		log.Fatalf("Failed to subscribe to account updates: %v", err)
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		var notification LogsNotification
		if err := json.Unmarshal(message, &notification); err != nil {
			// log.Printf("Error parsing message: %v", err)
			continue
		}
		if notification.Params.Result.Value.Err != nil {
			continue
		}
		if notification.Params.Result.Value.Logs == nil {
			continue
		}
		// log.Printf("Notifikasi: %v", notification)
		transactionResponse := transactionService.GetTransaction(notification.Params.Result.Value.Signature)

		splTokenTransfer := transactionService.GetSPLTokenTransaction(transactionResponse, walletAddress)
		if splTokenTransfer.TokenName != "" {
			fmt.Println(config.GetSPLTokenTransactionMessage(splTokenTransfer))
		}

		transferData := transactionService.GetWalletTransfer(transactionResponse, walletAddress)
		if len(transferData) > 0 {
			fmt.Println(config.GetTransferMessage(transferData[0]))
		}
	}
}
