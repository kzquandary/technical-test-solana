package config

func GetRPCURL() string {
	rpcURL := "https://mainnet.helius-rpc.com/?api-key="
	return rpcURL
}

func GetWSSURL() string {
	wssURL := "wss://mainnet.helius-rpc.com/?api-key="
	return wssURL
}
