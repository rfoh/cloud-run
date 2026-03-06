# Multi-stage build
FROM golang:1.25.4-alpine AS builder

WORKDIR /app

# Copiar go.mod e go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copiar código
COPY . .

# Build da aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

# Stage final - runtime
FROM alpine:3.19

RUN apk add --no-cache ca-certificates

WORKDIR /root/

# Copiar binary do builder
COPY --from=builder /app/server .

EXPOSE 8080

ENV PORT=8080

CMD ["./server"]
