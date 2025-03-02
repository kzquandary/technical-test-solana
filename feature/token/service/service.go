package service

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"soltracker/config"
	"soltracker/feature/token"
)

func GetTokenInfo(mint string) token.TokenInfoResponse {
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "test",
		"method":  "getAsset",
		"params": map[string]string{
			"id": mint,
		},
	}

	rpcURL := config.GetRPCURL()

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Error marshalling payload: %v", err)
	}

	// send request
	resp, err := http.Post(rpcURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	// read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// parse response
	var tokenInfoResponse token.TokenInfoResponse
	err = json.Unmarshal(body, &tokenInfoResponse)
	if err != nil {
		log.Fatalf("Error unmarshalling response: %v", err)
	}

	// show result
	// tokenInfoJSON, err := json.MarshalIndent(tokenInfoResponse, "", "  ")
	// if err != nil {
	// 	log.Fatalf("Error marshalling response JSON: %v", err)
	// }

	// fmt.Printf("Token Info Response:\n%s\n", string(tokenInfoJSON))

	return tokenInfoResponse
}

func GetTokenName(tokenInfoResponse token.TokenInfoResponse) string {
	return tokenInfoResponse.Result.Content.Metadata.Symbol
}
