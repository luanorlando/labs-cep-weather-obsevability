# 1. Etapa de Compilação (Builder) usando o Go 1.25 atualizado
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copia os arquivos de dependências primeiro (ajuda no cache do Docker)
COPY go.mod ./
# Se você tiver um go.sum na raiz, descomente a linha abaixo:
# COPY go.sum ./
RUN go mod download

# Copia o restante de todo o projeto (incluindo todas as pastas de serviço)
COPY . .

# Compila o binário estático apontando diretamente para o cmd do Serviço A
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./service-a/cmd/main.go

# 2. Etapa de Execução (Imagem Final super leve)
FROM alpine:latest
WORKDIR /app

# Copia o binário gerado na etapa anterior
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
