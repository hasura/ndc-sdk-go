package credentials

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

var (
	ErrAuthWebhookUriRequired = errors.New(
		"the env var HASURA_CREDENTIALS_PROVIDER_URI must be set and non-empty",
	)
	ErrEmptyCredentials = errors.New("empty credentials")
)

var defaultClient = CredentialClient{
	httpClient: http.DefaultClient,
}

var tracer = otel.Tracer("CredentialProvider")

// AcquireCredentials calls the credentials provider webhook to get the credentials for the given key.
// If force_refresh is true, the provider will ignore any cached credentials and fetch new ones.
// The credentials provider URI is read from the HASURA_CREDENTIALS_PROVIDER_URI environment variable.
// If the HASURA_CREDENTIALS_PROVIDER_BEARER_TOKEN environment variable is set, it will be used as a bearer token for the request.
func AcquireCredentials(ctx context.Context, key string, forceRefresh bool) (string, error) {
	return defaultClient.AcquireCredentials(ctx, key, forceRefresh)
}

// Payload is the credentials provider webhook response payload.
type Payload struct {
	Credentials string `json:"credentials"`
}

// CredentialClient is an HTTP client that  can requests the credentials provider webhook to get the credentials.
type CredentialClient struct {
	providerUri         *url.URL
	providerBearerToken string
	httpClient          *http.Client
	propagator          propagation.TextMapPropagator
}

// NewCredentialClient creates a CredentialClient instance.
func NewCredentialClient(httpClient *http.Client) (*CredentialClient, error) {
	client := &CredentialClient{
		httpClient: httpClient,
	}

	return client, client.reload()
}

func (cc *CredentialClient) reload() error {
	rawProviderUri, providerUriExists := os.LookupEnv("HASURA_CREDENTIALS_PROVIDER_URI")
	if !providerUriExists || rawProviderUri == "" {
		return ErrAuthWebhookUriRequired
	}

	providerUri, err := url.Parse(rawProviderUri)
	if err != nil {
		return fmt.Errorf("invalid HASURA_CREDENTIALS_PROVIDER_URI: %w", err)
	}

	if providerUri.Scheme != "http" && providerUri.Scheme != "https" {
		return errors.New("invalid HASURA_CREDENTIALS_PROVIDER_URI: allow http(s) scheme only")
	}

	cc.providerUri = providerUri
	cc.providerBearerToken = os.Getenv("HASURA_CREDENTIALS_PROVIDER_BEARER_TOKEN")
	cc.propagator = otel.GetTextMapPropagator()

	return nil
}

// AcquireCredentials calls the credentials provider webhook to get the credentials for the given key.
func (cc *CredentialClient) AcquireCredentials(
	ctx context.Context,
	key string,
	forceRefresh bool,
) (string, error) {
	ctx, span := tracer.Start(
		ctx,
		"AcquireCredentials",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attribute.String("internal.visibility", "user")),
	)
	defer span.End()

	if forceRefresh || cc.providerUri == nil {
		if err := cc.reload(); err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)

			return "", err
		}
	}

	span.SetAttributes(
		attribute.String("http.request.method", http.MethodGet),
		attribute.String("url.full", cc.providerUri.String()),
		attribute.String("server.address", cc.providerUri.Hostname()),
		attribute.String("server.port", cc.providerUri.Port()),
		attribute.String("network.protocol.name", "http"),
		attribute.Bool("force_refresh", forceRefresh),
	)

	requestUri := *cc.providerUri
	requestQuery := requestUri.Query()
	requestQuery.Set("key", key)
	requestQuery.Set("force_refresh", strconv.FormatBool(forceRefresh))
	requestUri.RawQuery = requestQuery.Encode()
	fullWebhookUri := requestUri.String()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullWebhookUri, nil)
	if err != nil {
		span.SetStatus(codes.Error, "failed to create http request")
		span.RecordError(err)

		return "", fmt.Errorf("error creating request: %w", err)
	}

	if cc.providerBearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+cc.providerBearerToken)
	}

	cc.propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := cc.httpClient.Do(req)
	if err != nil {
		span.SetStatus(codes.Error, "failed to do the http request")
		span.RecordError(err)

		return "", fmt.Errorf("error making request: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	span.SetAttributes(attribute.Int("http.response.status_code", resp.StatusCode))

	var payload Payload

	err = json.NewDecoder(resp.Body).Decode(&payload)
	if err != nil {
		span.SetStatus(codes.Error, "failed to read the response")
		span.RecordError(err)

		return "", fmt.Errorf("error reading response: %w", err)
	}

	if payload.Credentials == "" {
		span.SetStatus(codes.Error, ErrEmptyCredentials.Error())
		span.RecordError(err)

		return "", ErrEmptyCredentials
	}

	return payload.Credentials, nil
}
