# Dockerfile do stock-service
# Build multi-stage: compila o Go em uma imagem grande e copia apenas o binário para a imagem final.
# Resultado: imagem de produção pequena (~20MB) sem toolchain do Go.

# ─── Estágio 1: Build ──────────────────────────────────────────────────────
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copia e baixa dependências antes do código (cache de layers Docker)
COPY go.mod go.sum ./
RUN go mod download

# Copia o código-fonte e compila o binário
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o stock-service ./cmd/main.go

# ─── Estágio 2: Imagem final ───────────────────────────────────────────────
FROM alpine:3.19

WORKDIR /app

# Copia apenas o binário compilado (sem código-fonte ou toolchain)
COPY --from=builder /app/stock-service .

# Porta exposta pelo serviço
EXPOSE 8081

# Comando de inicialização
CMD ["./stock-service"]
