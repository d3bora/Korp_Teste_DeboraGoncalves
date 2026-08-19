# Dockerfile do billing-service (multi-stage build)
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Permite usar toolchain mais recente se necessário
ENV GOTOOLCHAIN=auto

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o billing-service ./cmd/main.go

# ─── Imagem final mínima ───────────────────────────────────────────────────
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/billing-service .

EXPOSE 8082

CMD ["./billing-service"]
