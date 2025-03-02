package transaction

type TransactionTransfer struct {
	Signature string
	From      string
	To        string
	Amount    string
}

type SPLTokenTransfer struct {
	Signature     string
	WalletAddress string
	TokenName     string
	SolAmount     string
	TokenAmount   string
}
