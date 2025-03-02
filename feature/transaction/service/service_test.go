package service

import (
	"fmt"
	"soltracker/feature/transaction"
	"testing"
)

func TestGetTransaction(t *testing.T) {
	type args struct {
		signature string
	}
	tests := []struct {
		name string
		args args
		want transaction.TransactionResponse
	}{
		{
			name: "Test Get Transaction",
			args: args{
				signature: "28wdNAME7vun2z21FbaTHJ7qAwbWXWtge6ZVANJJTwqJTpRVF3YcbZMAZHXUfyNKNX2itKeLThEeFJ9hZtMeLxtu",
			},
		},
	}
	for _, tt := range tests {
		fmt.Println(GetTransaction(tt.args.signature))
	}
}

func TestGetTransactionLogs(t *testing.T) {
	type args struct {
		transactionResponse transaction.TransactionResponse
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Test Get Transaction Logs",
			args: args{
				transactionResponse: GetTransaction("28wdNAME7vun2z21FbaTHJ7qAwbWXWtge6ZVANJJTwqJTpRVF3YcbZMAZHXUfyNKNX2itKeLThEeFJ9hZtMeLxtu"),
			},
		},
	}
	for _, tt := range tests {
		fmt.Println(GetTransactionLogs(tt.args.transactionResponse))
	}
}

func TestGetWalletTransfer(t *testing.T) {
	type args struct {
		transactionResponse transaction.TransactionResponse
		walletAddress       string
	}
	tests := []struct {
		name    string
		args    args
		want    []transaction.TransactionTransfer
		wantErr bool
	}{
		{
			name: "Test Get Wallet Transfer",
			args: args{
				transactionResponse: GetTransaction("47Ms157WxD9r3Ra3NB7KYmESuyVm6VeyXNKAnMcy8miA4zJZWNavEEzx4ToPob3GRRgqSeWxRKfcscxnT88ZpPLo"),
				walletAddress:       "1FFHHz6vuqUFyu5iCveqhqD3mB3BrWqkngpRoBdkqcU",
			},
		},
	}
	for _, tt := range tests {
		transfers := GetWalletTransfer(tt.args.transactionResponse, tt.args.walletAddress)
		if len(transfers) == 0 {
			t.Errorf("GetWalletTransfer() error = %v", "No transfer found")
		}
		for _, transfer := range transfers {
			fmt.Println(transfer)
		}
	}
}

func TestGetSPLTokenTransaction(t *testing.T) {
	type args struct {
		transactionResponse transaction.TransactionResponse
		walletAddress       string
	}
	tests := []struct {
		name string
		args args
		want *transaction.SPLTokenTransfer
	}{
		{
			name: "Test Get SPL Token Transaction",
			args: args{
				transactionResponse: GetTransaction("2tbE3Tz9xkPFgotNghbgWqTHZvu3AxPqgGb4USxjkWTPtpexHdL4SErJxJ2sPvqxg9qxnj34vpkzFBU8mZPfFQsm"),
				walletAddress:       "8bM5Ez4vTC1LAKENSL3rp5WAfm7EfbpMUvWYQH2hD5rF",
			},
		},
	}
	for _, tt := range tests {
		fmt.Println(GetSPLTokenTransaction(tt.args.transactionResponse, tt.args.walletAddress))
	}
}
