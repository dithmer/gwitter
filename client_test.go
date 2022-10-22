package gwitter_test

import (
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
