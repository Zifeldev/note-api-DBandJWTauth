FROM golang:1.24.2-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download


COPY cmd/ ./cmd/
COPY internals/ ./internals/
COPY templates/ ./templates/

RUN go build -o app ./cmd

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/internals/migrations ./internals/migrations 

RUN apk add --no-cache ca-certificates

EXPOSE 8484
CMD ["./app"]