package provider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type User struct {
	Username string `json:"username"`
	ID       string `json:"id"`
}

type Client struct {
	ClientId     string
	ClientSecret string
	HTTPClient   *http.Client
	BaseURL      string
	Realm        string
}

func NewClient(clientId string, clientSecret string, baseURL string, realm string) *Client {
	return &Client{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		HTTPClient:   &http.Client{},
		BaseURL:      baseURL,
		Realm:        realm,
	}
}

func (c *Client) GetToken() (string, error) {
	data := url.Values{}
	data.Set("client_id", c.ClientId)
	data.Set("grant_type", "client_credentials")
	data.Set("scope", "openid")
	data.Set("client_secret", c.ClientSecret)

	resp, err := c.HTTPClient.PostForm(c.BaseURL+"/auth/realms/"+c.Realm+"/protocol/openid-connect/token", data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var TokenResult TokenResponse
	_ = json.Unmarshal(body, &TokenResult)

	return TokenResult.AccessToken, nil
}

func (c *Client) GetUserId(username string) (string, error) {
	token, err := c.GetToken()
	if err != nil {
		return "", err
	}

	req, _ := http.NewRequest("GET", c.BaseURL+"/auth/admin/realms/"+c.Realm+"/users?username="+username, nil)
	bearer := "Bearer " + token
	req.Header.Add("Authorization", bearer)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var users []User
	err = json.Unmarshal(body, &users)
	if err != nil {
		return "", err
	}

	if len(users) != 0 {
		// Assuming we only care about the first user
		return users[0].ID, nil
	}

	return "", fmt.Errorf("user not found")
}
