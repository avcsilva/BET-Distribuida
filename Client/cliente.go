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

// Estrutura de dados para o Evento ativo
type Evento_ativo struct {
	Nome      string `json:"nome"`      // Nome do evento
	Descricao string `json:"descricao"` // Descrição do evento
}

// Estrutura de dados para o Evento criado
type Evento_criado struct {
	ativo              bool    // Status do evento (ativo ou inativo)
	nome               string  // Nome do evento
	descricao          string  // Descrição do evento
	porcentagemCriador float64 // Porcentagem que o criador irá receber
	resultado          string  // Resultado do evento
}

// Estrutura de dados para o Evento participado
type Evento_participado struct {
	ativo     bool   // Status do evento (ativo ou inativo)
	nome      string // Nome do evento
	descricao string // Descrição do evento
	resultado string // Resultado do evento
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

type Saldo_req struct {
	Id    int     `json:"id"`
	Saldo float64 `json:"saldo"`
}

type Participa_Evento_req struct {
	Id_evento  int     `json:"id_evento"`
	Id_cliente int     `json:"id_cliente"`
	Palpite    string  `json:"palpite"`
	Valor      float64 `json:"valor"`
}

type Altera_Resultado_req struct {
	Id_evento  int    `json:"id_evento"`
	Id_cliente int    `json:"id_cliente"`
	Resultado  string `json:"resultado"`
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
	fmt.Println("|\033[32m                 BET - Distribuida              	 \033[0m|")
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

func participarDoEvento(eventoID int, palpite string, valorAposta float64) bool {
	var dados = Participa_Evento_req{Id_evento: eventoID, Id_cliente: id, Palpite: palpite, Valor: valorAposta}

	json_valor, err := json.Marshal(dados) // Serializa o JSON
	if err != nil {
		fmt.Println("Erro ao serializar o JSON:", err)
		return false
	}
	fmt.Println("Serializado: sucesso")

	resposta, err := http.Post(url+"/participa_evento", "application/json", bytes.NewBuffer(json_valor)) // Faz a requisição POST
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

func alterar_saldo(saldo float64) bool {
	var alterar_saldo = Saldo_req{Id: id, Saldo: saldo}
	json_valor, err := json.Marshal(alterar_saldo) // Serializa o JSON
	if err != nil {
		fmt.Println("Erro ao serializar o JSON:", err)
		return false
	}
	fmt.Println("Serializado: sucesso")

	req, err := http.NewRequest(http.MethodPatch, url+"/alt_saldo", bytes.NewBuffer(json_valor)) // Cria uma requisição PATCH
	if err != nil {
		fmt.Println("Erro ao criar a requisição PATCH:", err)
		return false
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resposta, err := client.Do(req) // Envia a requisição PATCH
	if err != nil {
		fmt.Println("Erro ao fazer a requisição PATCH:", err)
		return false
	}
	fmt.Println("PATCH: sucesso")
	defer resposta.Body.Close()

	var resposta_map map[string]interface{}                                      // Mapa para decodificar o JSON
	if err := json.NewDecoder(resposta.Body).Decode(&resposta_map); err != nil { // Decodifica o JSON
		fmt.Println("Erro ao decodificar o JSON:", err)
		return false
	}
	fmt.Println("Decodificado: sucesso")

	return true
}

func obter_saldo() (float64, error) {
	resposta, err := http.Get(url + fmt.Sprintf("/cliente?id=%d", id)) // Faz a requisição GET com parâmetro de id
	if err != nil {
		fmt.Println("Erro ao fazer a requisição GET:", err)
		return 0, err
	}
	defer resposta.Body.Close()

	var resposta_map map[string]interface{}
	if err := json.NewDecoder(resposta.Body).Decode(&resposta_map); err != nil {
		fmt.Println("Erro ao decodificar o JSON:", err)
		return 0, err
	}

	saldo, ok := resposta_map["saldo"].(float64)
	if !ok {
		return 0, fmt.Errorf("erro ao converter o saldo")
	}
	return saldo, nil
}

func eventos_ativos() map[int]Evento_ativo {
	resposta, err := http.Get(url + "/eventos_ativos")
	if err != nil {
		fmt.Println("Erro ao fazer a requisição GET:", err)
		return nil
	}
	defer resposta.Body.Close()

	var eventos_ativos map[int]Evento_ativo
	if err := json.NewDecoder(resposta.Body).Decode(&eventos_ativos); err != nil {
		fmt.Println("Erro ao decodificar o JSON:", err)
		return nil
	}
	return eventos_ativos
}

// Função para receber todos os eventos relacionados ao cliente, tanto criados quanto participados
func eventos_cliente() map[string]interface{} {
	resposta, err := http.Get(fmt.Sprintf("%s/eventos_cliente?id=%d", url, id)) // Faz a requisição GET com parâmetro de id
	if err != nil {
		fmt.Println("Erro ao fazer a requisição GET:", err)
		return nil
	}
	defer resposta.Body.Close()

	var eventos_cliente_interface map[string]interface{}
	if err := json.NewDecoder(resposta.Body).Decode(&eventos_cliente_interface); err != nil {
		fmt.Println("Erro ao decodificar o JSON:", err)
		return nil
	}

	return eventos_cliente_interface
}

func alterarResultado(eventoID int, resultado string) bool {
    var dados = Altera_Resultado_req{Id_evento: eventoID, Id_cliente: id, Resultado: resultado}

    json_valor, err := json.Marshal(dados) // Serializa o JSON
    if err != nil {
        fmt.Println("Erro ao serializar o JSON:", err)
        return false
    }
    fmt.Println("Serializado: sucesso")

    req, err := http.NewRequest(http.MethodPatch, url+"/altera_resultado", bytes.NewBuffer(json_valor)) // Cria uma requisição PATCH
    if err != nil {
        fmt.Println("Erro ao criar a requisição PATCH:", err)
        return false
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resposta, err := client.Do(req) // Envia a requisição PATCH
    if err != nil {
        fmt.Println("Erro ao fazer a requisição PATCH:", err)
        return false
    }
    fmt.Println("PATCH: sucesso")
    defer resposta.Body.Close()

    var resposta_map map[string]interface{}                                      // Mapa para decodificar o JSON
    if err := json.NewDecoder(resposta.Body).Decode(&resposta_map); err != nil { // Decodifica o JSON
        fmt.Println("Erro ao decodificar o JSON:", err)
        return false
    }
    fmt.Println("Decodificado: sucesso")

    status, ok := resposta_map["status"].(string)
    if !ok || status != "alterado" {
        fmt.Println("Erro ao alterar o resultado.")
        return false
    }

    return true
}

func limpar_buffer() {
	// Limpa o buffer pendente
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n') // Lê o restante da entrada e descarta
}

// Função para validar a URL do servidor
func validarURL(url string) bool {
	_, err := http.Get(url + "/infos")
	return err == nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	//var nome string

	limpar_terminal()

	displayMessageWithColors("Bem vindo a melhor BET do cenario!!!", 2)

	// Entrada do endereço do servidor
	limpar_terminal()
	for {

		fmt.Println("Digite o endereço do servidor: ")
		url, _ = reader.ReadString('\n')
		url = strings.TrimSpace(url)
		if validarURL(url) {
			break
		} else {
			limpar_terminal()
			fmt.Println("\033[31mNão foi possível conectar ao servidor. Tente novamente.\033[0m")
		}
	}
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
		fmt.Println("5 - Enviar resultado de um evento")
		fmt.Println("6 - Depositar")
		fmt.Println("7 - Sacar")
		fmt.Println("0 - Encerrar sessão")

		// Leitura da seleção
		fmt.Println("Escolha uma opção: ")
		selecao, _ := reader.ReadString('\n')
		selecao = strings.TrimSpace(selecao)

		if selecao == "1" { //Participar de um evento
			limpar_terminal()
			eventos := eventos_ativos()
			if len(eventos) == 0 {
				fmt.Println("Não há eventos ativos no momento.")
			} else {
				var eventoID int
				fmt.Println("Eventos ativos:")
				fmt.Println("--------------------------------------------------")
				for id, evento := range eventos {
					fmt.Printf("ID: %d \nNome: %s \nDescrição: %s\n", id, evento.Nome, evento.Descricao)
					fmt.Println("--------------------------------------------------")
				}

				for {
					fmt.Println("Digite o ID do evento que deseja participar: ")
					input, _ := reader.ReadString('\n')
					input = strings.TrimSpace(input)
					tempEventoID, err := strconv.Atoi(input)

					if err != nil {
						fmt.Println("ID do evento inválido. Tente novamente.")
						continue
					}

					if evento, exists := eventos[tempEventoID]; exists {
						eventoID = tempEventoID
						limpar_terminal()
						fmt.Printf("Evento selecionado:\n")
						fmt.Printf("id: %d\n", eventoID)
						fmt.Printf("Nome: %s\n", evento.Nome)
						fmt.Printf("Descrição: %s\n", evento.Descricao)
						fmt.Println("--------------------------------------------------")
						break
					} else {
						fmt.Println("ID do evento inválido. Tente novamente.")
					}
				}

				fmt.Println("antes do palpite eventoID:", eventoID)
				var palpite string
				var valorAposta float64
				fmt.Println("Digite seu palpite: ")
				palpite, _ = reader.ReadString('\n')
				palpite = strings.TrimSpace(palpite)

				for {
					saldoAtual, err := obter_saldo()
					if err != nil {
						fmt.Println("Erro ao obter o saldo atual:", err)
						continue
					}
					fmt.Printf("Seu saldo atual é: %.2f\n", saldoAtual)
					fmt.Println("Digite o valor da aposta: ")
					input, _ := reader.ReadString('\n')
					input = strings.TrimSpace(input)
					valorAposta, err := strconv.ParseFloat(input, 64)

					if err != nil || valorAposta <= 0 || valorAposta > saldoAtual {
						fmt.Println("Valor da aposta inválido. Deve ser maior que zero e menor ou igual ao seu saldo atual.")
						continue
					}
					break
				}
				fmt.Println("eventoID:", eventoID)
				fmt.Println("palpite:", palpite)
				fmt.Println("valorAposta:", valorAposta)
				if participarDoEvento(eventoID, palpite, valorAposta) {
					fmt.Println("Aposta realizada com sucesso!")
				} else {
					fmt.Println("Erro ao realizar aposta.")
				}
			}
			fmt.Println("Pressione Enter para voltar ao menu.")
			reader.ReadString('\n')
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

		} else if selecao == "3" { //Ver evetos [Participados]
			limpar_terminal()
			eventos := eventos_cliente()
			eventos_participados, ok := eventos["eventos_participados"].(map[string]interface{})
			if !ok || len(eventos_participados) == 0 {
				fmt.Println("Você não participou de nenhum evento.")
			} else {
				fmt.Println("Eventos participados:")
				fmt.Println("--------------------------------------------------")
				for id, evento := range eventos_participados {
					evento_map := evento.(map[string]interface{})
					status := "inativo"
					if evento_map["ativo"].(bool) {
						status = "ativo"
					}
					resposta := evento_map["resultado"].(string)
					if resposta == "" {
						resposta = "Sem resposta"
					}
					fmt.Printf("ID: %s \nNome: %s \nDescrição: %s \nStatus: %s \nResposta: %s\n", id, evento_map["nome"], evento_map["descricao"], status, resposta)
					fmt.Println("--------------------------------------------------")
				}
			}
			fmt.Println("Pressione Enter para voltar ao menu.")
			reader.ReadString('\n')
		} else if selecao == "4" { //Ver evetos [Criados]
			limpar_terminal()
			eventos := eventos_cliente()
			eventos_criados, ok := eventos["eventos_criados"].(map[string]interface{})
			if !ok || len(eventos_criados) == 0 {
				fmt.Println("Você não criou nenhum evento.")
			} else {
				fmt.Println("Eventos criados:")
				fmt.Println("--------------------------------------------------")
				for id, evento := range eventos_criados {
					evento_map := evento.(map[string]interface{})
					status := "inativo"
					if evento_map["ativo"].(bool) {
						status = "ativo"
					}
					resposta := evento_map["resultado"].(string)
					if resposta == "" {
						resposta = "Sem resposta"
					}
					porcentagemCriador := evento_map["porcentagemCriador"].(float64)
					fmt.Printf("ID: %s \nNome: %s \nDescrição: %s \nStatus: %s \nResposta: %s \nPorcentagem do Criador: %.1f%%\n", id, evento_map["nome"], evento_map["descricao"], status, resposta, porcentagemCriador)
					fmt.Println("--------------------------------------------------")
				}
			}
			fmt.Println("Pressione Enter para voltar ao menu.")
			reader.ReadString('\n')

		} else if selecao == "5" { //Enviar resultado de um evento
			limpar_terminal()
    eventos := eventos_ativos()
    if len(eventos) == 0 {
        fmt.Println("Não há eventos ativos no momento.")
    } else {
        var eventoID int
        fmt.Println("Eventos ativos:")
        fmt.Println("--------------------------------------------------")
        for id, evento := range eventos {
            fmt.Printf("ID: %d \nNome: %s \nDescrição: %s\n", id, evento.Nome, evento.Descricao)
            fmt.Println("--------------------------------------------------")
        }

        for {
            fmt.Println("Digite o ID do evento que deseja alterar o resultado: ")
            input, _ := reader.ReadString('\n')
            input = strings.TrimSpace(input)
            tempEventoID, err := strconv.Atoi(input)

            if err != nil {
                fmt.Println("ID do evento inválido. Tente novamente.")
                continue
            }

            if evento, exists := eventos[tempEventoID]; exists {
                eventoID = tempEventoID
                limpar_terminal()
                fmt.Printf("Evento selecionado:\n")
                fmt.Printf("id: %d\n", eventoID)
                fmt.Printf("Nome: %s\n", evento.Nome)
                fmt.Printf("Descrição: %s\n", evento.Descricao)
                fmt.Println("--------------------------------------------------")
                break
            } else {
                fmt.Println("ID do evento inválido. Tente novamente.")
            }
        }

        fmt.Println("Digite o resultado do evento: ")
        resultado, _ := reader.ReadString('\n')
        resultado = strings.TrimSpace(resultado)

        fmt.Println("eventoID:", eventoID)
        fmt.Println("resultado:", resultado)
        if alterarResultado(eventoID, resultado) {
            fmt.Println("Resultado alterado com sucesso!")
        } else {
            fmt.Println("Erro ao alterar o resultado.")
        }
    }
    fmt.Println("Pressione Enter para voltar ao menu.")
    reader.ReadString('\n')

		} else if selecao == "6" { //Depositar
			limpar_terminal()
			reader := bufio.NewReader(os.Stdin)

			saldoAtual, err := obter_saldo()
			if err != nil {
				fmt.Println("Erro ao obter o saldo atual:", err)
				continue
			}
			fmt.Printf("Seu saldo atual é: %.2f\n", saldoAtual)

			for {
				fmt.Println("Digite o valor do depósito (ou digite '0' para encerrar): ")
				saldoStr, _ := reader.ReadString('\n') // Lê a entrada do usuário
				saldoStr = strings.TrimSpace(saldoStr) // Remove espaços extras e nova linha

				// Verifica se o usuário deseja sair
				if strings.ToLower(saldoStr) == "0" {
					break
				}

				// Converte a entrada para float64
				saldo, err := strconv.ParseFloat(saldoStr, 64)
				if err != nil {
					limpar_terminal()
					fmt.Println("Erro: valor inválido. Certifique-se de inserir um número válido.")
					continue // Volta ao início do loop
				}

				// Verifica se o valor é maior que 0
				if saldo <= 0 {
					limpar_terminal()
					fmt.Println("Erro: valor inválido. O valor deve ser maior que 0.")
					continue // Volta ao início do loop
				}

				alterar_saldo(saldo)
				limpar_terminal()
				saldoAtual, err := obter_saldo()
				if err != nil {
					fmt.Println("Erro ao obter o saldo atual:", err)
					continue
				}
				fmt.Printf("\033[32mO valor do deposito é: %.2f\033[0m\n\033[34mSeu saldo atual é: %.2f\n\033[0m\n", saldo, saldoAtual)
			}

		} else if selecao == "7" { //Sacar
			limpar_terminal()
			reader := bufio.NewReader(os.Stdin)

			saldoAtual, err := obter_saldo()
			if err != nil {
				fmt.Println("Erro ao obter o saldo atual:", err)
				continue
			}
			fmt.Printf("Seu saldo atual é: %.2f\n", saldoAtual)

			for {
				fmt.Println("Digite o valor do saque (ou digite '0' para encerrar): ")
				saldoStr, _ := reader.ReadString('\n') // Lê a entrada do usuário
				saldoStr = strings.TrimSpace(saldoStr) // Remove espaços extras e nova linha

				// Verifica se o usuário deseja sair
				if strings.ToLower(saldoStr) == "0" {
					break
				}

				// Converte a entrada para float64
				saldo, err := strconv.ParseFloat(saldoStr, 64)
				if err != nil {
					limpar_terminal()
					fmt.Println("Erro: valor inválido. Certifique-se de inserir um número válido.")
					continue // Volta ao início do loop
				}

				// Verifica se o valor é maior que 0 e menor ou igual ao saldo atual
				if saldo <= 0 || saldo > saldoAtual {
					limpar_terminal()
					fmt.Printf("Erro: valor inválido. O valor deve ser maior que 0 e menor ou igual ao saldo atual (%.2f).\n", saldoAtual)
					continue // Volta ao início do loop
				}

				alterar_saldo(-saldo)
				limpar_terminal()
				saldoAtual, err := obter_saldo()
				if err != nil {
					fmt.Println("Erro ao obter o saldo atual:", err)
					continue
				}
				fmt.Printf("\033[32mO valor do saque é: %.2f\033[0m\n\033[34mSeu saldo atual é: %.2f\n\033[0m\n", saldo, saldoAtual)
			}
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
