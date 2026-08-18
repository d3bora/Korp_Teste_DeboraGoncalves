// Pacote dto define os objetos de transferência de dados (Data Transfer Objects).
// Separa os contratos da API HTTP das entidades internas do domínio.
package dto

// CreateProductRequest é o corpo esperado na requisição de criação de produto.
type CreateProductRequest struct {
	Code        string `json:"code" binding:"required"`        // Código único obrigatório
	Description string `json:"description" binding:"required"` // Descrição obrigatória
	Balance     int    `json:"balance" binding:"min=0"`        // Saldo inicial (mínimo 0)
}

// UpdateProductRequest é o corpo esperado na atualização de produto.
// Todos os campos são opcionais (atualização parcial).
type UpdateProductRequest struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	Balance     int    `json:"balance" binding:"min=0"`
}

// DebitBalanceRequest é o corpo esperado na requisição de débito de saldo.
// Chamado internamente pelo billing-service.
type DebitBalanceRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"` // Quantidade a debitar (mínimo 1)
}

// ErrorResponse é o formato padrão de resposta de erro da API.
type ErrorResponse struct {
	Error string `json:"error"`
}
