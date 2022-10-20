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
	type Response struct {
		//nolint: tagliatelle
		Token string `json:"access_token"`
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		j, err := json.Marshal(Response{Token: "TEST_TOKEN"})
		if err != nil {
			t.Fatal(err)
		}
		fmt.Fprintln(w, string(j))
	}))

	hclient := server.Client()
	client, err := gwitter.NewClient(hclient, "TEST_API_KEY", "TEST_API_KEY_SECRET")
	if err != nil {
		t.Errorf("failed to create client: %v", err)
	}

	fmt.Println(client.Token)
	if client != nil && client.Token != "TEST_TOKEN" {
		t.Errorf("token is not set")
	}
}
