package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	apiAuthToken = "26c465ff73fcbbfc3e79f528e9cc59fde54081ba"
)

type Response struct {
	CodReuniao                   int    `json:"codReuniao"`
	CodReuniaoPrincipal          int    `json:"codReuniaoPrincipal"`
	TxtTituloReuniao             string `json:"txtTituloReuniao"`
	TxtSiglaOrgao                string `json:"txtSiglaOrgao"`
	TxtApelido                   string `json:"txtApelido"`
	TxtNomeOrgao                 string `json:"txtNomeOrgao"`
	CodEstadoReuniao             int    `json:"codEstadoReuniao"`
	TxtTipoReuniao               string `json:"txtTipoReuniao"`
	TxtObjeto                    string `json:"txtObjeto"`
	TxtLocal                     string `json:"txtLocal"`
	BolReuniaoConjunta           bool   `json:"bolReuniaoConjunta"`
	BolHabilitarEventoInterativo bool   `json:"bolHabilitarEventoInterativo"`
	IDYoutube                    string `json:"idYoutube"`
	CodEstadoTransmissaoYoutube  int    `json:"codEstadoTransmissaoYoutube"`
	DatReuniaoString             string `json:"datReuniaoString"`
}

func main() {
	router := gin.Default()

	// Variável para armazenar o resultado anterior
	var resultadoAnterior bool

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
		apiURL := "https://sapl.sapezal.mt.leg.br/api/sessao-plenaria/21/"
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

		var myStoredVariable Response
		json.Unmarshal(body, &myStoredVariable)
		// return myStoredVariable.BolHabilitarEventoInterativo
		// Converte o corpo da resposta para string
		resultadoAtual := myStoredVariable.BolHabilitarEventoInterativo

		fmt.Println("Resposta da API:", resultadoAtual)
		// Verifica se houve uma mudança no resultado
		if resultadoAtual != resultadoAnterior {
			// Atualiza o resultado anterior
			resultadoAnterior = resultadoAtual

			// Lógica para executar diferentes URLs com base na resposta
			switch resultadoAtual {
			case true:
				executarURL("http://192.168.1.80/Bcron.htm?C1=2&f2=on&n4=0&n5=8&n6=0&Fc=CONTAR&f4=on")
			case false:
				executarURL("http://192.168.1.80/Bcron.htm?C1=2&f2=on&n4=0&n5=8&n6=0&Fp=PARAR&f4=on")
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
