// Pacote handler contém os controllers HTTP do microsserviço de Faturamento.
package handler

import (
	"net/http"
	"strconv"

	"github.com/deboragoncalves/korp/billing-service/internal/application/service"
	"github.com/deboragoncalves/korp/billing-service/internal/domain/entity"
	"github.com/deboragoncalves/korp/billing-service/internal/presentation/dto"
	"github.com/gin-gonic/gin"
)

// InvoiceHandler agrupa os controllers HTTP de Notas Fiscais.
type InvoiceHandler struct {
	service *service.InvoiceService
}

// NewInvoiceHandler cria o handler via Dependency Injection.
func NewInvoiceHandler(s *service.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{service: s}
}

// GetAll godoc
// GET /invoices
// Retorna todas as notas fiscais com seus itens.
func (h *InvoiceHandler) GetAll(c *gin.Context) {
	invoices, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, invoices)
}

// GetByID godoc
// GET /invoices/:id
// Retorna uma nota fiscal pelo ID.
func (h *InvoiceHandler) GetByID(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "ID inválido"})
		return
	}

	invoice, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, invoice)
}

// Create godoc
// POST /invoices
// Cria uma nova NF com numeração sequencial e status Aberta.
func (h *InvoiceHandler) Create(c *gin.Context) {
	var req dto.CreateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	// Mapeia DTO → entidade do domínio
	items := make([]entity.InvoiceItem, len(req.Items))
	for i, itemReq := range req.Items {
		items[i] = entity.InvoiceItem{
			ProductID:   itemReq.ProductID,
			ProductCode: itemReq.ProductCode,
			ProductName: itemReq.ProductName,
			Quantity:    itemReq.Quantity,
		}
	}

	invoice := &entity.Invoice{Items: items}

	if err := h.service.Create(invoice); err != nil {
		c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, invoice)
}

// Print godoc
// POST /invoices/:id/print
// Imprime a NF: debita estoque e muda status para Fechada.
// Retorna erro 422 se a NF não estiver Aberta ou o estoque for insuficiente.
// Retorna erro 503 se o stock-service estiver indisponível.
func (h *InvoiceHandler) Print(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "ID inválido"})
		return
	}

	invoice, err := h.service.PrintInvoice(id)
	if err != nil {
		// Distingue falha de serviço externo (503) de erro de negócio (422)
		statusCode := http.StatusUnprocessableEntity
		if isServiceUnavailable(err.Error()) {
			statusCode = http.StatusServiceUnavailable
		}
		c.JSON(statusCode, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.PrintResponse{
		Message: "Nota fiscal impressa com sucesso",
		Status:  string(invoice.Status),
	})
}

// isServiceUnavailable identifica erros de conectividade com serviços externos.
func isServiceUnavailable(errMsg string) bool {
	return len(errMsg) > 0 && (contains(errMsg, "indisponível") || contains(errMsg, "inacessível"))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && stringContains(s, substr))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func parseUintParam(c *gin.Context, param string) (uint, error) {
	val, err := strconv.ParseUint(c.Param(param), 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(val), nil
}
