package main

import (
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
var token bool = false // Variável para controle de token. Indica se o servidor possui ou não o token
var trava_token bool = false // Variável para controle de trava do token. Indica se o servidor deve ou não segurar o token
var tarefa_ok bool = true // Variável para controle de tarefa. Indica se a tarefa do servidor foi concluída ou não

// Variáveis para lógica do sistema
var id_cont_cliente int = 1 // Variável para controle de IDs dos clientes
var id_cont_evento int = 1 // Variável para controle de IDs dos eventos

// Estruturas de dados
// Estrutura de dados para o Cliente
type Cliente struct {
	id   int // ID do cliente
	nome string // Nome do cliente
    saldo float64 // Saldo do cliente
    id_eventos_criados []int // Lista de IDs dos eventos criados pelo cliente
    id_eventos_participados []int // Lista de IDs dos eventos que o cliente participou
}

// Estrutura de dados para o Evento
type Evento struct {
    id              int // ID do evento
    ativo           bool // Status do evento (ativo ou inativo)
    id_criador      int // ID do cliente criador do evento
    nome            string // Nome do evento
    descricao       string // Descrição do evento
    participantes 	map[int]float64 // Mapa de participantes baseado em seu ID e valor de aposta (Ex.: {1: 30, 2: 25})
    palpite         map[int]string // Mapa de palpites baseado em seu ID e resultado (Ex.: {1: "A", 2: "B"})
    resultado       string // Resultado do evento
}

// Estrutura de dados para as informações locais do servidor em questão
type Infos_local struct {
	servidores []string // Lista de servidores quem que deverá se comunicar
	porta      string  // Porta do servidor
    qual_serv string // Servidor em questão
    clientes  map[int]Cliente // Mapa de clientes baseado em seu ID (Ex.: {1: Cliente1, 2: Cliente2})
    eventos   map[int]Evento // Mapa de eventos baseado em seu ID (Ex.: {1: Evento1, 2: Evento2})
}

// Fim das estruturas de dados

// Funções para o Token Ring
// Função para validar existência de token no sistema com 3 servidores
func existe_token(serv_local *Infos_local) {
    for{
        // Verifica se a variável token se torna verdadeira dentro de um tempo limite
        // Ou seja, verifica se há um token circulando pelo sistema em um período de 3 segundos
        for i := 0; i < 3; i++ {
            if token { // Se houver token, o servidor atual não precisa gerar um novo
                fmt.Println("Token encontrado.")
                continue
            }
            // Conta 1 segundo
            time.Sleep(1 * time.Second)
            fmt.Printf("%ss, token ainda não encontrado.\n", strconv.Itoa(i + 1))
        } // Caso não haja token em 3 segundos, o servidor atual gerará um novo token

        if (serv_local.qual_serv == "A") { // Tempo de espera para cada servidor gerar um novo token
            time.Sleep(100 * time.Millisecond) // Servidor A espera apenas 0,1 segundos para gerar um novo token
            if token { // Se o token já foi gerado por outro servidor, o servidor atual não precisa gerar um novo
                fmt.Println("Token encontrado.")
                continue
            }
            token = true
            fmt.Println("Token gerado servidor A.")
        } else if (serv_local.qual_serv == "B") { // Tempo de espera para cada servidor gerar um novo token
            time.Sleep(300 * time.Millisecond) // Servidor B espera 0,3 segundos para gerar um novo token
            if token { // Se o token já foi gerado por outro servidor, o servidor atual não precisa gerar um novo
                fmt.Println("Token encontrado.")
                continue
            }
            token = true
            fmt.Println("Token gerado servidor B.")
        } else if (serv_local.qual_serv == "C") { // Tempo de espera para cada servidor gerar um novo token
            time.Sleep(600 * time.Millisecond) // Servidor C espera 0,6 segundos para gerar um novo token
            if token { // Se o token já foi gerado por outro servidor, o servidor atual não precisa gerar um novo
                fmt.Println("Token encontrado.")
                continue
            }
            token = true
            fmt.Println("Token gerado servidor C.")
        }
    }
}

// Função para enviar o token via HTTP para o próximo servidor
func envia_req_token(servidor string) bool {
    resposta, erro := http.Post(servidor + "/token", "application/json", nil) // Envia requisição POST para o servidor
    if erro != nil { // Se houver erro, retorna false
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
    for{
        if (token) { // Se o servidor possuir o token, irá mantê-lo por 1 seguno ou até terminar sua tarefa
            fmt.Println("Token em posse do servidor.")
            trava_token = true
            for (trava_token  || !tarefa_ok) { // Enquanto a trava estiver ativa ou a tarefa não estiver concluída, o servidor segura o token
                fmt.Println("Token em espera.")
                time.Sleep(1 * time.Second)
                trava_token = false
            }
            fmt.Println("Token em liberação.")

            // Com trava não ativa e tarefa concluída, o servidor passa o token para o próximo
            token = false // O servidor atual não possui mais o token
            
            for index, servidor := range serv_local.servidores { // O servidor atual tentará passar o token para os próximos servidores
                if (envia_req_token(servidor)) { // Se o token for enviado com sucesso, o servidor atual encerra o loop
                    fmt.Printf("Token enviado para o servidor %s\n", servidor)
                    break
                } else if (index == 0){ // Se o token não for enviado com sucesso, o servidor atual tentará enviar para o próximo servidor
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
        if erro != nil { // Se houver erro, o servidor não será atualizado
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
                    id: id,
                    nome: nome,
                    saldo: saldo,
                    id_eventos_criados: id_ev_cri,
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

                // Obtenção do resultado do evento
                resultado, ok := evento["resultado"].(string)
                if !ok {
                    fmt.Printf("Erro ao obter resultado de evento do servidor %s\n", servidor)
                }

                // Atualização do evento
                serv_local.eventos[id] = Evento{
                    id: id,
                    ativo: ativo,
                    id_criador: id_criador_int,
                    nome: nome,
                    descricao: descricao,
                    participantes: participantes,
                    palpite: palpite,
                    resultado: resultado,
                }
            }
        } else {
            fmt.Printf("Erro ao obter informações de eventos do servidor %s\n", servidor)
        }
        // Fim obtenção dos eventos
    }
}

// Função para definir as informações locais do servidor
func define_info() Infos_local{
	var qual_serv, tipo_serv string
    fmt.Printf("Tipo de servidor (local [LOC] ou laboratório [LAB]): ")
    fmt.Scan(&tipo_serv)
    tipo_serv = strings.ToUpper(tipo_serv)
	fmt.Printf("Servidor (A, B ou C): ")
    fmt.Scan(&qual_serv)
    qual_serv = strings.ToUpper(qual_serv)
    for{
        if tipo_serv == "LOC"{
            if qual_serv == "A"{
                servidores := []string{serv_loc_B, serv_loc_C}
                serv_local := Infos_local{
                    servidores: servidores, 
                    porta: ":8080", 
                    qual_serv: qual_serv,
                    clientes: make(map[int]Cliente),
                    eventos: make(map[int]Evento),
                }
                return serv_local
            } else if qual_serv == "B" {
                servidores := []string{serv_loc_C, serv_loc_A}
                serv_local := Infos_local{
                    servidores: servidores, 
                    porta: ":8081", 
                    qual_serv: qual_serv,
                    clientes: make(map[int]Cliente),
                    eventos: make(map[int]Evento),
                }
                return serv_local
            } else if qual_serv == "C" {
                servidores := []string{serv_loc_A, serv_loc_B}
                serv_local := Infos_local{
                    servidores: servidores, 
                    porta: ":8082", 
                    qual_serv: qual_serv,
                    clientes: make(map[int]Cliente),
                    eventos: make(map[int]Evento),
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
        } else if tipo_serv == "LAB"{
            if qual_serv == "A"{
                servidores := []string{serv_lab_B, serv_lab_C}
                serv_local := Infos_local{
                    servidores: servidores,
                    porta: ":8080",
                    qual_serv: qual_serv,
                    clientes: make(map[int]Cliente),
                    eventos: make(map[int]Evento),
                }
                return serv_local
            } else if qual_serv == "B"{
                servidores := []string{serv_lab_A, serv_lab_C}
                serv_local := Infos_local{
                    servidores: servidores, 
                    porta: ":8080", 
                    qual_serv: qual_serv,
                    clientes: make(map[int]Cliente),
                    eventos: make(map[int]Evento),
                }
                return serv_local
            } else if qual_serv == "C"{
                servidores := []string{serv_lab_A, serv_lab_B}
                serv_local := Infos_local{
                    servidores: servidores, 
                    porta: ":8080", 
                    qual_serv: qual_serv,
                    clientes: make(map[int]Cliente),
                    eventos: make(map[int]Evento),
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
func define_metodo_get(serv_local *Infos_local, serv *gin.Engine){
    // Método GET para retornar informações armazenadas pelo servidor, como clientes e eventos e as contagens de IDs
    serv.GET("/infos", func(c *gin.Context){
        // Criação de estrutura para informações de clientes
        clientes := make(map[int]map[string]interface{})
        for id, cliente := range serv_local.clientes{
            clientes[id] = map[string]interface{}{
                "nome": cliente.nome,
                "saldo": cliente.saldo,
                "id_eventos_criados": cliente.id_eventos_criados,
                "id_eventos_participados": cliente.id_eventos_participados,
            }
        }

        // Criação de estrutura para informações de eventos
        eventos := make(map[int]map[string]interface{})
        for id, evento := range serv_local.eventos{
            eventos[id] = map[string]interface{}{
                "ativo": evento.ativo,
                "id_criador": evento.id_criador,
                "nome": evento.nome,
                "descricao": evento.descricao,
                "participantes": evento.participantes,
                "palpite": evento.palpite,
                "resultado": evento.resultado,
            }
        }

        // Criação de estrutura para o retorno geral de informações
        infos := map[string]interface{}{
            "clientes": clientes,
            "eventos": eventos,
            "id_cont_cliente": id_cont_cliente,
            "id_cont_evento": id_cont_evento,
        }

        c.JSON(http.StatusOK, infos) // Retorna as informações armazenadas pelo servidor
    })
}

// Função para definir os métodos POST do servidor
func define_metodo_post(serv_local *Infos_local, serv *gin.Engine){
    // Método POST para recebimento do token
    serv.POST("/token", func(c *gin.Context){
        token = true // O servidor recebe o token
        c.JSON(http.StatusOK, gin.H{"message": "Token recebido com sucesso."}) // Retorna mensagem de sucesso
    })
}

// Função para definir os métodos PATCH do servidor
func define_metodo_patch(serv_local *Infos_local, serv *gin.Engine){}

// Função para definir o servidor com os métodos POST, GET e PATCH
func define_servidor(serv_local *Infos_local) *gin.Engine{
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

// tratameto de dados
func calcularPremio() {
    premiacao := 0.0 //Variável para armazenar o valor total das apostas
    total_ganhadores := 0.0 //Variável para armazenar o total apostado pelos ganhadores
    participantes := map[int]float64{ //Simulacao Mapa que associa o ID dos participantes aos seus respectivos valores de aposta
        01: 30,
        02: 25,
        03: 40,
    }
    palpite := map[int]string{ // Simulacao Mapa que associa o ID dos participantes aos seus palpites
        01: "A",
        02: "B",
        03: "A",
    }

    resultado := "A" // Simulacao String representando o resultado correto do palpite
    ganhadores := make(map[int]float64) // Mapa para armazenar os ganhadores e seus valores de aposta

    // calcula o valor total 
    for _, valor := range participantes{
        premiacao += valor
    }
    for ip, palpt := range palpite{
        if palpt == resultado{
            ganhadores[ip] = participantes[ip]
        }
    }
    //pagamento 
    for id, valor := range ganhadores{
        total_ganhadores += valor
        ganho := (participantes[id] / total_ganhadores) * premiacao
        //atribuir ganho ao saldo do Cliente
        fmt.Println(ganho) //para nn reclamar
    }
}


func main() {
    serv_local := define_info()
    servidor := define_servidor(&serv_local)
    atualiza_infos(&serv_local)
    go passa_token(&serv_local)
    go existe_token(&serv_local)
    servidor.Run(serv_local.porta)
}
