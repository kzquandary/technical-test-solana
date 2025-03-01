package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gorilla/websocket"
)

const (
	wssURL = "wss://mainnet.helius-rpc.com/?api-key=API_KEY"
)

// AccountNotification for account updates
type AccountNotification struct {
	Params struct {
		Result struct {
			Context struct {
				Slot uint64 `json:"slot"`
			} `json:"context"`
			Value struct {
				Lamports uint64 `json:"lamports"`
			} `json:"value"`
		} `json:"result"`
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
		"method": "accountSubscribe",
		"params": [
			"%s",
			{"encoding": "jsonParsed", "commitment": "finalized"}
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
			time.Sleep(5 * time.Second)
			continue
		}

		var notification AccountNotification
		if err := json.Unmarshal(message, &notification); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		processAccountChange(walletAddress, notification)
	}
}

func processAccountChange(walletAddress string, notification AccountNotification) {
	fmt.Printf("Balance updated for %s: %d lamports (%.6f SOL) at slot %d\n",
		walletAddress,
		notification.Params.Result.Value.Lamports,
		float64(notification.Params.Result.Value.Lamports)/1e9,
		notification.Params.Result.Context.Slot)
}
