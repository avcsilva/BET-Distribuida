package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var id int = 0
var url string //Ex.: "http://localhost:8080"

type Evento struct {
	id                 int
	ativo              bool
	id_criador         int
	nome               string
	descricao          string
	participantes      map[int]float64 //id dos participantes e valor pago
	palpite            map[int]string  //id dos participantes e palpite
	porcentagemCriador float64
	resultado          string
}

type Cadastro_req struct {
	Id   int    `json:"id"`
	Nome string `json:"nome"`
}

type Cria_Evento_req struct {
	Id                 int     `json:"id"`
	Id_event           int     `json:"id_event"`
	Nome               string  `json:"nome"`
	Descricao          string  `json:"descricao"`
	PorcentagemCriador float64 `json:"porcentagemCriador"`
}

// Função para limpar o terminal
func limpar_terminal() {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")

	default: //linux e mac
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	erro := cmd.Run()
	if erro != nil {
		fmt.Println("Erro ao limpar o terminal:", erro)
		return
	}
}

// Função para exibir uma mensagem alternando cores
func displayMessageWithColors(message string, seconds int) {
	// Códigos de cores ANSI
	colors := []string{
		"\033[31m", // Vermelho
		"\033[32m", // Verde
		"\033[33m", // Amarelo
		"\033[34m", // Azul
		"\033[35m", // Magenta
		"\033[36m", // Ciano
	}

	// Define o tempo de duração
	start := time.Now()
	duration := time.Duration(seconds) * time.Second

	for {
		// Verifica se o tempo limite foi atingido
		if time.Since(start) > duration {
			break
		}

		for _, color := range colors {
			// Exibe a mensagem alternando as cores
			fmt.Printf("%s%s\033[0m\r", color, message) // \033[0m reseta as cores
			time.Sleep(500 * time.Millisecond)          // Delay de 500ms

			// Verifica novamente se o tempo foi atingido dentro do loop interno
			if time.Since(start) > duration {
				break
			}
		}
	}
}

// Função para exibir o cabeçalho com o endereço do servidor para conexão
func cabecalho(endereco string) {
	limpar_terminal()

	tamanho := len(endereco)
	espacamento := ""
	if tamanho < 33 {
		espacamento = strings.Repeat(" ", 33-tamanho)
	}
	fmt.Println("=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-")
	fmt.Println("|\033[32m                 BET - Distribuida         	 \033[0m|")
	fmt.Println("|--------------------------------------------------------|")
	fmt.Println("|\033[34m            Conectado:", endereco+espacamento+"\033[0m|")
	fmt.Print("=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-\n\n")
}

func cadastrar(nome string) bool {
	var cadastro = Cadastro_req{Nome: nome}   // ID é gerado pelo servidor
	json_valor, err := json.Marshal(cadastro) // Serializa o JSON
	if err != nil {
		fmt.Println("Erro ao serializar o JSON:", err)
		return false
	}
	fmt.Println("Serializado: sucesso")

	resposta, err := http.Post(url+"/cadastro", "application/json", bytes.NewBuffer(json_valor)) // Faz a requisição POST
	if err != nil {
		fmt.Println("Erro ao fazer a requisição POST:", err)
		return false
	}
	defer resposta.Body.Close()
	fmt.Println("POST: sucesso")

	var resposta_map map[string]interface{}                                      // Mapa para decodificar o JSON
	if err := json.NewDecoder(resposta.Body).Decode(&resposta_map); err != nil { // Decodifica o JSON
		fmt.Println("Erro ao decodificar o JSON:", err)
		return false
	}
	fmt.Println("Decodificado: sucesso")
	id_receb, ok := resposta_map["id"].(float64) // Converte o ID para int
	if !ok {
		fmt.Println("Erro ao converter o ID")
		return false
	}
	fmt.Println("ID recebido:", id_receb)
	id = int(id_receb) // Atribui o ID recebido
	return true
}

func criar_evento(nome string, descricao string, porcentagemCriador float64) bool {
	var cria_evento = Cria_Evento_req{Id: id, Nome: nome, Descricao: descricao, PorcentagemCriador: porcentagemCriador}
	json_valor, err := json.Marshal(cria_evento) // Serializa o JSON
	if err != nil {
		fmt.Println("Erro ao serializar o JSON:", err)
		return false
	}
	fmt.Println("Serializado: sucesso")

	resposta, err := http.Post(url+"/cria_evento", "application/json", bytes.NewBuffer(json_valor)) // Faz a requisição POST
	if err != nil {
		fmt.Println("Erro ao fazer a requisição POST:", err)
		return false
	}
	fmt.Println("POST: sucesso")
	defer resposta.Body.Close()

	var resposta_map map[string]interface{}                                      // Mapa para decodificar o JSON
	if err := json.NewDecoder(resposta.Body).Decode(&resposta_map); err != nil { // Decodifica o JSON
		fmt.Println("Erro ao decodificar o JSON:", err)
		return false
	}
	fmt.Println("Decodificado: sucesso")

	return true
}

func limpar_buffer() {
	// Limpa o buffer pendente
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n') // Lê o restante da entrada e descarta
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	//var nome string

	limpar_terminal()

	displayMessageWithColors("Bem vindo a melhor BET do cenario!!!", 2)

	// Entrada do endereço do servidor
	fmt.Println("Digite o endereço do servidor: ")
	url, _ = reader.ReadString('\n')
	url = strings.TrimSpace(url)
	limpar_terminal()

	// Entrada do nome de acesso
	fmt.Println("Digite seu nome de acesso: ")
	nomeUser, _ := reader.ReadString('\n')
	nomeUser = strings.TrimSpace(nomeUser)

	if !cadastrar(nomeUser) {
		fmt.Println("Erro ao cadastrar usuário. Tente novamente.")
		return
	}
	// time.Sleep(3 * time.Second)

	limpar_terminal()
	loop := 1
	for loop == 1 {
		limpar_terminal()

		cabecalho(url)
		fmt.Println("	  Menu")
		fmt.Println("1 - Participar de um evento")
		fmt.Println("2 - Criar um evento")
		fmt.Println("3 - Ver eventos [Participados]")
		fmt.Println("4 - Ver eventos [Criados]")
		fmt.Println("5 - Depositar")
		fmt.Println("6 - Sacar")
		fmt.Println("0 - Encerrar sessão")
		
		// Leitura da seleção
		fmt.Println("Escolha uma opção: ")
		selecao, _ := reader.ReadString('\n')
		selecao = strings.TrimSpace(selecao)

		if selecao == "1" { //Participar de um evento
			limpar_terminal()
			//ver lista com os nomes do evento
			//sugestao 1 - colocar o id do evento do lado e o usuario digita
			//sugestao 2 - o usario digita o nome igual as rotas
			//segestao 3 - ambos (AI COMPLICA)
			//aparece os detalhes do evento, nome, descricao e criado
			//opcoes voltar, participar
			//em participar - cobra o pagameto
			//opcoes voltar e confirmar
		} else if selecao == "2" { //Criar um envento
			limpar_terminal()

			// Leitura do nome
			fmt.Println("Defina um nome para seu evento: ")
			nome, _ := reader.ReadString('\n')
			nome = strings.TrimSpace(nome) // Remove espaços em branco e \n


			// Leitura da descrição
			fmt.Println("Defina a descrição do evento: ")
			descricao, _ := reader.ReadString('\n')
			descricao = strings.TrimSpace(descricao) // Remove espaços em branco e \n


			// Leitura da porcentagem
			var porcentagemCriador float64
			for {
				fmt.Println("Defina a porcentagem que você irá receber (são permitidos de 0% a 50%): ")
				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(input) // Remove espaços em branco e \n
				porcentagem, err := strconv.ParseFloat(input, 64)
				if err != nil || porcentagem < 0 || porcentagem > 50 {
					fmt.Println("Porcentagem inválida, digite novamente.")
					continue
				}
				porcentagemCriador = porcentagem
				break
			}
			criar_evento(nome, descricao, porcentagemCriador)
			time.Sleep(3 * time.Second)

			// atribuido valores
			criar := Evento{
				id:                 0,
				ativo:              true,
				id_criador:         0,
				nome:               nome,
				descricao:          descricao,
				participantes:      nil,
				palpite:            nil,
				porcentagemCriador: 0,
				resultado:          "",
			}

			//encaminhar infomaçoes para o servidor
			fmt.Println(criar) //para não gerar erro
		} else if selecao == "3" { //Ver evetos [Participados]
			limpar_terminal()
			//exibir lista de envento que o usuario participa ou participou
			//sugestao 1 - colocar o id do evento do lado e o usuario digita
			//sugestao 2 - o usario digita o nome igual as rotas
			//segestao 3 - ambos (AI COMPLICA)
			//aparece os voltar e detalhe (RPZ TEM QUE INFORMAR O GANHADOR) (REEBOLSO NÃO É UMA OPCAO, JOGOU PQ QUIS)
		} else if selecao == "4" { //Ver evetos [Criados]
			limpar_terminal()
			//opcao ativo e desativado
			//sugestao 1 - colocar o id do evento do lado e o usuario digita
			//sugestao 2 - o usario digita o nome igual as rotas
			//segestao 3 - ambos (AI COMPLICA)
			//desativado - opcoes volta e detalhe (NAO PODE EDITAR NADA)
			//ativo - opcoes voltar, informar vencedor (isso poes como evento finalizado) (NADA MAIS PODE SER ALTERADO)

		} else if selecao == "5" { //Depositar
			limpar_terminal()
			//solicitar valor
			//encaminhar para o servidor
			//receber confirmação
		} else if selecao == "6" { //Sacar
			limpar_terminal()
			//solicitar valor
			//encaminhar para o servidor
			//receber confirmação
		} else if selecao == "0" {
			limpar_terminal()
			displayMessageWithColors("Te vejo em breve :D", 3)
			loop = 0
		} else {
			limpar_terminal()
			println("Essa opção não existe")
			time.Sleep(1 * time.Second)
		}
	}
}
