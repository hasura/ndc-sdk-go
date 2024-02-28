# build context at repo root: docker build -f Dockerfile .
FROM golang:1.22 AS builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -v -o ndc-cli .

# stage 2: production image
FROM gcr.io/distroless/base-debian12:nonroot

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/ndc-cli /ndc-cli

# Run the web service on container startup.
CMD ["/ndc-cli", "serve"]