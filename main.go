package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// RequestToken response from Withings API
// https://developer.withings.com/oauth2/#operation/oauth2-getaccesstoken
type RequestToken struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
	Body   struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		Scope        string `json:"scope"`
		ExpiresIn    string `json:"expires_in"`
		TokenType    string `json:"token_type"`
	} `json:"body"`
}

func main() {
	const withingsAPIBaseURL = "https://wbsapi.withings.net"

	accessToken := os.Getenv("WITHINGS_API_ACCESS_TOKEN")
	if accessToken == "" {
		clientID := os.Getenv("WITHINGS_APP_CLIENT_ID")
		clientSecret := os.Getenv("WITHINGS_APP_CLIENT_SECRET")

		if clientID == "" || clientSecret == "" {
			fmt.Println("Set your Withings API application up with `WITHINGS_APP_CLIENT_ID` and `WITHINGS_APP_CLIENT_SECRET` envvars.")
			return
		}

		const scopes = "user.info,user.metrics"
		_, accessToken = oauthFlow(withingsAPIBaseURL, clientID, clientSecret, scopes)
	}
}

func oauthFlow(withingsAPIBaseURL string, clientID string, clientSecret string, scopes string) (string, string) {
	authCode := ""
	fmt.Printf("Go to https://account.withings.com/oauth2_user/authorize2?response_type=code&client_id=%s&scope=%s&state=issyl0-withings&redirect_uri=http://localhost\n", clientID, scopes)
	fmt.Println("Enter the value of `code` from the returned query string:")
	fmt.Scanln(&authCode)

	url := fmt.Sprintf("%s/v2/oauth2?action=requesttoken&grant_type=authorization_code&client_id=%s&client_secret=%s&code=%s&redirect_uri=http://localhost", withingsAPIBaseURL, clientID, clientSecret, authCode)
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	parsedRequestToken := RequestToken{}
	json.Unmarshal(body, &parsedRequestToken)

	accessToken := parsedRequestToken.Body.AccessToken
	fmt.Printf("To avoid reauthenticating every time, run `export WITHINGS_API_ACCESS_TOKEN=%s`\n", accessToken)
	return authCode, accessToken
}