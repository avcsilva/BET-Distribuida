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
var token bool = false
var trava_token bool = false
var tarefa_ok bool = true

// Variáveis para lógica do sistema
var id_cont_cliente int = 1
var id_cont_evento int = 1

// Estruturas de dados
// Estrutura de dados para o Cliente
type Cliente struct {
	id   int
	nome string
    saldo float64
}

// Estrutura de dados para o Evento
type Evento struct {
    id              int
    ativo           bool
    id_criador      int
    nome            string
    descricao       string
    participantes 	map[int]float64
    palpite         map[int]string
    resultado       string
}

// Estrutura de dados para as informações locais do servidor em questão
type Infos_local struct {
	servidores []string
	porta      string
    qual_serv string
}

// Fim das estruturas de dados

// Funções para o Token Ring
// Função para validar existência de token no sistema com 3 servidores
func existe_token(serv_local *Infos_local) {
    // Verifica se a variável token se torna verdadeira dentro de um tempo limite
    for i := 0; i < 3; i++ {
        if token == true {
            return
        }
        // Conta 1 segundo
        time.Sleep(1 * time.Second)
    } // Caso não haja token em 3 segundos, o servidor atual gerará um novo token

    if (serv_local.qual_serv == "A") {
        time.Sleep(100 * time.Millisecond)
        if token {
            return
        }
        token = true
    } else if (serv_local.qual_serv == "B") {
        time.Sleep(300 * time.Millisecond)
        if token {
            return
        }
        token = true
    } else if (serv_local.qual_serv == "C") {
        time.Sleep(600 * time.Millisecond)
        if token { 
            return
        }
        token = true
    }
    return
}

// Função para enviar o token via HTTP para o próximo servidor
func envia_req_token(servidor string) bool {
    resposta, erro := http.Post(servidor + "/token", "application/json", nil)
    if erro != nil {
        return false
    }

    if resposta.StatusCode == 200 {
        return true
    } else {
        return false
    }
}

// Função para segurar e passar o token para o próximo servidor
func passa_token(serv_local *Infos_local) {
    for{
        if (token) {
            trava_token = true
            for (trava_token  || !tarefa_ok) {
                time.Sleep(1 * time.Second)
                trava_token = false
            }

            token = false
            
            for index, servidor := range serv_local.servidores {
                if (envia_req_token(servidor)) {
                    fmt.Printf("Token enviado para o servidor %s\n", servidor)
                    break
                } else if (index == 0){
                    fmt.Printf("Erro ao enviar token para o servidor %s\n. Tentando enviar ao próximo.", servidor)
                } else {
                    fmt.Printf("Erro ao enviar token para o servidor %s\n", servidor)
                    fmt.Printf("Token retornado ao servidor %s\n", serv_local.qual_serv)
                    token = true
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
                serv_local := Infos_local{servidores, "8080", qual_serv}
                return serv_local
            } else if qual_serv == "B" {
                servidores := []string{serv_loc_C, serv_loc_A}
                serv_local := Infos_local{servidores, "8081", qual_serv}
                return serv_local
            } else if qual_serv == "C" {
                servidores := []string{serv_loc_A, serv_loc_B}
                serv_local := Infos_local{servidores, "8082", qual_serv}
                return serv_local
            } else {
                fmt.Printf("Servidor inválido.")
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
                serv_local := Infos_local{servidores, "8080", qual_serv}
                return serv_local
            } else if qual_serv == "B"{
                servidores := []string{serv_lab_A, serv_lab_C}
                serv_local := Infos_local{servidores, "8080", qual_serv}
                return serv_local
            } else if qual_serv == "C"{
                servidores := []string{serv_lab_A, serv_lab_B}
                serv_local := Infos_local{servidores, "8080", qual_serv}
                return serv_local
            } else {
                fmt.Printf("Servidor inválido.")
                fmt.Printf("Tipo de servidor (local [LOC] ou laboratório [LAB]): ")
                fmt.Scan(&tipo_serv)
                tipo_serv = strings.ToUpper(tipo_serv)
                fmt.Printf("Servidor (A, B ou C): ")
                fmt.Scan(&qual_serv)
                qual_serv = strings.ToUpper(qual_serv)
            }
        } else {
            fmt.Printf("Tipo de servidor inválido.")
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
func define_metodo_get(serv_local *Infos_local, serv *gin.Engine, id_cont *int){}

// Função para definir os métodos POST do servidor
func define_metodo_post(serv_local *Infos_local, serv *gin.Engine, id_cont *int){}

// Função para definir os métodos PATCH do servidor
func define_metodo_patch(serv_local *Infos_local, serv *gin.Engine){}

// Função para definir o servidor com os métodos POST, GET e PATCH
func define_servidor(serv_local *Infos_local, id_cont *int) *gin.Engine{
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
	define_metodo_get(serv_local, r, id_cont)

	// Define os métodos POST
	define_metodo_post(serv_local, r, id_cont)

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
}
