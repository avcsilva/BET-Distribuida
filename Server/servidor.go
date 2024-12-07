package main

import (
	"fmt"
	"net/http"
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
                continue
            }
            // Conta 1 segundo
            time.Sleep(1 * time.Second)
        } // Caso não haja token em 3 segundos, o servidor atual gerará um novo token

        if (serv_local.qual_serv == "A") { // Tempo de espera para cada servidor gerar um novo token
            time.Sleep(100 * time.Millisecond) // Servidor A espera apenas 0,1 segundos para gerar um novo token
            if token { // Se o token já foi gerado por outro servidor, o servidor atual não precisa gerar um novo
                continue
            }
            token = true
        } else if (serv_local.qual_serv == "B") { // Tempo de espera para cada servidor gerar um novo token
            time.Sleep(300 * time.Millisecond) // Servidor B espera 0,3 segundos para gerar um novo token
            if token { // Se o token já foi gerado por outro servidor, o servidor atual não precisa gerar um novo
                continue
            }
            token = true
        } else if (serv_local.qual_serv == "C") { // Tempo de espera para cada servidor gerar um novo token
            time.Sleep(600 * time.Millisecond) // Servidor C espera 0,6 segundos para gerar um novo token
            if token { // Se o token já foi gerado por outro servidor, o servidor atual não precisa gerar um novo
                continue
            }
            token = true
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
            trava_token = true
            for (trava_token  || !tarefa_ok) { // Enquanto a trava estiver ativa ou a tarefa não estiver concluída, o servidor segura o token
                time.Sleep(1 * time.Second)
                trava_token = false
            }

            // Com trava não ativa e tarefa concluída, o servidor passa o token para o próximo
            token = false // O servidor atual não possui mais o token
            
            for index, servidor := range serv_local.servidores { // O servidor atual tentará passar o token para os próximos servidores
                if (envia_req_token(servidor)) { // Se o token for enviado com sucesso, o servidor atual encerra o loop
                    fmt.Printf("Token enviado para o servidor %s\n", servidor)
                    break
                } else if (index == 0){ // Se o token não for enviado com sucesso, o servidor atual tentará enviar para o próximo servidor
                    fmt.Printf("Erro ao enviar token para o servidor %s\n. Tentando enviar ao próximo.", servidor)
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
func define_metodo_get(serv_local *Infos_local, serv *gin.Engine){}

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
    go passa_token(&serv_local)
    go existe_token(&serv_local)
    servidor.Run(serv_local.porta)
}
