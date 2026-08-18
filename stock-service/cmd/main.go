// main.go é o entrypoint do microsserviço de Estoque.
// Inicializa o banco, executa migrations, monta as dependências e sobe o servidor HTTP.
package main

import (
	"log"
	"os"

	"github.com/deboragoncalves/korp/stock-service/internal/application/service"
	"github.com/deboragoncalves/korp/stock-service/internal/infrastructure/database"
	infraRepo "github.com/deboragoncalves/korp/stock-service/internal/infrastructure/repository"
	"github.com/deboragoncalves/korp/stock-service/internal/presentation/handler"
	"github.com/deboragoncalves/korp/stock-service/internal/presentation/router"
)

func main() {
	// ─── 1. Conectar ao banco de dados ─────────────────────────────────────
	if err := database.Connect(); err != nil {
		log.Fatalf("❌ Erro ao conectar ao banco: %v", err)
	}

	// ─── 2. Executar migrations automáticas ────────────────────────────────
	if err := database.Migrate(); err != nil {
		log.Fatalf("❌ Erro ao executar migrations: %v", err)
	}

	// ─── 3. Montar as dependências (Dependency Injection manual) ───────────
	// Repository (Adapter concreto) → Service → Handler
	productRepo := infraRepo.NewGormProductRepository(database.DB)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	// ─── 4. Configurar e iniciar o servidor HTTP ───────────────────────────
	r := router.Setup(productHandler)

	port := getEnv("PORT", "8081")
	log.Printf("🚀 stock-service rodando na porta %s", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Erro ao iniciar servidor: %v", err)
	}
}

// getEnv retorna o valor da variável de ambiente ou um padrão.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
