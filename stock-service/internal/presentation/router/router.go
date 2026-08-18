// Pacote router configura as rotas e middlewares do servidor HTTP (Gin).
package router

import (
	"github.com/deboragoncalves/korp/stock-service/internal/presentation/handler"
	"github.com/gin-gonic/gin"
)

// Setup configura todas as rotas da API de Estoque e retorna o engine do Gin.
func Setup(productHandler *handler.ProductHandler) *gin.Engine {
	r := gin.New()

	// ─── Middlewares globais ────────────────────────────────────────────────
	r.Use(gin.Logger())   // Log de todas as requisições
	r.Use(gin.Recovery()) // Recupera de panics sem derrubar o serviço

	// ─── Health Check ───────────────────────────────────────────────────────
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "stock-service"})
	})

	// ─── Rotas de Produtos ──────────────────────────────────────────────────
	// Agrupadas sob /products para organização e versionamento futuro
	products := r.Group("/products")
	{
		products.GET("", productHandler.GetAll)            // Listar todos
		products.GET("/:id", productHandler.GetByID)       // Buscar por ID
		products.POST("", productHandler.Create)           // Criar produto
		products.PUT("/:id", productHandler.Update)        // Atualizar produto
		products.DELETE("/:id", productHandler.Delete)     // Remover produto

		// Endpoint interno: débito de saldo (chamado pelo billing-service)
		products.PATCH("/:id/debit", productHandler.DebitBalance)
	}

	return r
}
