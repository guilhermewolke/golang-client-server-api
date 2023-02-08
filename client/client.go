package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/guilhermewolke/golang-client-server-api/types"
)

// main Consome o endpoint do server, lê o JSON entregue pelo serviço e escreve o valor em um arquivo
func main() {
	log.Printf("(client) main.CotacaoHandler - Início do método")
	ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(fmt.Sprintf("Erro do request: %v", err))
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(fmt.Sprintf("Erro do do: %v", err))
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		panic(fmt.Sprintf("Erro do body: %v", err))
	}

	var cotacao types.CotacaoResponse

	if err = json.Unmarshal(body, &cotacao); err != nil {
		panic(fmt.Sprintf("Erro do Unmarshal: %v", err))
	}

	log.Printf("(client) main.CotacaoHandler - cotacao: %#v", cotacao)
	// escrever no arquivo cotacao.txt o valor da cotacao
	file, err := os.Create("cotacao.txt")

	if err != nil {
		panic(fmt.Sprintf("Erro do os.Create: %v", err))
	}

	defer file.Close()
	data := []byte(cotacao.BID)
	if _, err = file.Write(data); err != nil {
		panic(fmt.Sprintf("Erro do file.Write: %v", err))
	}
	log.Printf("(client) main.CotacaoHandler - Fim do método")
}
