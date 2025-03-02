package service

import (
	"fmt"
	"soltracker/feature/token"
	"testing"
)

func TestGetTokenInfo(t *testing.T) {
	type args struct {
		mint string
	}
	tests := []struct {
		name string
		args args
		want token.TokenInfoResponse
	}{
		{
			name: "Test Get Token Info",
			args: args{
				mint: "So11111111111111111111111111111111111111112",
			},
		},
	}
	for _, tt := range tests {
		fmt.Println(GetTokenInfo(tt.args.mint))
		fmt.Println(GetTokenName(GetTokenInfo(tt.args.mint)))
	}
}
