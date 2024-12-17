package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Endereços de servidores para testes e usos locais
var serv_loc_A string = "http://localhost:8080"
var serv_loc_B string = "http://localhost:8081"
var serv_loc_C string = "http://localhost:8082"

// Endereços de servidores para testes e usos em laboratório
var serv_lab_A string = "http://172.16.103.14:8080"
var serv_lab_B string = "http://172.16.103.13:8080"
var serv_lab_C string = "http://172.16.103.12:8080"

// Variáveis para o uso de Token Ring
var token bool = false       // Variável para controle de token. Indica se o servidor possui ou não o token
var token_exist bool = false // Variável para controle de existência de token no sistema. Indica se o token existe ou não no sistema
var trava_token bool = false // Variável para controle de trava do token. Indica se o servidor deve ou não segurar o token
var tarefa_ok bool = true    // Variável para controle de tarefa. Indica se a tarefa do servidor foi concluída ou não

// Variáveis para lógica do sistema
var id_cont_cliente int = 1 // Variável para controle de IDs dos clientes
var id_cont_evento int = 1  // Variável para controle de IDs dos eventos

// Estruturas de dados
// Estrutura de dados para o Cliente
type Cliente struct {
	id                      int     // ID do cliente
	nome                    string  // Nome do cliente
	saldo                   float64 // Saldo do cliente
	id_eventos_criados      []int   // Lista de IDs dos eventos criados pelo cliente
	id_eventos_participados []int   // Lista de IDs dos eventos que o cliente participou
}

// Estrutura de dados para o Evento
type Evento struct {
	id                 int             // ID do evento
	ativo              bool            // Status do evento (ativo ou inativo)
	id_criador         int             // ID do cliente criador do evento
	nome               string          // Nome do evento
	descricao          string          // Descrição do evento
	participantes      map[int]float64 // Mapa de participantes baseado em seu ID e valor de aposta (Ex.: {1: 30, 2: 25})
	palpite            map[int]string  // Mapa de palpites baseado em seu ID e resultado (Ex.: {1: "A", 2: "B"})
	porcentagemCriador float64         // Porcentagem do valor total do evento que o criador receberá
	resultado          string          // Resultado do evento
}

// Estrutura de dados para as informações locais do servidor em questão
type Infos_local struct {
	servidores []string        // Lista de servidores quem que deverá se comunicar
	porta      string          // Porta do servidor
	qual_serv  string          // Servidor em questão
	clientes   map[int]Cliente // Mapa de clientes baseado em seu ID (Ex.: {1: Cliente1, 2: Cliente2})
	eventos    map[int]Evento  // Mapa de eventos baseado em seu ID (Ex.: {1: Evento1, 2: Evento2})
}

// Estrutura de dados para as requisições de cadastro de cliente
type Cadastro_req struct {
	Id   int    `json:"id"`
	Nome string `json:"nome"`
}

// Estrutura de dados para as requisições de criação de evento
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

// Fim das estruturas de dados

// Funções para o Token Ring
// Função para validar existência de token no sistema com 3 servidores
func existe_token(serv_local *Infos_local) {
	for {
		// Verifica se a variável token se torna verdadeira dentro de um tempo limite
		// Ou seja, verifica se há um token circulando pelo sistema em um período de 3 segundos
		for i := 0; i < 3; i++ {
			if token || token_exist { // Se houver token, o servidor atual não precisa gerar um novo
				token_exist = false
				i = -1
				continue
			}
			// Conta 1 segundo
			time.Sleep(1 * time.Second)
			fmt.Printf("%ss, token ainda não encontrado.\n", strconv.Itoa(i+1))
		} // Caso não haja token em 3 segundos, o servidor atual gerará um novo token

		if serv_local.qual_serv == "A" { // Tempo de espera para cada servidor gerar um novo token
			time.Sleep(100 * time.Millisecond) // Servidor A espera apenas 0,1 segundos para gerar um novo token
			if token || token_exist {          // Se o token já foi gerado por outro servidor, o servidor atual não precisa gerar um novo
				fmt.Println("Token encontrado.")
				token_exist = false
				continue
			}
			fmt.Println("Token gerado servidor A.")
		} else if serv_local.qual_serv == "B" { // Tempo de espera para cada servidor gerar um novo token
			time.Sleep(300 * time.Millisecond) // Servidor B espera 0,3 segundos para gerar um novo token
			if token || token_exist {          // Se o token já foi gerado por outro servidor, o servidor atual não precisa gerar um novo
				fmt.Println("Token encontrado.")
				token_exist = false
				continue
			}
			fmt.Println("Token gerado servidor B.")
		} else if serv_local.qual_serv == "C" { // Tempo de espera para cada servidor gerar um novo token
			time.Sleep(600 * time.Millisecond) // Servidor C espera 0,6 segundos para gerar um novo token
			if token || token_exist {          // Se o token já foi gerado por outro servidor, o servidor atual não precisa gerar um novo
				fmt.Println("Token encontrado.")
				token_exist = false
				continue
			}
			fmt.Println("Token gerado servidor C.")
		}
		token = true // O servidor atual possui o token
		// Envia confirmação de existência de token no sistema para o próximo servidor
		for _, servidor := range serv_local.servidores { // O servidor atual tentará enviar a confirmação para os próximos servidores
			if !(envia_req_token_exist(servidor)) { // Se a confirmação não for enviada com sucesso, o servidor atual tentará enviar para o próximo servidor
				fmt.Printf("Erro ao enviar confirmação de existência de token para o servidor %s\n", servidor)
			}
		}
	}
}

// Função para enviar o token via HTTP para o próximo servidor
func envia_req_token(servidor string) bool {
	resposta, erro := http.Post(servidor+"/token", "application/json", nil) // Envia requisição POST para o servidor
	if erro != nil {                                                        // Se houver erro, retorna false
		return false
	}

	if resposta.StatusCode == 200 { // Se a resposta for 200, retorna true
		return true
	} else {
		return false
	}
}

// Função para enviar a confirmção de existência de token via HTTP para o próximo servidor
func envia_req_token_exist(servidor string) bool {
	resposta, erro := http.Post(servidor+"/token_exist", "application/json", nil) // Envia requisição POST para o servidor
	if erro != nil {                                                              // Se houver erro, retorna false
		return false
	}

	if resposta.StatusCode == 200 { // Se a resposta for 200, retorna true
		return true
	} else {
		return false
	}
}

// Função para segurar e passar o token para o próximo servidor
func passa_token(serv_local *Infos_local) {
	for {
		if token { // Se o servidor possuir o token, irá mantê-lo por 1 seguno ou até terminar sua tarefa
			fmt.Println("Token em posse do servidor.")
			trava_token = true
			for trava_token || !tarefa_ok { // Enquanto a trava estiver ativa ou a tarefa não estiver concluída, o servidor segura o token
				fmt.Println("Token em espera.")
				time.Sleep(1 * time.Second)
				trava_token = false
			}
			fmt.Println("Token em liberação.")

			// Com trava não ativa e tarefa concluída, o servidor passa o token para o próximo
			token = false // O servidor atual não possui mais o token

			for index, servidor := range serv_local.servidores { // O servidor atual tentará passar o token para os próximos servidores
				if envia_req_token(servidor) { // Se o token for enviado com sucesso, o servidor atual encerra o loop
					fmt.Printf("Token enviado para o servidor %s\n", servidor)
					break
				} else if index == 0 { // Se o token não for enviado com sucesso, o servidor atual tentará enviar para o próximo servidor
					fmt.Printf("Erro ao enviar token para o servidor %s. Tentando enviar ao próximo.\n", servidor)
					continue
				} else { // Se o token não for enviado com sucesso, o servidor atual permanecerá com o token
					fmt.Printf("Erro ao enviar token para o servidor %s\n", servidor)
					fmt.Printf("Token retornado ao servidor %s\n", serv_local.qual_serv)
					token = true // Mantém o token, pois não há outros servidores para enviar
					break
				}
			}
		}
	}
}

// Função para atualização de informações do servidor, como clientes e eventos e as contagens de IDs
// Atualização realizada com base nas informações dos outros servidores
// Em ideia principal, deve ser chamada apenas na inicialização do servidor
func atualiza_infos(serv_local *Infos_local) {
	for _, servidor := range serv_local.servidores {
		resposta, erro := http.Get(servidor + "/infos") // Requisição GET para obter as informações do servidor
		if erro != nil {                                // Se houver erro, o servidor não será atualizado
			fmt.Printf("Erro ao solicitar informações do servidor %s\n", servidor)
			continue // O servidor tentará obter informações do próximo servidor
		}

		defer resposta.Body.Close() // Fecha o corpo da resposta após o uso

		// Criação de estrutura para armazenar as informações obtidas do servidor
		infos := make(map[string]interface{})
		if erro := json.NewDecoder(resposta.Body).Decode(&infos); erro != nil { // Decodifica as informações obtidas
			fmt.Printf("Erro ao decodificar informações do servidor %s\n", servidor)
			continue // O servidor tentará obter informações do próximo servidor
		}

		// Mensagem de confirmação de recebimento
		// juntamente com as informações recebidas
		fmt.Printf("Informações obtidas do servidor %s: %v\n", servidor, infos)

		// Atualização das informações locais do servidor
		// com base nas informações obtidas

		// Início obtenção da contagem de IDs de clientes
		if id_cli, ok := infos["id_cont_cliente"].(float64); ok {
			id_cont_cliente = int(id_cli)
		} else {
			fmt.Printf("Erro ao obter contagem de IDs de clientes do servidor %s\n", servidor)
		}
		// Fim obtenção da contagem de IDs de clientes

		// Início obtenção da contagem de IDs de eventos
		if id_eve, ok := infos["id_cont_evento"].(float64); ok { // Atualização da contagem de IDs de eventos
			id_cont_evento = int(id_eve)
		} else {
			fmt.Printf("Erro ao obter contagem de IDs de eventos do servidor %s\n", servidor)
		}
		// Fim obtenção da contagem de IDs de eventos

		// Início obtenção dos clientes
		if clientes_externos, ok := infos["clientes"].(map[string]interface{}); ok { // Atualização dos clientes
			for index, item := range clientes_externos {
				cliente, ok := item.(map[string]interface{}) // Conversão do item para um mapa de informações de cliente
				if !ok {
					fmt.Printf("Erro ao converter informações de clientes do servidor %s\n", servidor)
				}

				// Obtenção do ID do cliente
				id, erro := strconv.Atoi(index)
				if erro != nil {
					fmt.Printf("Erro ao converter ID de cliente do servidor %s\n", servidor)
				}

				// Obtenção do nome do cliente
				nome, ok := cliente["nome"].(string)
				if !ok {
					fmt.Printf("Erro ao obter nome de cliente do servidor %s\n", servidor)
				}

				// Obtenção do saldo do cliente
				saldo, ok := cliente["saldo"].(float64)
				if !ok {
					fmt.Printf("Erro ao obter saldo de cliente do servidor %s\n", servidor)
				}

				// Início obtenção dos IDs de eventos criados pelo cliente
				id_ev_cri_slice, ok := cliente["id_eventos_criados"].([]interface{})
				if !ok {
					fmt.Printf("Erro ao obter IDs de eventos criados por cliente do servidor %s\n", servidor)
				}
				id_ev_cri := make([]int, len(id_ev_cri_slice))
				for i, valor := range id_ev_cri_slice {
					id_ev, ok := valor.(float64)
					if !ok {
						fmt.Printf("Erro ao converter ID de evento criado por cliente do servidor %s\n", servidor)
					}
					id_ev_cri[i] = int(id_ev)
				}
				// Fim obtenção dos IDs de eventos criados pelo cliente

				// Início obtenção dos IDs de eventos participados pelo cliente
				id_ev_part_slice, ok := cliente["id_eventos_participados"].([]interface{})
				if !ok {
					fmt.Printf("Erro ao obter IDs de eventos participados por cliente do servidor %s\n", servidor)
				}
				id_ev_part := make([]int, len(id_ev_part_slice))
				for i, valor := range id_ev_part_slice {
					id_ev, ok := valor.(float64)
					if !ok {
						fmt.Printf("Erro ao converter ID de evento participado por cliente do servidor %s\n", servidor)
					}
					id_ev_part[i] = int(id_ev)
				}
				// Fim obtenção dos IDs de eventos participados pelo cliente

				// Atualização do cliente
				serv_local.clientes[id] = Cliente{
					id:                      id,
					nome:                    nome,
					saldo:                   saldo,
					id_eventos_criados:      id_ev_cri,
					id_eventos_participados: id_ev_part,
				}
			}
		} else {
			fmt.Printf("Erro ao obter informações de clientes do servidor %s\n", servidor)
		}
		// Fim obtenção dos clientes

		// Início obtenção dos eventos
		if eventos_externos, ok := infos["eventos"].(map[string]interface{}); ok { // Atualização dos eventos
			for index, item := range eventos_externos {
				evento, ok := item.(map[string]interface{}) // Conversão do item para um mapa de informações de evento
				if !ok {
					fmt.Printf("Erro ao converter informações de eventos do servidor %s\n", servidor)
				}

				// Obtenção do ID do evento
				id, erro := strconv.Atoi(index)
				if erro != nil {
					fmt.Printf("Erro ao converter ID de evento do servidor %s\n", servidor)
				}

				// Obtenção do status do evento
				ativo, ok := evento["ativo"].(bool)
				if !ok {
					fmt.Printf("Erro ao obter status de evento do servidor %s\n", servidor)
				}

				// Obtenção do ID do cliente criador do evento
				id_criador, ok := evento["id_criador"].(float64)
				if !ok {
					fmt.Printf("Erro ao obter ID do criador de evento do servidor %s\n", servidor)
				}
				id_criador_int := int(id_criador)

				// Obtenção do nome do evento
				nome, ok := evento["nome"].(string)
				if !ok {
					fmt.Printf("Erro ao obter nome de evento do servidor %s\n", servidor)
				}

				// Obtenção da descrição do evento
				descricao, ok := evento["descricao"].(string)
				if !ok {
					fmt.Printf("Erro ao obter descrição de evento do servidor %s\n", servidor)
				}

				// Início obtenção dos participantes do evento
				participantes_externos, ok := evento["participantes"].(map[string]interface{})
				if !ok {
					fmt.Printf("Erro ao obter participantes de evento do servidor %s\n", servidor)
				}
				participantes := make(map[int]float64)
				for i, valor := range participantes_externos {
					id_part, erro := strconv.Atoi(i)
					if erro != nil {
						fmt.Printf("Erro ao converter ID de participante de evento do servidor %s\n", servidor)
					}
					valor_part, ok := valor.(float64)
					if !ok {
						fmt.Printf("Erro ao obter valor de participante de evento do servidor %s\n", servidor)
					}
					participantes[id_part] = valor_part
				}
				// Fim obtenção dos participantes do evento

				// Início obtenção dos palpites dos participantes do evento
				palpite_externos, ok := evento["palpite"].(map[string]interface{})
				if !ok {
					fmt.Printf("Erro ao obter palpites de evento do servidor %s\n", servidor)
				}
				palpite := make(map[int]string)
				for i, valor := range palpite_externos {
					id_part, erro := strconv.Atoi(i)
					if erro != nil {
						fmt.Printf("Erro ao converter ID de participante de evento do servidor %s\n", servidor)
					}
					palpite_part, ok := valor.(string)
					if !ok {
						fmt.Printf("Erro ao obter palpite de participante de evento do servidor %s\n", servidor)
					}
					palpite[id_part] = palpite_part
				}
				// Fim obtenção dos palpites dos participantes do evento

				// Obtenção da porcentagem do criador
				porcentagemCriador, ok := evento["porcentagemCriador"].(float64)
				if !ok {
					fmt.Printf("Erro ao obter porcentagem do criador de evento do servidor %s\n", servidor)
				}

				// Obtenção do resultado do evento
				resultado, ok := evento["resultado"].(string)
				if !ok {
					fmt.Printf("Erro ao obter resultado de evento do servidor %s\n", servidor)
				}

				// Atualização do evento
				serv_local.eventos[id] = Evento{
					id:                 id,
					ativo:              ativo,
					id_criador:         id_criador_int,
					nome:               nome,
					descricao:          descricao,
					participantes:      participantes,
					palpite:            palpite,
					porcentagemCriador: porcentagemCriador,
					resultado:          resultado,
				}
			}
		} else {
			fmt.Printf("Erro ao obter informações de eventos do servidor %s\n", servidor)
		}
		// Fim obtenção dos eventos
	}
}

// Função para definir as informações locais do servidor
func define_info() Infos_local {
	var qual_serv, tipo_serv string
	fmt.Printf("Tipo de servidor (local [LOC] ou laboratório [LAB]): ")
	fmt.Scan(&tipo_serv)
	tipo_serv = strings.ToUpper(tipo_serv)
	fmt.Printf("Servidor (A, B ou C): ")
	fmt.Scan(&qual_serv)
	qual_serv = strings.ToUpper(qual_serv)
	for {
		if tipo_serv == "LOC" {
			if qual_serv == "A" {
				servidores := []string{serv_loc_B, serv_loc_C}
				serv_local := Infos_local{
					servidores: servidores,
					porta:      ":8080",
					qual_serv:  qual_serv,
					clientes:   make(map[int]Cliente),
					eventos:    make(map[int]Evento),
				}
				return serv_local
			} else if qual_serv == "B" {
				servidores := []string{serv_loc_C, serv_loc_A}
				serv_local := Infos_local{
					servidores: servidores,
					porta:      ":8081",
					qual_serv:  qual_serv,
					clientes:   make(map[int]Cliente),
					eventos:    make(map[int]Evento),
				}
				return serv_local
			} else if qual_serv == "C" {
				servidores := []string{serv_loc_A, serv_loc_B}
				serv_local := Infos_local{
					servidores: servidores,
					porta:      ":8082",
					qual_serv:  qual_serv,
					clientes:   make(map[int]Cliente),
					eventos:    make(map[int]Evento),
				}
				return serv_local
			} else {
				fmt.Println("Servidor inválido.")
				fmt.Printf("Tipo de servidor (local [LOC] ou laboratório [LAB]): ")
				fmt.Scan(&tipo_serv)
				tipo_serv = strings.ToUpper(tipo_serv)
				fmt.Printf("Servidor (A, B ou C): ")
				fmt.Scan(&qual_serv)
				qual_serv = strings.ToUpper(qual_serv)
			}
		} else if tipo_serv == "LAB" {
			if qual_serv == "A" {
				servidores := []string{serv_lab_B, serv_lab_C}
				serv_local := Infos_local{
					servidores: servidores,
					porta:      ":8080",
					qual_serv:  qual_serv,
					clientes:   make(map[int]Cliente),
					eventos:    make(map[int]Evento),
				}
				return serv_local
			} else if qual_serv == "B" {
				servidores := []string{serv_lab_A, serv_lab_C}
				serv_local := Infos_local{
					servidores: servidores,
					porta:      ":8080",
					qual_serv:  qual_serv,
					clientes:   make(map[int]Cliente),
					eventos:    make(map[int]Evento),
				}
				return serv_local
			} else if qual_serv == "C" {
				servidores := []string{serv_lab_A, serv_lab_B}
				serv_local := Infos_local{
					servidores: servidores,
					porta:      ":8080",
					qual_serv:  qual_serv,
					clientes:   make(map[int]Cliente),
					eventos:    make(map[int]Evento),
				}
				return serv_local
			} else {
				fmt.Println("Servidor inválido.")
				fmt.Printf("Tipo de servidor (local [LOC] ou laboratório [LAB]): ")
				fmt.Scan(&tipo_serv)
				tipo_serv = strings.ToUpper(tipo_serv)
				fmt.Printf("Servidor (A, B ou C): ")
				fmt.Scan(&qual_serv)
				qual_serv = strings.ToUpper(qual_serv)
			}
		} else {
			fmt.Println("Tipo de servidor inválido.")
			fmt.Printf("Tipo de servidor (local [LOC] ou laboratório [LAB]): ")
			fmt.Scan(&tipo_serv)
			tipo_serv = strings.ToUpper(tipo_serv)
			fmt.Printf("Servidor (A, B ou C): ")
			fmt.Scan(&qual_serv)
			qual_serv = strings.ToUpper(qual_serv)
		}
	}
}

// Função para definir os métodos GET do servidor
func define_metodo_get(serv_local *Infos_local, serv *gin.Engine) {
	// Método GET para retornar informações armazenadas pelo servidor, como clientes e eventos e as contagens de IDs
	serv.GET("/infos", func(c *gin.Context) {
		// Criação de estrutura para informações de clientes
		clientes := make(map[int]map[string]interface{})
		for id, cliente := range serv_local.clientes {
			clientes[id] = map[string]interface{}{
				"nome":                    cliente.nome,
				"saldo":                   cliente.saldo,
				"id_eventos_criados":      cliente.id_eventos_criados,
				"id_eventos_participados": cliente.id_eventos_participados,
			}
		}

		// Criação de estrutura para informações de eventos
		eventos := make(map[int]map[string]interface{})
		for id, evento := range serv_local.eventos {
			eventos[id] = map[string]interface{}{
				"ativo":         evento.ativo,
				"id_criador":    evento.id_criador,
				"nome":          evento.nome,
				"descricao":     evento.descricao,
				"participantes": evento.participantes,
				"palpite":       evento.palpite,
				"resultado":     evento.resultado,
			}
		}

		// Criação de estrutura para o retorno geral de informações
		infos := map[string]interface{}{
			"clientes":        clientes,
			"eventos":         eventos,
			"id_cont_cliente": id_cont_cliente,
			"id_cont_evento":  id_cont_evento,
		}

		c.JSON(http.StatusOK, infos) // Retorna as informações armazenadas pelo servidor
	})

	// Método GET para retornar eventos ativos
	serv.GET("/eventos_ativos", func(c *gin.Context) {
		// Criação de estrutura para armazenar os eventos ativos
		eventos_ativos := make(map[int]map[string]string)
		for id, evento := range serv_local.eventos {
			if evento.ativo { // Verifica se o evento está ativo
				eventos_ativos[id] = map[string]string{
					"nome":          evento.nome,
					"descricao":     evento.descricao,
				}
			}
		}

		c.JSON(http.StatusOK, eventos_ativos) // Retorna os eventos ativos
	})

	// Método GET para retornar eventos relacionados a um cliente, tanto criados quanto participados
	serv.GET("/eventos_cliente", func(c *gin.Context) {
		cliente_id := c.Query("id") // Obtém o ID do cliente
		id, erro := strconv.Atoi(cliente_id)
		if erro != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido."})
			return
		}

		// Criação de estrutura para armazenar os eventos relacionados ao cliente
		eventos_criados := make(map[int]map[string]interface{})
		eventos_participados := make(map[int]map[string]interface{})

		for _, evento := range serv_local.eventos {
			if evento.id_criador == id { // Verifica se o cliente é o criador do evento
				eventos_criados[evento.id] = map[string]interface{}{
					"ativo":     evento.ativo,
					"nome":      evento.nome,
					"descricao": evento.descricao,
					"porcentagemCriador": evento.porcentagemCriador,
					"resultado": evento.resultado,
				}
				continue // Se o cliente é o criador, não é necessário verificar se ele é um participante
			}
			for id_part, _ := range evento.participantes {
				if id_part == id { // Verifica se o cliente é um participante do evento
					eventos_participados[evento.id] = map[string]interface{}{
						"ativo":    evento.ativo,
						"nome":      evento.nome,
						"descricao": evento.descricao,
						"resultado": evento.resultado,
					}
				}
			}
		}

		// Criação de estrutura para o retorno dos eventos relacionados ao cliente
		eventos_cliente := map[string]interface{}{
			"eventos_criados":      eventos_criados,
			"eventos_participados": eventos_participados,
		}

		c.JSON(http.StatusOK, eventos_cliente) // Retorna os eventos relacionados ao cliente
	})
}

// Função para definir os métodos POST do servidor
func define_metodo_post(serv_local *Infos_local, serv *gin.Engine) {
	// Método POST para recebimento do token
	serv.POST("/token", func(c *gin.Context) {
		token = true                                                           // O servidor recebe o token
		c.JSON(http.StatusOK, gin.H{"message": "Token recebido com sucesso."}) // Retorna mensagem de sucesso

		// Envia confirmação de existência de token no sistema para o próximo servidor
		for _, servidor := range serv_local.servidores { // O servidor atual tentará enviar a confirmação para os próximos servidores
			if !(envia_req_token_exist(servidor)) { // Se a confirmação não for enviada com sucesso, o servidor atual tentará enviar para o próximo servidor
				fmt.Printf("Erro ao enviar confirmação de existência de token para o servidor %s\n", servidor)
			}
		}
	})

	// Método POST para confirmação de existência de token no sistema
	serv.POST("/token_exist", func(c *gin.Context) {
		token_exist = true
		c.JSON(http.StatusOK, gin.H{"message": "Token existente."})
	})

	// Método POST para cadastro de um cliente nos servidores
	serv.POST("/cadastro", func(c *gin.Context) {
		var cadastro Cadastro_req                           // Cria uma variável para armazenar o cadastro do cliente
		if err := c.ShouldBindJSON(&cadastro); err != nil { // Faz o bind do JSON recebido para a variável de cadastro
			fmt.Println("Erro ao fazer o bind JSON (/cadastro):", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fonte := c.GetHeader("X-Source") // Verifica o cabeçalho X-Source para saber se veio de um servidor ou cliente
		if fonte == "servidor" {         // Caso seja de um servidor, realizar apenas cadastro próprio do cliente
			serv_local.clientes[cadastro.Id] = Cliente{
				id:                      cadastro.Id,
				nome:                    cadastro.Nome,
				saldo:                   0,
				id_eventos_criados:      []int{},
				id_eventos_participados: []int{},
			}
			id_cont_cliente = cadastro.Id + 1
			c.JSON(http.StatusOK, gin.H{"status": "cadastrado"})
			return
		}

		for !token { // Enquanto o servidor não possuir o token, ele aguardará
			time.Sleep(10 * time.Millisecond) // Adiciona uma pausa de  para evitar o uso de 100% da CPU
		}

		tarefa_ok = false // O servidor atual ainda não concluiu sua tarefa

		//Verificando se o cliente já está cadastrado
		for _, cliente := range serv_local.clientes {
			if cliente.nome == cadastro.Nome { // Caso o cliente já esteja cadastrado, responde com o ID do cliente
				c.JSON(http.StatusOK, gin.H{"status": "logado", "id": cliente.id}) // Responde com o ID do cliente já cadastrado
				tarefa_ok = true                                                   // O servidor atual concluiu sua tarefa
				return
			}
		}

		//Cadastrando o cliente localmente
		serv_local.clientes[id_cont_cliente] = Cliente{
			id:                      id_cont_cliente,
			nome:                    cadastro.Nome,
			saldo:                   0,
			id_eventos_criados:      []int{},
			id_eventos_participados: []int{},
		}

		// Enviando cadastro de cliente aos outros servidores
		cadastro = Cadastro_req{ // Monta a estrutuda de dados para enviar os servidores
			Id:   id_cont_cliente,
			Nome: cadastro.Nome,
		}
		json_valor, err := json.Marshal(cadastro) // Serializa o JSON para enviar aos servidores
		if err != nil {
			fmt.Println("Erro ao serializar o JSON (/cadastro):", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			tarefa_ok = true // O servidor atual concluiu sua tarefa
			return
		}
		for _, server := range serv_local.servidores { // Loop para acessar outros servidores
			go func(server string) {
				req, err := http.NewRequest("POST", server+"/cadastro", bytes.NewBuffer(json_valor)) // Cria uma requisição POST para o servidor
				if err != nil {
					fmt.Printf("Failed to create request to server %s: %v\n", server, err)
					return
				}
				req.Header.Set("Content-Type", "application/json") // Adiciona o cabeçalho Content-Type para identificar que é um JSON
				req.Header.Set("X-Source", "servidor")             // Adiciona o cabeçalho X-Source para identificar que é uma requisição de servidor

				client := &http.Client{}    // Cria um cliente HTTP
				resp, err := client.Do(req) // Envia a requisição
				if err != nil {
					fmt.Printf("Failed to send to server %s: %v\n", server, err)
					return
				}
				defer resp.Body.Close()
			}(server)
		}
		c.JSON(http.StatusOK, gin.H{"status": "cadastrado", "id": id_cont_cliente}) // Responde com o status de cadastrado e o ID do cliente
		id_cont_cliente++                                                           // Incrementa o contador de ID
		tarefa_ok = true                                                            // O servidor atual concluiu sua tarefa
	})

	// Método POST para criação de um evento nos servidores
	serv.POST("/cria_evento", func(c *gin.Context) {
		var cria_evento Cria_Evento_req                        // Cria uma variável para armazenar a criação do evento
		if err := c.ShouldBindJSON(&cria_evento); err != nil { // Faz o bind do JSON recebido para a variável de criação do evento
			fmt.Println("Erro ao fazer o bind JSON (/cria_evento):", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fonte := c.GetHeader("X-Source") // Verifica o cabeçalho X-Source para saber se veio de um servidor ou cliente
		if fonte == "servidor" {         // Caso seja de um servidor, realizar apenas a criação própria do evento
			serv_local.eventos[cria_evento.Id_event] = Evento{
				id:                 cria_evento.Id_event,
				ativo:              true,
				id_criador:         cria_evento.Id,
				nome:               cria_evento.Nome,
				descricao:          cria_evento.Descricao,
				participantes:      make(map[int]float64),
				palpite:            make(map[int]string),
				porcentagemCriador: cria_evento.PorcentagemCriador,
				resultado:          "",
			}
			// Atualizando a lista de eventos criados pelo cliente
			cliente := serv_local.clientes[cria_evento.Id]
			cliente.id_eventos_criados = append(cliente.id_eventos_criados, cria_evento.Id_event)
			serv_local.clientes[cria_evento.Id] = cliente

			id_cont_evento = cria_evento.Id_event + 1
			c.JSON(http.StatusOK, gin.H{"status": "criado"})
			return
		}

		for !token { // Enquanto o servidor não possuir o token, ele aguardará
			time.Sleep(10 * time.Millisecond) // Adiciona uma pausa de 10ms para evitar o uso de 100% da CPU
		}

		tarefa_ok = false // O servidor atual ainda não concluiu sua tarefa

		// Cadastrando o evento localmente
		serv_local.eventos[id_cont_evento] = Evento{
			id:                 id_cont_evento,
			ativo:              true,
			id_criador:         cria_evento.Id,
			nome:               cria_evento.Nome,
			descricao:          cria_evento.Descricao,
			participantes:      make(map[int]float64),
			palpite:            make(map[int]string),
			porcentagemCriador: cria_evento.PorcentagemCriador,
			resultado:          "",
		}

		// Atualizando a lista de eventos criados pelo cliente
		cliente := serv_local.clientes[cria_evento.Id]
		cliente.id_eventos_criados = append(cliente.id_eventos_criados, id_cont_evento)
		serv_local.clientes[cria_evento.Id] = cliente

		// Enviando criação de evento aos outros servidores
		cria_evento = Cria_Evento_req{ // Monta a estrutuda de dados para enviar os servidores
			Id:                 cria_evento.Id,
			Id_event:           id_cont_evento,
			Nome:               cria_evento.Nome,
			Descricao:          cria_evento.Descricao,
			PorcentagemCriador: cria_evento.PorcentagemCriador,
		}
		json_valor, err := json.Marshal(cria_evento) // Serializa o JSON para enviar aos servidores
		if err != nil {
			fmt.Println("Erro ao serializar o JSON (/cria_evento):", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			tarefa_ok = true // O servidor atual concluiu sua tarefa
			return
		}
		for _, server := range serv_local.servidores { // Loop para acessar outros servidores
			go func(server string) {
				req, err := http.NewRequest("POST", server+"/cria_evento", bytes.NewBuffer(json_valor)) // Cria uma requisição POST para o servidor
				if err != nil {
					fmt.Printf("Failed to create request to server %s: %v\n", server, err)
					return
				}
				req.Header.Set("Content-Type", "application/json") // Adiciona o cabeçalho Content-Type para identificar que é um JSON
				req.Header.Set("X-Source", "servidor")             // Adiciona o cabeçalho X-Source para identificar que é uma requisição de servidor

				client := &http.Client{}    // Cria um cliente HTTP
				resp, err := client.Do(req) // Envia a requisição
				if err != nil {
					fmt.Printf("Failed to send to server %s: %v\n", server, err)
					return
				}
				defer resp.Body.Close()
			}(server)
		}
		c.JSON(http.StatusOK, gin.H{"status": "criado", "id": id_cont_evento}) // Responde com o status de criado e o ID do evento
		id_cont_evento++                                                       // Incrementa o contador de ID
		tarefa_ok = true                                                       // O servidor atual concluiu sua tarefa
	})
}

// Função para definir os métodos PATCH do servidor
func define_metodo_patch(serv_local *Infos_local, serv *gin.Engine) {
	// Método PATCH para alteração do saldo de um cliente
	serv.PATCH("/alt_saldo", func(c *gin.Context) {
		var alt_saldo Saldo_req  // Cria uma variável para armazenar a alteração do saldo
		if err := c.ShouldBindJSON(&alt_saldo); err != nil {       // Faz o bind do JSON recebido para a variável de alteração do saldo
			fmt.Println("Erro ao fazer o bind JSON (/alt_saldo):", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fonte := c.GetHeader("X-Source") // Verifica o cabeçalho X-Source para saber se veio de um servidor ou cliente
		if fonte == "servidor" {         // Caso seja de um servidor, realizar apenas a alteração própria do saldo do cliente
			if !alterarSaldo(alt_saldo.Id, alt_saldo.Saldo, serv_local) { // Verifica se o cliente existe e altera o saldo
				c.JSON(http.StatusNotFound, gin.H{"status": "não encontrado"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": "alterado"}) // Responde com o status de alterado
		}

		for !token { // Enquanto o servidor não possuir o token, ele aguardará
			time.Sleep(10 * time.Millisecond) // Adiciona uma pausa de 10ms para evitar o uso de 100% da CPU
		}

		tarefa_ok = false // O servidor atual ainda não concluiu sua tarefa

		if !alterarSaldo(alt_saldo.Id, alt_saldo.Saldo, serv_local) { // Verifica se o cliente existe e altera o saldo
			c.JSON(http.StatusNotFound, gin.H{"status": "não encontrado"})
			return
		}

		// Enviando alteração de saldo aos outros servidores
		json_valor, err := json.Marshal(alt_saldo) // Serializa o JSON para enviar aos servidores
		if err != nil {
			fmt.Println("Erro ao serializar o JSON (/alt_saldo):", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			tarefa_ok = true // O servidor atual concluiu sua tarefa
			return
		}
		for _, server := range serv_local.servidores { // Loop para acessar outros servidores
			go func(server string) {
				req, err := http.NewRequest("PATCH", server+"/alt_saldo", bytes.NewBuffer(json_valor)) // Cria uma requisição PATCH para o servidor
				if err != nil {
					fmt.Printf("Failed to create request to server %s: %v\n", server, err)
					return
				}
				req.Header.Set("Content-Type", "application/json") // Adiciona o cabeçalho Content-Type para identificar que é um JSON
				req.Header.Set("X-Source", "servidor")             // Adiciona o cabeçalho X-Source para identificar que é uma requisição de servidor

				client := &http.Client{}    // Cria um cliente HTTP
				resp, err := client.Do(req) // Envia a requisição
				if err != nil {
					fmt.Printf("Failed to send to server %s: %v\n", server, err)
					return
				}
				defer resp.Body.Close()
			}(server)
		}

		c.JSON(http.StatusOK, gin.H{"status": "alterado"}) // Responde com o status de alterado
		tarefa_ok = true  // O servidor atual concluiu sua tarefa
	})
}

// Função para definir o servidor com os métodos POST, GET e PATCH
func define_servidor(serv_local *Infos_local) *gin.Engine {
	r := gin.Default()

	// Configuração do middleware CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8088"}, // Permitir a origem do frontend
		AllowMethods:     []string{"GET", "POST", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "X-Source"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Define os métodos GET
	define_metodo_get(serv_local, r)

	// Define os métodos POST
	define_metodo_post(serv_local, r)

	// Define os métodos PATCH
	define_metodo_patch(serv_local, r)

	return r
}

func alterarSaldo(id int, novoSaldo float64, il *Infos_local) bool {
	// Verificar se o cliente existe no mapa
	cliente, encontrado := il.clientes[id]
	if !encontrado {
		return false // Retorna falso se o cliente não for encontrado
	}

	// Atualizar o saldo do cliente diretamente no mapa
	cliente.saldo += novoSaldo
	il.clientes[id] = cliente // Reatribuir o cliente ao mapa
	return true
}

// tratameto de dados
func calcularPremio(e Evento, serv_local *Infos_local) {
	premiacao := 0.0                    //Variável para armazenar o valor total das apostas
	total_ganhadores := 0.0             //Variável para armazenar o total apostado pelos ganhadores
	ganhadores := make(map[int]float64) // Mapa para armazenar os ganhadores e seus valores de aposta

	// calcula o valor total
	for _, valor := range e.participantes {
		premiacao += valor
	}

	// Calcula o valor do ganho do criador
	ganhoCriador := (e.porcentagemCriador / 100) * premiacao
	alterarSaldo(e.id_criador, ganhoCriador, serv_local) // Passa o id do criador, o valor do ganho e a estrutura de dados
	premiacao -= ganhoCriador

	for ip, palpt := range e.palpite {
		if palpt == e.resultado {
			ganhadores[ip] = e.participantes[ip]
		}
	}
	// Calcula o total apostado pelos ganhadores
	for _, valor := range ganhadores {
		total_ganhadores += valor
	}
	//pagamento
	for id := range ganhadores {
		ganho := (e.participantes[id] / total_ganhadores) * premiacao
		//atribuir ganho ao saldo do Cliente
		alterarSaldo(id, ganho, serv_local) // Passa o id do ganhador, o valor do ganho e a estrutura de dados
	}
}

func resultadoEvento(id int, resultado string, il *Infos_local) bool {
	// Verificar se o evento existe no mapa
	evento, encontrado := il.eventos[id]
	if !encontrado {
		return false // Retorna falso se o evento não for encontrado
	}

	// Atualizar o resultado do evento diretamente no mapa
	evento.resultado = resultado
	il.eventos[id] = evento // Reatribuir o evento ao mapa

	calcularPremio(evento, il) // Calcula o prêmio e atribui aos ganhadores
	return true
}

func main() {
	serv_local := define_info()
	servidor := define_servidor(&serv_local)
	atualiza_infos(&serv_local)
	go passa_token(&serv_local)
	go existe_token(&serv_local)
	servidor.Run(serv_local.porta)
}
