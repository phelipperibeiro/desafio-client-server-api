package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	uuid "github.com/satori/go.uuid"
)

type USDBRL struct {
	ID         string `gorm:"primaryKey"`
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

type ExchangeRate struct {
	Bid string `json:"bid"`
}

var Db *sql.DB

func main() {
	Db, _ = sql.Open("sqlite3", "db.sqlite")
	createTable(Db)

	log.Println("Iniciando servidor")

	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", handler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func handler(responseWriter http.ResponseWriter, request *http.Request) {

	log.Println("consulta iniciada")

	defer log.Println("consulta finalizada")

	if request.URL.Path != "/cotacao" {
		http.NotFound(responseWriter, request)
		return
	}

	usdbrl, err := getQuotation()
	if err != nil {
		http.Error(responseWriter, "Erro ao obter cotação", http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(responseWriter).Encode(ExchangeRate{Bid: usdbrl.Bid})
}

func getQuotation() (*USDBRL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var usdbrlMap map[string]USDBRL
	err = json.Unmarshal(body, &usdbrlMap)
	if err != nil {
		return nil, err
	}

	log.Println("JSON formatado:")
	log.Println(string(body))

	usdbrl, ok := usdbrlMap["USDBRL"]
	if !ok {
		return nil, fmt.Errorf("Chave USDBRL não encontrada no JSON")
	}

	usdbrl.ID = uuid.NewV4().String()

	err = createUSDBRL(Db, &usdbrl)
	return &usdbrl, err
}

func createTable(db *sql.DB) {
	table := `CREATE TABLE IF NOT EXISTS USDBRL (
		id STRING PRIMARY KEY,
		code STRING,
		codein STRING,
		name STRING,
		high STRING,
		low STRING,
		varBid STRING,
		pctChange STRING,
		bid STRING,
		ask STRING,
		timestamp STRING,
		create_date TIMESTAMP
	);`

	_, err := db.Exec(table)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func createUSDBRL(db *sql.DB, usdbrl *USDBRL) error {
	log.Println("inserindo dados no DB")
	defer log.Println("dados inserido!!!")

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	insert := `INSERT INTO USDBRL (id, code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP);`

	_, err := db.ExecContext(
		ctx,
		insert,
		usdbrl.ID,
		usdbrl.Code,
		usdbrl.Codein,
		usdbrl.Name,
		usdbrl.High,
		usdbrl.Low,
		usdbrl.VarBid,
		usdbrl.PctChange,
		usdbrl.Bid,
		usdbrl.Ask,
		usdbrl.Timestamp,
	)
	return err
}
