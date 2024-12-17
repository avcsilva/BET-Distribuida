# Use a imagem base oficial do Golang
FROM golang:1.23.1

# Defina o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copie o arquivo go.mod e go.sum, se existirem
COPY go.mod go.sum ./

# Baixe as dependências
RUN go mod download

# Copie o arquivo servidor.go para o diretório de trabalho
COPY servidor.go .

# Compile o aplicativo Go
RUN go build -o servidor servidor.go

# Exponha a porta que o servidor irá usar
EXPOSE 8080

# Comando para executar o aplicativo
CMD ["./servidor"]