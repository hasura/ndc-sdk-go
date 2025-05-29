package connector

import (
	"fmt"
	"net/http"

	"github.com/Masterminds/semver/v3"
	"github.com/hasura/ndc-sdk-go/schema"
	"go.opentelemetry.io/otel/metric"
)

// validate the requested NDC version from the engine. Skip the validation if the header is empty.
func (s *Server[Configuration, State]) withNDCVersionCheck(w http.ResponseWriter, r *http.Request) {
	requestedVersion := r.Header.Get(schema.XHasuraNDCVersion)
	if requestedVersion != "" {
		logger := GetLogger(r.Context())

		parsedVersion, err := semver.NewVersion(requestedVersion)
		if err != nil {
			writeError(
				w,
				logger,
				schema.BadRequestError(
					fmt.Sprintf(
						"Invalid %s header, expected a semver version string, got: %s",
						schema.XHasuraNDCVersion,
						requestedVersion,
					),
					nil,
				),
				http.StatusUnprocessableEntity,
			)
			s.increaseFailureCounterMetric(r, http.StatusBadRequest)

			return
		}

		// The connector may implement the patch version ahead of the engine
		// so the connector should validate major and minor version only.
		if s.ndcVersionConstraint.Major() != parsedVersion.Major() || s.ndcVersionConstraint.Minor() != parsedVersion.Minor() {
			writeError(
				w,
				logger,
				schema.BadRequestError(
					fmt.Sprintf(
						"NDC version range ^%s does not match implemented version %s",
						schema.NDCVersion,
						requestedVersion,
					),
					nil,
				),
				http.StatusUnprocessableEntity,
			)
			s.increaseFailureCounterMetric(r, http.StatusBadRequest)

			return
		}
	}
}

func (s *Server[Configuration, State]) withAuth(w http.ResponseWriter, r *http.Request) {
	logger := GetLogger(r.Context())
	// authorize the secret token in the request header if exists.
	if s.options.ServiceTokenSecret != "" &&
		r.Header.Get("Authorization") != ("Bearer "+s.options.ServiceTokenSecret) {
		writeJson(w, logger, http.StatusUnauthorized, schema.ErrorResponse{
			Message: "Unauthorized",
			Details: map[string]any{
				"cause": "Bearer token does not match.",
			},
		})

		s.increaseFailureCounterMetric(r, http.StatusUnauthorized)

		return
	}
}

func (s *Server[Configuration, State]) increaseFailureCounterMetric(
	r *http.Request,
	statusCode int,
) {
	attributes := metric.WithAttributes(failureStatusAttribute, httpStatusAttribute(statusCode))

	switch r.URL.Path {
	case apiPathQuery:
		s.telemetry.queryCounter.Add(r.Context(), 1, attributes)
	case apiPathQueryExplain:
		s.telemetry.queryExplainCounter.Add(r.Context(), 1, attributes)
	case apiPathMutation:
		s.telemetry.mutationCounter.Add(r.Context(), 1, attributes)
	case apiPathMutationExplain:
		s.telemetry.mutationExplainCounter.Add(r.Context(), 1, attributes)
	default:
	}
}
