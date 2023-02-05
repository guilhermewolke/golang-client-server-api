package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/guilhermewolke/golang-client-server-api/types"
)

func main() {
	http.HandleFunc("/cotacao", CotacaoHandler)
	http.ListenAndServe(":8080", nil)
}

// CotacaoHandler Realiza uma requisição HTTP para consumir a API contendo o câmbio de Dolar e Real, no endereço 'https://economia.awesomeapi.com.br/json/last/USD-BRL'
// e retorna em seguida o resultado para o client, em formato JSON.

// Além disso, o serviço deve ser capaz de realizar o registro em banco de dados SQLite de cada cotação recebida, respeitando limites de tempo de timeout pré-estabelecidos pelo desafio:
//  - 200ms para a chamada da API
//  - 100ms para persistir o resultado no banco de dados SQLite

func CotacaoHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("main.CotacaoHandler - Início do Método")
	// Passos a seguir
	// 1) Realizar a requisição no serviço externo de cotação de Dolar

	cotacao, err := getCotacao()
	if err != nil {
		panic(err)
	}

	log.Printf("main.CotacaoHandler - cotacao convertido em objeto: %#v", cotacao)

	// 2) Registrar no banco de dados SQLite
	// err = putCotacao(cotacao)

	// if err != nil {
	// 	panic(err)
	// }

	log.Println("main.CotacaoHandler - Fim do Método")
}

func getCotacao() (map[string]types.CotacaoDTO, error) {
	log.Printf("main.getCotacao -Início do método")
	var cotacao map[string]types.CotacaoDTO

	ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Millisecond)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, types.URL_API, nil)
	if err != nil {
		return cotacao, err
	}

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return cotacao, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return cotacao, err
	}

	log.Printf("main.getCotacao - body retornado: %#v", string(body))

	if err = json.Unmarshal(body, &cotacao); err != nil {
		return cotacao, err
	}

	return cotacao, nil
}
