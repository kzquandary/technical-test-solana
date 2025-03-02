package transaction

type TransactionResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		BlockTime int `json:"blockTime"`
		Meta      struct {
			ComputeUnitsConsumed int         `json:"computeUnitsConsumed"`
			Err                  interface{} `json:"err"`
			Fee                  int         `json:"fee"`
			InnerInstructions    []struct {
				Index        int `json:"index"`
				Instructions []struct {
					Accounts       []int  `json:"accounts"`
					Data           string `json:"data"`
					ProgramIDIndex int    `json:"programIdIndex"`
					StackHeight    int    `json:"stackHeight"`
				} `json:"instructions"`
			} `json:"innerInstructions"`
			LoadedAddresses struct {
				Readonly []string `json:"readonly"`
				Writable []string `json:"writable"`
			} `json:"loadedAddresses"`
			LogMessages       []string      `json:"logMessages"`
			PostBalances      []interface{} `json:"postBalances"`
			PostTokenBalances []struct {
				AccountIndex  int    `json:"accountIndex"`
				Mint          string `json:"mint"`
				Owner         string `json:"owner"`
				ProgramID     string `json:"programId"`
				UITokenAmount struct {
					Amount         string  `json:"amount"`
					Decimals       int     `json:"decimals"`
					UIAmount       float64 `json:"uiAmount"`
					UIAmountString string  `json:"uiAmountString"`
				} `json:"uiTokenAmount"`
			} `json:"postTokenBalances"`
			PreBalances      []interface{} `json:"preBalances"`
			PreTokenBalances []struct {
				AccountIndex  int    `json:"accountIndex"`
				Mint          string `json:"mint"`
				Owner         string `json:"owner"`
				ProgramID     string `json:"programId"`
				UITokenAmount struct {
					Amount         string  `json:"amount"`
					Decimals       int     `json:"decimals"`
					UIAmount       float64 `json:"uiAmount"`
					UIAmountString string  `json:"uiAmountString"`
				} `json:"uiTokenAmount"`
			} `json:"preTokenBalances"`
			Rewards []interface{} `json:"rewards"`
			Status  struct {
				Ok interface{} `json:"Ok"`
			} `json:"status"`
		} `json:"meta"`
		Slot        int `json:"slot"`
		Transaction struct {
			Message struct {
				AccountKeys         []string `json:"accountKeys"`
				AddressTableLookups []struct {
					AccountKey      string `json:"accountKey"`
					ReadonlyIndexes []int  `json:"readonlyIndexes"`
					WritableIndexes []int  `json:"writableIndexes"`
				} `json:"addressTableLookups"`
				Header struct {
					NumReadonlySignedAccounts   int `json:"numReadonlySignedAccounts"`
					NumReadonlyUnsignedAccounts int `json:"numReadonlyUnsignedAccounts"`
					NumRequiredSignatures       int `json:"numRequiredSignatures"`
				} `json:"header"`
				Instructions []struct {
					Accounts       []int       `json:"accounts"`
					Data           string      `json:"data"`
					ProgramIDIndex int         `json:"programIdIndex"`
					StackHeight    interface{} `json:"stackHeight"`
				} `json:"instructions"`
				RecentBlockhash string `json:"recentBlockhash"`
			} `json:"message"`
			Signatures []string `json:"signatures"`
		} `json:"transaction"`
		// Version string `json:"version"`
	} `json:"result"`
	ID int `json:"id"`
}
