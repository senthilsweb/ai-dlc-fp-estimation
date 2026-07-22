# Stage 1: Build the Go binary
FROM golang:1.22-alpine AS builder
RUN apk add --no-cache git
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
COPY server/ server/
COPY app/ app/
COPY data/ data/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /fp-estimator .

# Stage 2: Minimal runtime
FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=builder /fp-estimator /usr/local/bin/fp-estimator

EXPOSE 8080
ENTRYPOINT ["fp-estimator"]
