package credentials

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"go.opentelemetry.io/otel"
)

func TestAcquireCredentials(t *testing.T) {
	t.Run("when HASURA_CREDENTIALS_PROVIDER_URI is not set", func(t *testing.T) {
		// Unset the environment variable
		os.Unsetenv("HASURA_CREDENTIALS_PROVIDER_URI")

		_, err := AcquireCredentials(context.TODO(), "key", false)
		if err != ErrAuthWebhookUriRequired {
			t.Errorf("expected error %v, got %v", ErrAuthWebhookUriRequired, err)
		}
	})

	t.Run("when HASURA_CREDENTIALS_PROVIDER_URI is set", func(t *testing.T) {
		t.Run("when the request is successful", func(t *testing.T) {
			bearerToken := "random-token"
			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Query().Get("key") != "key" {
						t.Errorf("expected key=key; got %s", r.URL.Query().Get("key"))
					}

					if r.URL.Query().Get("force_refresh") != "false" {
						t.Errorf(
							"expected force_refresh=false; got %s",
							r.URL.Query().Get("force_refresh"),
						)
					}

					if r.Header.Get("Authorization") != ("Bearer " + bearerToken) {
						t.Errorf(
							"expected Authorization=Bearer %s; got %s",
							bearerToken,
							r.Header.Get("Authorization"),
						)
					}

					fmt.Fprint(w, "{ \"credentials\": \"api-key\" }")
				}),
			)

			defer server.Close()

			// Set the environment variable
			os.Setenv("HASURA_CREDENTIALS_PROVIDER_URI", server.URL)
			os.Setenv("HASURA_CREDENTIALS_PROVIDER_BEARER_TOKEN", bearerToken)
			defer os.Unsetenv("HASURA_CREDENTIALS_PROVIDER_URI")
			defer os.Unsetenv("HASURA_CREDENTIALS_PROVIDER_BEARER_TOKEN")

			credentials, err := AcquireCredentials(context.TODO(), "key", false)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if credentials != "api-key" {
				t.Errorf("expected credentials to be 'api-key', got '%s'", credentials)
			}
		})

		t.Run("when the server does not exist", func(t *testing.T) {
			os.Setenv("HASURA_CREDENTIALS_PROVIDER_URI", "http://localhost:0000")

			_, err := AcquireCredentials(context.TODO(), "key", false)
			if err == nil {
				t.Error("expected an error, got nil")
			}
		})
	})

	t.Run("when the response does not have credentials", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "{}")
		}))

		defer server.Close()

		serverUri, err := url.Parse(server.URL)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		client := &CredentialClient{
			providerUri: serverUri,
			httpClient:  server.Client(),
			propagator:  otel.GetTextMapPropagator(),
		}

		_, err = client.AcquireCredentials(context.TODO(), "key", false)
		if err != ErrEmptyCredentials {
			t.Errorf("expected an empty credentails error, got: %s\n", err)
		}
	})
}
