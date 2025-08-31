package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type cotacao struct {
	Bid string `json:"bid"`
}

func RequisitarServidorLocal() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()

	cotacao, err := obterCotacaoDeApiLocal(ctx)
	if err != nil {
		log.Printf("erro ao obter cotacao de api local: %v", err)
		return
	}

	err = salvarCotacaoEmArquivo(cotacao)
	if err != nil {
		log.Printf("erro ao salvar cotacao em arquivo: %v", err)
		return
	}
}

func obterCotacaoDeApiLocal(ctx context.Context) (*cotacao, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao requisitar servidor: %w", err)
	}
	defer resp.Body.Close()

	var cotacao cotacao

	err = json.NewDecoder(resp.Body).Decode(&cotacao)
	if err != nil {
		return nil, fmt.Errorf("erro ao parsear resposta: %w", err)
	}

	return &cotacao, nil
}

func salvarCotacaoEmArquivo(c *cotacao) error {
	file, err := os.OpenFile("./cotacao.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo de cotação: %w", err)
	}
	defer file.Close()

	registro := fmt.Sprintf("Dólar: %s\n", c.Bid)

	_, err = file.WriteString(registro)
	if err != nil {
		return fmt.Errorf("erro ao salvar cotação %s: %w", c.Bid, err)
	}

	return nil
}
