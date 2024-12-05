package main

import (
	"fmt"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var serv_loc_A string = "http://localhost:8080"
var serv_loc_B string = "http://localhost:8081"
var serv_loc_C string = "http://localhost:8082"

var serv_lab_A string = "http://172.16.103.14:8080"
var serv_lab_B string = "http://172.16.103.13:8080"
var serv_lab_C string = "http://172.16.103.12:8080"

type Cliente struct {
	id   int
	nome string
    saldo int
}

type Evento struct {
    id int
    ativo bool
    id_criador int
    nome string
    descricao string
}

type Infos_local struct {
	servidores []string
	porta      string
}

func define_info() {
	var qual_serv string
	fmt.Printf("Servidor (A, B ou C): ")
    fmt.Scan(&qual_serv)
    for{
        switch strings.ToUpper(qual_serv) {
        case "A":
        case "B":
        case "C":
        default:
            fmt.Println("Servidor inválido. Digite novamente.")
            fmt.Printf("Servidor (A, B ou C): ")
            fmt.Scan(&qual_serv)
            continue
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

func main() {

}
