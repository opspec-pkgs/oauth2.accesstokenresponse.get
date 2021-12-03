package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/cli/oauth/device"
)

func main() {
	clientID := os.Getenv("CLIENT_ID")

	scopes := []string{}
	err := json.Unmarshal([]byte(os.Getenv("SCOPES")), &scopes)
	if err != nil {
		panic(err)
	}

	tokenEndpoint := os.Getenv("TOKEN_ENDPOINT")
	deviceAuthorizationEndpoint := os.Getenv("DEVICE_AUTHORIZATION_ENDPOINT")
	httpClient := http.DefaultClient

	code, err := device.RequestCode(httpClient, deviceAuthorizationEndpoint, clientID, scopes)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Copy code: %s\n", code.UserCode)
	fmt.Printf("then open: %s\n", fmt.Sprintf("%s?user_code=%s", code.VerificationURI, code.UserCode))

	accessTokenResponse, err := device.PollToken(httpClient, tokenEndpoint, clientID, code)
	if err != nil {
		panic(err)
	}

	accessTokenResponseBytes, err := json.Marshal(map[string]string{
		"refreshToken": accessTokenResponse.RefreshToken,
		"scope":        accessTokenResponse.Scope,
		"token":        accessTokenResponse.Token,
		"type":         accessTokenResponse.Type,
	})
	if err != nil {
		panic(err)
	}

	os.WriteFile("/accessTokenResponse", accessTokenResponseBytes, 0666)
}
