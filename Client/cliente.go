package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

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
			time.Sleep(500 * time.Millisecond)         // Delay de 500ms

			// Verifica novamente se o tempo foi atingido dentro do loop interno
			if time.Since(start) > duration {
				break
			}
		}
	}
}

func main() {
	limpar_terminal()
	type Evento struct {
		id         		int
		ativo      		bool
		id_criador 		int
		nome      		string
		descricao		string
		participantes 	map[int]float64  //id dos participantes e valor pago
		palpite			map[int]string //id dos participantes e palpite
		porcentagemCriador float64
		resultado		string
	}

	displayMessageWithColors("Bem vindo a melhor BET do cenario!!!", 2)

	limpar_terminal()
	loop := 1
	for loop == 1 {
		limpar_terminal()
		var selecao string
		
		fmt.Println("	  Menu")
		fmt.Println("1 - Participar de um evento")
		fmt.Println("2 - Criar um evento")
		fmt.Println("3 - Ver eventos [Participados]")
		fmt.Println("4 - Ver eventos [Criados]")
		fmt.Println("5 - Depositar")
		fmt.Println("6 - Sacar") 
		fmt.Println("0 - Encerrar sessão")
		fmt.Scan(&selecao)

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

			// solicitar informacoes do usuario
			var nome string
			fmt.Println("Defina um nome para seu evento: ")
			fmt.Scan(&nome)
			var descricao string
			fmt.Println("Defina a descrição do evento: ")
			fmt.Scan(&descricao)
			var porcentagemCriador float64
			fmt.Println("Defina a porcentagem que o voce irá receber (são pemitido de 0% a 50%): ")
			fmt.Scan(&porcentagemCriador)
			for porcentagemCriador < 0 || porcentagemCriador > 50 {
				fmt.Println("Porcentagem invalida, digite novamente: ")
				fmt.Scan(&porcentagemCriador)
			}

			// atribuido valores
			criar := Evento{
				id:  0,
				ativo: true,
				id_criador: 0,
				nome: nome,
				descricao: descricao,
				participantes: nil,
				palpite: nil,
				porcentagemCriador: 0,
				resultado: "",
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
