package gwitter_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dithmer/gwitter"
)

func TestNewClient(t *testing.T) {
	t.Parallel()
	type Response struct {
		//nolint: tagliatelle
		Token string `json:"access_token"`
	}
	t.Run("defining a new client should automatically authenticate", func(t *testing.T) {
		t.Parallel()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			j, err := json.Marshal(Response{Token: "TEST_TOKEN"})
			if err != nil {
				t.Fatal(err)
			}
			fmt.Fprintln(w, string(j))
		}))
		defer server.Close()

		client, err := gwitter.NewClient(server.URL, server.Client(), "TEST_API_KEY", "TEST_API_KEY_SECRET")
		if err != nil {
			t.Errorf("failed to create client: %v", err)
		}

		if client != nil && client.Token != "TEST_TOKEN" {
			t.Errorf("token is not set")
		}
	})

	t.Run("defining a new client should return an error if authentication fails", func(t *testing.T) {
		t.Parallel()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer server.Close()

		_, err := gwitter.NewClient(server.URL, server.Client(), "TEST_API_KEY", "TEST_API_KEY_SECRET")
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("defining a new client should call the correct endpoint", func(t *testing.T) {
		t.Parallel()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/oauth2/token" {
				t.Errorf("expected /oauth2/token, got %v", r.URL.Path)
			}
			fmt.Fprintln(w, "{}")
		}))
		defer server.Close()

		_, err := gwitter.NewClient(server.URL, server.Client(), "TEST_API_KEY", "TEST_API_KEY_SECRET")
		if err != nil {
			t.Errorf("failed to create client: %v", err)
		}
	})
}

func TestDoAuthenticatedRequest(t *testing.T) {
	t.Parallel()

	t.Run("should return an error if the request fails with unauthorized", func(t *testing.T) {
		t.Parallel()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer server.Close()

		mockClient := &gwitter.Client{
			HttpClient: server.Client(),
			URL:        server.URL,
			Token:      "TEST_TOKEN",
		}

		req, _ := http.NewRequestWithContext(context.TODO(), http.MethodGet, server.URL+"/auth/endpoint", nil)

		_, err := mockClient.DoAuthenticatedRequest(req) //nolint: bodyclose
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})

	t.Run("should return an error if the clients token is empty", func(t *testing.T) {
		t.Parallel()

		mockClient := &gwitter.Client{
			HttpClient: nil,
			URL:        gwitter.DefaultURL,
			Token:      "",
		}

		req, _ := http.NewRequestWithContext(context.TODO(), http.MethodGet, mockClient.URL+"/auth/endpoint", nil)

		_, err := mockClient.DoAuthenticatedRequest(req)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}
