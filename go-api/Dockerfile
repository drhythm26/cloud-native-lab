FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o api .


FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/api /api
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/api"]
