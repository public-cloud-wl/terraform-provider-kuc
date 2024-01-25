package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
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
	Token        string
	mutex        sync.Mutex
}

func NewClient(clientId string, clientSecret string, baseURL string, realm string, httpClient *http.Client) *Client {
	return &Client{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		HTTPClient:   httpClient,
		BaseURL:      baseURL,
		Realm:        realm,
	}
}
func (c *Client) GetToken(ctx context.Context) (string, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.Token != "" {
		return c.Token, nil
	}

	data := url.Values{}
	data.Set("client_id", c.ClientId)
	data.Set("grant_type", "client_credentials")
	data.Set("scope", "openid")
	data.Set("client_secret", c.ClientSecret)

	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/auth/realms/"+c.Realm+"/protocol/openid-connect/token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		return "", err
	}

	var TokenResult TokenResponse
	err = json.Unmarshal(body, &TokenResult)
	if err != nil {
		return "", err
	}

	c.Token = TokenResult.AccessToken
	return c.Token, nil
}

func (c *Client) GetUserId(ctx context.Context, username string) (string, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", c.BaseURL+"/auth/admin/realms/"+c.Realm+"/users?username="+username, nil)
	if err != nil {
		return "", err
	}

	bearer := "Bearer " + token
	req.Header.Add("Authorization", bearer)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		return "", err
	}

	var users []User
	err = json.Unmarshal(body, &users)
	if err != nil {
		return "", err
	}

	if len(users) != 0 {
		return users[0].ID, nil
	}

	return "", fmt.Errorf("user not found")
}

func (c *Client) GetUser(ctx context.Context, id string) (User, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return User{}, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", c.BaseURL+"/auth/admin/realms/"+c.Realm+"/users/"+id, nil)
	if err != nil {
		return User{}, err
	}

	bearer := "Bearer " + token
	req.Header.Add("Authorization", bearer)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		return User{}, err
	}

	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
