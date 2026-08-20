// Pacote client contém o HTTP client para comunicação com o stock-service.
// Implementa tratamento de falha: se o stock-service estiver indisponível,
// retorna um erro claro que será repassado ao usuário pelo frontend.
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// StockClient encapsula chamadas HTTP ao microsserviço de Estoque.
type StockClient struct {
	baseURL    string
	httpClient *http.Client
}

// DebitRequest é o payload enviado ao stock-service para débito de saldo.
type DebitRequest struct {
	Quantity int `json:"quantity"`
}

// ProductResponse é a resposta do stock-service após débito.
type ProductResponse struct {
	ID      uint   `json:"id"`
	Code    string `json:"code"`
	Balance int    `json:"balance"`
}

// NewStockClient cria o client com URL e timeout configuráveis via ambiente.
func NewStockClient() *StockClient {
	baseURL := os.Getenv("STOCK_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8081" // padrão para desenvolvimento local
	}

	return &StockClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second, // Timeout para evitar travamentos
		},
	}
}

// DebitProductBalance debita a quantidade do saldo de um produto no stock-service.
// Retorna erro detalhado se o serviço estiver indisponível ou o saldo for insuficiente.
func (c *StockClient) DebitProductBalance(productID uint, quantity int) error {
	url := fmt.Sprintf("%s/products/%d/debit", c.baseURL, productID)

	payload, err := json.Marshal(DebitRequest{Quantity: quantity})
	if err != nil {
		return fmt.Errorf("erro ao serializar requisição: %w", err)
	}

	// Usa PATCH para corresponder à rota PATCH /:id/debit do stock-service
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("erro ao criar requisição: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		// ⚠️ Cenário de falha: stock-service indisponível
		return fmt.Errorf("serviço de estoque indisponível — tente novamente em instantes: %w", err)
	}
	defer resp.Body.Close()

	// Trata erros de negócio retornados pelo stock-service (ex: saldo insuficiente)
	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Error string `json:"error"`
		}
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("erro no serviço de estoque: %s", errResp.Error)
	}

	return nil
}

// HealthCheck verifica se o stock-service está acessível.
// Usado para diagnóstico e feedback ao usuário.
func (c *StockClient) HealthCheck() error {
	url := fmt.Sprintf("%s/health", c.baseURL)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("stock-service inacessível: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("stock-service retornou status %d", resp.StatusCode)
	}
	return nil
}
