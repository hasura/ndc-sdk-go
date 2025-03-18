package connector

import (
	"errors"
	"os"

	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

var (
	errAuthWebhookUriRequired = errors.New("the env var HASURA_CREDENTIALS_PROVIDER_URI must be set and non-empty")
)

// AccquireCredentials calls the credentials provider webhook to get the credentials for the given key.
// If force_refresh is true, the provider will ignore any cached credentials and fetch new ones.
// The credentials provider URI is read from the HASURA_CREDENTIALS_PROVIDER_URI environment variable.
// If the HASURA_CREDENTIALS_PROVIDER_BEARER_TOKEN environment variable is set, it will be used as a bearer token for the request.
func AccquireCredentials(key string, force_refresh bool) (string, error) {
	authWebhookUri, authWebhookUriExists := os.LookupEnv("HASURA_CREDENTIALS_PROVIDER_URI")
	if !authWebhookUriExists || authWebhookUri == "" {
		return "", errAuthWebhookUriRequired
	}
	bearerToken, bearerTokenExists := os.LookupEnv("HASURA_CREDENTIALS_PROVIDER_BEARER_TOKEN")

	params := url.Values{}
	params.Add("key", key)
	params.Add("force_refresh", strconv.FormatBool(force_refresh))

	fullWebhookUri := fmt.Sprintf("%s?%s", authWebhookUri, params.Encode())

	req, err := http.NewRequest("GET", fullWebhookUri, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	if bearerTokenExists {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	return string(body), nil
}
