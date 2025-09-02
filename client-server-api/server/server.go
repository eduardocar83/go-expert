package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type USDBRL struct {
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

type USDBRLWrapper struct {
	USDBRL USDBRL `json:"USDBRL"`
}

type Cotacao struct {
	Bid string `json:"bid"`
}

var DB *sql.DB

func InicializarServidor() {
	initDB()
	defer DB.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", cotacaoHandler)
	log.Fatal(http.ListenAndServe("localhost:8080", mux))
}

func initDB() {
	var err error

	DB, err = sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := `
		CREATE TABLE IF NOT EXISTS cotacao (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		created_at TEXT DEFAULT CURRENT_TIMESTAMP,
		bid TEXT
		);`

	_, err = DB.Exec(sqlStmt)

	if err != nil {
		log.Fatalf("Erro ao criar a tabela <cotacao>: %v ", err)
	}
}

func cotacaoHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctxAPI, cancelAPI := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancelAPI()

	cotacao, err := obterCotacaoDeApiExterna(ctxAPI)
	if err != nil {
		log.Printf("Erro ao requisitar api externa de cotação de dolar: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctxDB, cancelDB := context.WithTimeout(ctx, time.Millisecond*10)
	defer cancelDB()

	err = salvarCotacao(ctxDB, cotacao)
	if err != nil {
		log.Printf("Erro ao salvar cotação de dolar: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cotacao)
}

func obterCotacaoDeApiExterna(ctx context.Context) (*Cotacao, error) {
	url := "https://economia.awesomeapi.com.br/json/last/USD-BRL"

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar request para [%s]: %w", url, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao requisitar endpoint [%s]: %w", url, err)
	}
	defer resp.Body.Close()

	var usdbrlWrapper USDBRLWrapper

	if err := json.NewDecoder(resp.Body).Decode(&usdbrlWrapper); err != nil {
		return nil, fmt.Errorf("erro ao parsear resposta do endpoint [%s]: %w", url, err)
	}

	return &Cotacao{Bid: usdbrlWrapper.USDBRL.Bid}, nil
}

func salvarCotacao(ctx context.Context, cotacao *Cotacao) error {
	_, err := DB.ExecContext(ctx, "insert into cotacao(bid) values(?)", cotacao.Bid)
	if err != nil {
		return fmt.Errorf("erro ao inserir cotacao %s: %w", cotacao.Bid, err)
	}
	return nil
}
