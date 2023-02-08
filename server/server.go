package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/guilhermewolke/golang-client-server-api/types"
	_ "github.com/mattn/go-sqlite3"
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
	// Primeiro, cria-se o contexto com 200ms para o consumo da API...
	ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Millisecond)
	defer cancel()

	// Passos a seguir
	// 1) Realizar a requisição no serviço externo de cotação de Dolar

	cotacao, err := getCotacao(ctx)
	if err != nil {
		panic(fmt.Sprintf("Erro do getCotacao: %v", err))
	}

	log.Printf("main.CotacaoHandler - cotacao convertido em objeto: %#v", cotacao)

	// 2) Registrar no banco de dados SQLite
	// Alterando o timeout do contexto para o tempo limite
	ctx = nil
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Nanosecond)

	lastID, err := putCotacao(ctx, cotacao)

	if err != nil {
		panic(fmt.Sprintf("Erro do put cotacao: %v", err))
	}

	log.Printf("main.CotacaoHandler - Último ID inserido: %d", lastID)

	// 3 Retornando o json de resposta para o client
	payload := types.CotacaoResponse{BID: cotacao["USDBRL"].BID}

	json.NewEncoder(w).Encode(payload)
	log.Println("main.CotacaoHandler - Fim do Método")
}

func getCotacao(c context.Context) (map[string]types.CotacaoDataDTO, error) {
	log.Printf("main.getCotacao -Início do método")
	var cotacao map[string]types.CotacaoDataDTO

	request, err := http.NewRequestWithContext(c, http.MethodGet, types.URL_API, nil)
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

	if err = json.Unmarshal(body, &cotacao); err != nil {
		return cotacao, err
	}

	return cotacao, nil
}

func putCotacao(c context.Context, cotacao map[string]types.CotacaoDataDTO) (int64, error) {
	log.Printf("main.putCotacao -Início do método")
	obj := cotacao["USDBRL"]
	// Conectar-se ao banco
	db, err := sql.Open("sqlite3", "cotacao.db")
	if err != nil {
		return int64(0), err
	}

	defer db.Close()

	q := `
		INSERT INTO
			cotacoes(
				code,
				codein,
				name,
				high,
				low,
				varbid,
				pctchange,
				bid,
				ask,
				timestamp,
				create_date
			)
			VALUES (
				?, 
				?, 
				?, 
				?, 
				?, 
				?, 
				?, 
				?, 
				?, 
				?, 
				?
			);
	`

	select {
	case <-c.Done():
		return int64(0), errors.New("Timeout do contexto excedido")
	}

	result, err := db.Exec(q,
		obj.Code,
		obj.CodeIN,
		obj.Name,
		obj.High,
		obj.Low,
		obj.VarBid,
		obj.PctChange,
		obj.BID,
		obj.Ask,
		obj.Timestamp,
		obj.CreateDate,
	)

	if err != nil {
		return int64(0), err
	}

	// Inserir o registro na tabela
	lastID, err := result.LastInsertId()

	if err != nil {
		return int64(0), err
	}

	return lastID, nil
}
