# Use a imagem base oficial do Golang
FROM golang:1.23.1

# Defina o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copie o arquivo cliente.go para o diretório de trabalho
COPY cliente.go .

# Compile o aplicativo Go
RUN go build -o cliente cliente.go

# Comando para executar o aplicativo
CMD ["./cliente"]