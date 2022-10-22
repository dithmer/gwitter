package gwitter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

type Endpoint string

const (
	DefaultURL = "https://api.twitter.com"

	GetBearerToken Endpoint = "/oauth2/token" //nolint: gosec
)

type Client struct {
	HttpClient *http.Client
	URL        string
	Token      string
}

func NewDefaultClient(apiKey string, apiKeySecret string) (*Client, error) {
	return NewClient(DefaultURL, http.DefaultClient, apiKey, apiKeySecret)
}

func NewClient(url string, hc *http.Client, apiKey string, apiKeySecret string) (*Client, error) {
	client := &Client{
		HttpClient: hc,
		Token:      "",
		URL:        url,
	}

	err := client.authenticate(apiKey, apiKeySecret)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate: %w", err)
	}

	return client, nil
}

func (c *Client) DoAuthenticatedRequest(r *http.Request) (*http.Response, error) {
	handleError := func(err error) (*http.Response, error) {
		return nil, errors.Wrap(err, "failed to do authenticated request")
	}

	if c.Token == "" {
		return handleError(fmt.Errorf("token is empty"))
	}

	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))

	resp, err := c.HttpClient.Do(r)
	if err != nil {
		return handleError(err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return handleError(errors.New(resp.Status))
	}

	return resp, nil
}

func (c *Client) authenticate(apiKey string, apiKeySecret string) error {
	values := url.Values{"grant_type": {"client_credentials"}}

	req, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodPost,
		buildURL(c.URL, GetBearerToken),
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(apiKey, apiKeySecret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	resp, err := c.HttpClient.Do(req)
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

func buildURL(baseURL string, endpoint Endpoint) string {
	return fmt.Sprintf("%s%s", baseURL, endpoint)
}
