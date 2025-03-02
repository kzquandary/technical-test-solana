package config

import (
	"fmt"
	"soltracker/feature/transaction"
)

func GetTransferMessage(transferData transaction.TransactionTransfer) string {
	message := fmt.Sprintf("%s has sent %s SOL to %s", transferData.From, transferData.Amount, transferData.To)
	return message
}

func GetSPLTokenTransactionMessage(splTokenTransfer transaction.SPLTokenTransfer) string {
	message := fmt.Sprintf("%s bought %s with %s SOL", splTokenTransfer.WalletAddress, splTokenTransfer.TokenName, splTokenTransfer.SolAmount)
	return message
}
