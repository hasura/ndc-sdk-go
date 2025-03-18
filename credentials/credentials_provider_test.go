package credentials

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestAcquireCredentials(t *testing.T) {
	t.Run("when HASURA_CREDENTIALS_PROVIDER_URI is not set", func(t *testing.T) {
		// Unset the environment variable
		os.Unsetenv("HASURA_CREDENTIALS_PROVIDER_URI")

		_, err := AcquireCredentials("key", false)
		if err != errAuthWebhookUriRequired {
			t.Errorf("expected error %v, got %v", errAuthWebhookUriRequired, err)
		}
	})

	t.Run("when HASURA_CREDENTIALS_PROVIDER_URI is set", func(t *testing.T) {
		t.Run("when the request is successful", func(t *testing.T) {
			bearerToken := "random-token"
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Query().Get("key") != "key" {
					t.Errorf("expected key=key; got %s", r.URL.Query().Get("key"))
				}

				if r.URL.Query().Get("force_refresh") != "false" {
					t.Errorf("expected force_refresh=false; got %s", r.URL.Query().Get("force_refresh"))
				}

				if r.Header.Get("Authorization") != ("Bearer " + bearerToken) {
					t.Errorf("expected Authorization=Bearer %s; got %s", bearerToken, r.Header.Get("Authorization"))
				}

				fmt.Fprint(w, "credentials")
			}))

			defer server.Close()

			// Set the environment variable
			os.Setenv("HASURA_CREDENTIALS_PROVIDER_URI", server.URL)
			os.Setenv("HASURA_CREDENTIALS_PROVIDER_BEARER_TOKEN", bearerToken)

			credentials, err := AcquireCredentials("key", false)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if credentials != "credentials" {
				t.Errorf("expected credentials to be 'credentials', got '%s'", credentials)
			}
		})

		t.Run("when the request fails", func(t *testing.T) {
			os.Setenv("HASURA_CREDENTIALS_PROVIDER_URI", "http://localhost:0000")

			_, err := AcquireCredentials("key", false)
			if err == nil {
				t.Error("expected an error, got nil")
			}
		})
	})
}
