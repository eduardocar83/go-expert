package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultCEP     = "50610230"
	defaultTimeout = 1
)

type Endereco struct {
	CEP     string
	Estado  string
	Cidade  string
	Bairro  string
	Rua     string
	Servico string
}

type BrasilApiCepResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
}

type ViaCepResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func main() {
	chViaCep := make(chan *Endereco)
	chBrasilApi := make(chan *Endereco)

	go func() {
		endereco, err := consultarViaCep(defaultCEP)
		if err != nil {
			log.Println(err)
			return
		}
		chViaCep <- endereco
	}()

	go func() {
		endereco, err := consultarBrasilApiCep(defaultCEP)
		if err != nil {
			log.Println(err)
			return
		}
		chBrasilApi <- endereco
	}()

	select {
	case enderecoViaCep := <-chViaCep:
		log.Println(enderecoViaCep)
	case enderecoBrasilApi := <-chBrasilApi:
		log.Println(enderecoBrasilApi)
	case <-time.After(time.Second * defaultTimeout):
		log.Printf("Timeout: Após %d segundos não obtivemos nenhuma resposta com sucesso das apis de cep\n", defaultTimeout)
	}
}

func consultarViaCep(cep string) (*Endereco, error) {
	log.Println("Iniciando consulta Via Cep")

	reqURL := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", url.PathEscape(cep))

	resp, err := http.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("erro ao requisitar viaCep [%s]: %w", reqURL, err)
	}
	defer resp.Body.Close()

	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta de viaCep [%s]: %w", reqURL, err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("a requisicao a viaCep [%s] retornou o erro [%d] e mensagem [%s]", reqURL, resp.StatusCode, string(payload))
	}

	var response ViaCepResponse
	if err = json.Unmarshal(payload, &response); err != nil {
		return nil, fmt.Errorf("erro ao parsear resposta de viaCep [%s]: %w", reqURL, err)
	}

	log.Println("Finalizando consulta Via Cep")

	return &Endereco{
		CEP:     response.Cep,
		Estado:  response.Estado,
		Cidade:  response.Localidade,
		Bairro:  response.Bairro,
		Rua:     response.Logradouro,
		Servico: "ViaCep",
	}, nil
}

func consultarBrasilApiCep(cep string) (*Endereco, error) {
	log.Println("Iniciando consulta Brasil Api")

	reqURL := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", url.PathEscape(cep))

	resp, err := http.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("erro ao requisitar BrasilApi [%s]: %w", reqURL, err)
	}
	defer resp.Body.Close()

	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta de BrasilApi [%s]: %w", reqURL, err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("a requisicao a BrasilApi [%s] retornou o erro [%d] e a mensagem [%s]", reqURL, resp.StatusCode, string(payload))
	}

	var response BrasilApiCepResponse
	if err = json.Unmarshal(payload, &response); err != nil {
		return nil, fmt.Errorf("erro ao parsear resposta de BrasilApi [%s]: %w", reqURL, err)
	}

	log.Println("Finalizando consulta Brasil Api")

	return &Endereco{
		CEP:     response.Cep,
		Estado:  response.State,
		Cidade:  response.City,
		Bairro:  response.Neighborhood,
		Rua:     response.Street,
		Servico: "BrasilApi",
	}, nil
}
