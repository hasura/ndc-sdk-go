package connector


import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestAccquireCredentials(t *testing.T) {
	t.Run("when HASURA_CREDENTIALS_PROVIDER_URI is not set", func(t *testing.T) {
		// Unset the environment variable
		os.Unsetenv("HASURA_CREDENTIALS_PROVIDER_URI")

		_, err := AccquireCredentials("key", false)
		if err != errAuthWebhookUriRequired {
			t.Errorf("expected error %v, got %v", errAuthWebhookUriRequired, err)
		}
	})

	t.Run("when HASURA_CREDENTIALS_PROVIDER_URI is set", func(t *testing.T) {
		// Set the environment variable
		os.Setenv("HASURA_CREDENTIALS_PROVIDER_URI", "http://localhost:8080")

		t.Run("when the request is successful", func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, "credentials")
			}))
			defer server.Close()

			// Set the environment variable
			os.Setenv("HASURA_CREDENTIALS_PROVIDER_URI", server.URL)

			credentials, err := AccquireCredentials("key", false)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if credentials != "credentials" {
				t.Errorf("expected credentials to be 'credentials', got '%s'", credentials)
			}
		})

		t.Run("when the request fails", func(t *testing.T) {
			os.Setenv("HASURA_CREDENTIALS_PROVIDER_URI", "http://localhost:0000")

			_, err := AccquireCredentials("key", false)
			if err == nil {
				t.Error("expected an error, got nil")
			}
		})
	})
}