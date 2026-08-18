// Configuração das rotas e middlewares do billing-service.
package router

import (
	"github.com/deboragoncalves/korp/billing-service/internal/presentation/handler"
	"github.com/gin-gonic/gin"
)

// Setup monta todas as rotas da API de Faturamento.
func Setup(invoiceHandler *handler.InvoiceHandler) *gin.Engine {
	r := gin.New()

	// ─── Middlewares globais ────────────────────────────────────────────────
	r.Use(gin.Logger())
	r.Use(gin.Recovery()) // Evita crash em panics inesperados

	// ─── Health Check ───────────────────────────────────────────────────────
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "billing-service"})
	})

	// ─── Rotas de Notas Fiscais ─────────────────────────────────────────────
	invoices := r.Group("/invoices")
	{
		invoices.GET("", invoiceHandler.GetAll)          // Listar todas as NFs
		invoices.GET("/:id", invoiceHandler.GetByID)     // Buscar NF por ID
		invoices.POST("", invoiceHandler.Create)         // Criar nova NF

		// Endpoint de impressão: debita estoque e fecha a NF
		invoices.POST("/:id/print", invoiceHandler.Print)
	}

	return r
}
