FROM golang:1.24-bookworm AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -o /out/release-tracker-api \
    ./cmd/api

FROM gcr.io/distroless/static-debian12:nonroot
ENV PORT=8080
WORKDIR /app
COPY --from=builder /out/release-tracker-api /app/release-tracker-api
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/release-tracker-api"]
