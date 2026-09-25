# --------------------------------------------------------
# Estágio 1: Base de Desenvolvimento e Testes
# --------------------------------------------------------
FROM golang:1.25.2-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["go", "test", "./...", "-v"]


# --------------------------------------------------------
# Estágio 2: Compilação do Executável
# --------------------------------------------------------
FROM builder AS compiler

# Compila o binário de forma estática para /app/cloudrun-observability
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/cloudrun-observability ./cmd/main.go


# --------------------------------------------------------
# Estágio 3: Imagem final de Execução (Produção com Scratch)
# --------------------------------------------------------
# Mantemos sem o "AS runner" já que o docker-compose não pede mais essa etapa por nome
FROM --platform=linux/amd64 scratch

# Copia os certificados necessários para conexões HTTPS externas
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app

# Copia o binário exatamente de onde foi gerado no estágio anterior
COPY --from=compiler /app/cloudrun-observability .

EXPOSE 8080

# Inicia o servidor usando o binário correto
ENTRYPOINT ["./cloudrun-observability"]
