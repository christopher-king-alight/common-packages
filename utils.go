package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type SecretResponse struct {
	ClientId     string `json:"WIZ_CLIENT_ID"`
	ClientSecret string `json:"WIZ_CLIENT_SECRET"`
	Endpoint     string `json:"WIZ_ENDPOINT"`
}
type AccessToken struct {
	Token string `json:"access_token"`
}

type ResponseData struct {
	Data interface{} `json:"data"`
}

func LoginToWiz(clientId, clientSecret, endpoint string) (AccessToken, error) {

	auth_data := url.Values{}

	auth_data.Set("grant_type", "client_credentials")
	auth_data.Set("audience", "wiz-api")
	auth_data.Set("client_id", clientId)
	auth_data.Set("client_secret", clientSecret)

	client := &http.Client{}

	r, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(auth_data.Encode()))
	if err != nil {
		log.Fatal("Error creating request in login: ", err)
		return AccessToken{}, err
	}
	r.Header.Add("Encoding", "UTF-8")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	fmt.Println("Getting token.")
	resp, err := client.Do(r)
	if err != nil {
		log.Fatal("Error getting token: ", err)
		return AccessToken{}, err
	}
	fmt.Println("Authentication response: " + resp.Status)

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error authenticating: %v", fmt.Errorf("status code: %d", resp.StatusCode))
		return AccessToken{}, err
	}
	bodyBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("Failed reading response body: %v", err)
		return AccessToken{}, err
	}

	at := AccessToken{}
	jsonErr := json.Unmarshal(bodyBytes, &at)

	if jsonErr != nil {
		log.Fatalf("Failed parsing JSON body: %v", jsonErr)
		return AccessToken{}, err
	}

	return at, nil
}
