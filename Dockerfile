# 1. IMAGEM BASE: Usa uma versão leve do Golang baseada em Alpine Linux
FROM golang:1.25-alpine

# 2. DIRETÓRIO DE TRABALHO: Cria e define a pasta /app dentro do container
WORKDIR /app

# 3. GERENCIADOR DE PACOTES: Copia os arquivos go.mod e go.sum primeiro
COPY go.mod go.sum ./

# 4. DEPENDÊNCIAS: Baixa todas as bibliotecas necessárias para o projeto
RUN go mod download

# 5. CÓPIA DO PROJETO: Copia todo o restante do código da sua máquina para o container
COPY . .

# 6. COMPILAÇÃO: Compila o projeto Golang e gera um executável chamado 'api-audesp'
RUN go build -o api-audesp main.go

# 7. PORTA: Avisa ao Docker que o nosso servidor escuta na porta 8080
EXPOSE 8080

# 8. EXECUÇÃO: O comando que o container vai rodar quando for iniciado
CMD ["./api-audesp"]
