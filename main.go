package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	apiAuthToken = "seu_token_de_autenticacao_aqui"
)

func main() {
	router := gin.Default()

	// Variável para armazenar o resultado anterior
	var resultadoAnterior string

	// Função para executar uma URL específica
	executarURL := func(url string) {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Erro ao executar a URL:", err)
			return
		}
		defer resp.Body.Close()

		// Lógica adicional, se necessário, para lidar com a resposta da URL
		fmt.Println("URL executada com sucesso:", url)
	}

	// Função para realizar a consulta à API com autenticação via token
	consultaAPI := func() {
		// Substitua a URL pela API que você deseja consultar
		apiURL := "https://api.example.com/data"
		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			fmt.Println("Erro ao criar a solicitação à API:", err)
			return
		}

		// Adiciona o token de autenticação ao cabeçalho da solicitação
		req.Header.Set("Authorization", "Bearer "+apiAuthToken)

		// Realiza a solicitação
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Erro ao consultar a API:", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Erro ao ler a resposta da API:", err)
			return
		}

		// Converte o corpo da resposta para string
		resultadoAtual := string(body)

		// Verifica se houve uma mudança no resultado
		if resultadoAtual != resultadoAnterior {
			// Atualiza o resultado anterior
			resultadoAnterior = resultadoAtual

			// Lógica para executar diferentes URLs com base na resposta
			switch resultadoAtual {
			case "A":
				executarURL("http://192.168.1.80/cron/start")
			case "B":
				executarURL("http://192.168.1.80/cron/stop")
			default:
				fmt.Println("Resposta desconhecida da API:", resultadoAtual)
			}
		}
	}

	// Rota para iniciar a consulta à API
	router.GET("/consultar-api", func(c *gin.Context) {
		consultaAPI()
		c.JSON(http.StatusOK, gin.H{"message": "Consulta bem-sucedida"})
	})

	// Rotina para realizar a consulta a cada segundo
	go func() {
		for {
			consultaAPI()
			time.Sleep(1 * time.Second)
		}
	}()

	// Inicia o servidor
	router.Run(":8080")
}
