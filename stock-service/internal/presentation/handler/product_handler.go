// Pacote handler contém os controllers HTTP do microsserviço de Estoque.
// Cada método mapeia para um endpoint REST e delega a lógica ao service.
package handler

import (
	"net/http"
	"strconv"

	"github.com/deboragoncalves/korp/stock-service/internal/application/service"
	"github.com/deboragoncalves/korp/stock-service/internal/domain/entity"
	"github.com/deboragoncalves/korp/stock-service/internal/presentation/dto"
	"github.com/gin-gonic/gin"
)

// ProductHandler agrupa os handlers HTTP relacionados a produtos.
type ProductHandler struct {
	service *service.ProductService
}

// NewProductHandler cria o handler via Dependency Injection.
func NewProductHandler(s *service.ProductService) *ProductHandler {
	return &ProductHandler{service: s}
}

// GetAll godoc
// GET /products
// Retorna lista de todos os produtos cadastrados.
func (h *ProductHandler) GetAll(c *gin.Context) {
	products, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

// GetByID godoc
// GET /products/:id
// Retorna um produto pelo ID.
func (h *ProductHandler) GetByID(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "ID inválido"})
		return
	}

	product, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

// Create godoc
// POST /products
// Cria um novo produto com os dados fornecidos no body.
func (h *ProductHandler) Create(c *gin.Context) {
	var req dto.CreateProductRequest

	// binding:"required" valida automaticamente os campos obrigatórios
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	product := &entity.Product{
		Code:        req.Code,
		Description: req.Description,
		Balance:     req.Balance,
	}

	if err := h.service.Create(product); err != nil {
		c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// Update godoc
// PUT /products/:id
// Atualiza os dados de um produto existente.
func (h *ProductHandler) Update(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "ID inválido"})
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	updated := &entity.Product{
		Code:        req.Code,
		Description: req.Description,
		Balance:     req.Balance,
	}

	product, err := h.service.Update(id, updated)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

// Delete godoc
// DELETE /products/:id
// Remove um produto pelo ID.
func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "ID inválido"})
		return
	}

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// DebitBalance godoc
// PATCH /products/:id/debit
// Debita a quantidade informada do saldo do produto.
// Endpoint chamado internamente pelo billing-service durante impressão de NF.
func (h *ProductHandler) DebitBalance(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "ID inválido"})
		return
	}

	var req dto.DebitBalanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	product, err := h.service.DebitBalance(id, req.Quantity)
	if err != nil {
		// 422 Unprocessable Entity: saldo insuficiente ou regra de negócio violada
		c.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

// parseUintParam extrai e converte um parâmetro de rota para uint.
func parseUintParam(c *gin.Context, param string) (uint, error) {
	val, err := strconv.ParseUint(c.Param(param), 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(val), nil
}
