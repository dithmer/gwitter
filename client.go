package gwitter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Endpoint string

const (
	URL string = "https://api.twitter.com"

	GetBearerToken Endpoint = "/oauth2/token" //nolint: gosec
)

type Client struct {
	httpClient *http.Client
	Token      string
}

func NewClient(httpClient *http.Client, apiKey string, apiKeySecret string) (*Client, error) {
	client := &Client{
		httpClient: httpClient,
	}

	err := client.authenticate(apiKey, apiKeySecret)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate: %w", err)
	}

	return client, nil
}

func (c *Client) authenticate(apiKey string, apiKeySecret string) error {
	values := url.Values{"grant_type": {"client_credentials"}}

	req, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodPost,
		buildURL(GetBearerToken),
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(apiKey, apiKeySecret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	type Response struct {
		// This comes from twitter api so I can't change it
		//nolint: tagliatelle
		Token string `json:"access_token"`
	}

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	c.Token = response.Token

	return nil
}

func buildURL(endpoint Endpoint) string {
	return fmt.Sprintf("%s%s", URL, endpoint)
}
