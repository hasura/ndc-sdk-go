package credentials

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var errAuthWebhookUriRequired = errors.New("the env var HASURA_CREDENTIALS_PROVIDER_URI must be set and non-empty")

var defaultClient = CredentialClient{
	httpClient: http.DefaultClient,
}

// AcquireCredentials calls the credentials provider webhook to get the credentials for the given key.
// If force_refresh is true, the provider will ignore any cached credentials and fetch new ones.
// The credentials provider URI is read from the HASURA_CREDENTIALS_PROVIDER_URI environment variable.
// If the HASURA_CREDENTIALS_PROVIDER_BEARER_TOKEN environment variable is set, it will be used as a bearer token for the request.
func AcquireCredentials(key string, forceRefresh bool) (string, error) {
	return defaultClient.AcquireCredentials(key, forceRefresh)
}

// CredentialClient is an HTTP client that  can requests the credentials provider webhook to get the credentials.
type CredentialClient struct {
	providerUri         string
	providerBearerToken string
	httpClient          *http.Client
}

// NewCredentialClient creates a CredentialClient instance.
func NewCredentialClient(httpClient *http.Client) (*CredentialClient, error) {
	client := &CredentialClient{
		httpClient: httpClient,
	}

	return client, client.reload()
}

func (cc *CredentialClient) reload() error {
	providerUri, providerUriExists := os.LookupEnv("HASURA_CREDENTIALS_PROVIDER_URI")
	if !providerUriExists || providerUri == "" {
		return errAuthWebhookUriRequired
	}

	cc.providerUri = providerUri
	cc.providerBearerToken = os.Getenv("HASURA_CREDENTIALS_PROVIDER_BEARER_TOKEN")

	return nil
}

// AcquireCredentials calls the credentials provider webhook to get the credentials for the given key.
func (cc *CredentialClient) AcquireCredentials(key string, forceRefresh bool) (string, error) {
	if forceRefresh || cc.providerUri == "" {
		if err := cc.reload(); err != nil {
			return "", err
		}
	}

	params := url.Values{}
	params.Add("key", key)
	params.Add("force_refresh", strconv.FormatBool(forceRefresh))

	fullWebhookUri := fmt.Sprintf("%s?%s", cc.providerUri, params.Encode())

	req, err := http.NewRequest(http.MethodGet, fullWebhookUri, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	if cc.providerBearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+cc.providerBearerToken)
	}

	resp, err := cc.httpClient.Do(req)
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
