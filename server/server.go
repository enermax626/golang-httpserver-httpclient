package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

var db *sql.DB
var client = http.DefaultClient
var cotacaoURL = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

type BidResponse struct {
	Bid string `json:"bid"`
}

type USDBRLResponse struct {
	USDBRL Currency `json:"USDBRL"`
}

type Currency struct {
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

func init() {
	log.Println("Initializing database connection...")
	db, err := sql.Open("sqlite", "./cotacoes.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// Test the database connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS cotacoes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		bid TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	)`
	_, err = db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	log.Println("Database initialized successfully.")
}

func StartServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", getBidHandler)

	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func getBidHandler(rw http.ResponseWriter, request *http.Request) {
	bidResponse, err := getBid()
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	marshalResponse, err := json.Marshal(bidResponse)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = rw.Write(marshalResponse)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getBid() (*BidResponse, error) {
	data, err := fetchData()
	if err != nil {
		return nil, err
	}
	bidResponse := BidResponse{
		Bid: data.USDBRL.Bid,
	}
	err = persistBid(bidResponse)
	if err != nil {
		return nil, err
	}
	return &bidResponse, nil
}

func persistBid(bidResponse BidResponse) error {
	ctxDbPersistence, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	query := `INSERT INTO cotacoes (bid) VALUES (?)`
	_, err := db.ExecContext(ctxDbPersistence, query, bidResponse.Bid)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func fetchData() (*USDBRLResponse, error) {
	ctxHttpRequest, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	request, err := http.NewRequestWithContext(ctxHttpRequest, "GET", cotacaoURL, nil)
	if err != nil {
		return nil, err
	}
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var responseModel USDBRLResponse
	err = json.Unmarshal(body, &responseModel)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &responseModel, nil
}
