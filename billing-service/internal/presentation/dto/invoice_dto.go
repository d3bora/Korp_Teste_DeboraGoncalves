// DTOs do microsserviço de Faturamento.
// Separam os contratos HTTP das entidades internas do domínio.
package dto

// InvoiceItemRequest representa um produto ao criar uma NF.
type InvoiceItemRequest struct {
	ProductID   uint   `json:"product_id" binding:"required"`
	ProductCode string `json:"product_code" binding:"required"`
	ProductName string `json:"product_name" binding:"required"`
	Quantity    int    `json:"quantity" binding:"required,min=1"`
}

// CreateInvoiceRequest é o body esperado na criação de uma NF.
type CreateInvoiceRequest struct {
	Items []InvoiceItemRequest `json:"items" binding:"required,min=1"`
}

// ErrorResponse é o formato padrão de resposta de erro.
type ErrorResponse struct {
	Error string `json:"error"`
}

// PrintResponse é a resposta após impressão bem-sucedida de uma NF.
type PrintResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}
