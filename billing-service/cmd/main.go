// main.go — Entrypoint do microsserviço de Faturamento.
// Monta todas as dependências e inicia o servidor HTTP.
package main

import (
	"log"
	"os"

	"github.com/deboragoncalves/korp/billing-service/internal/application/service"
	"github.com/deboragoncalves/korp/billing-service/internal/infrastructure/client"
	"github.com/deboragoncalves/korp/billing-service/internal/infrastructure/database"
	infraRepo "github.com/deboragoncalves/korp/billing-service/internal/infrastructure/repository"
	"github.com/deboragoncalves/korp/billing-service/internal/presentation/handler"
	"github.com/deboragoncalves/korp/billing-service/internal/presentation/router"
)

func main() {
	// ─── 1. Conectar ao banco de dados ─────────────────────────────────────
	if err := database.Connect(); err != nil {
		log.Fatalf("❌ Erro ao conectar ao banco: %v", err)
	}

	// ─── 2. Executar migrations ────────────────────────────────────────────
	if err := database.Migrate(); err != nil {
		log.Fatalf("❌ Erro ao executar migrations: %v", err)
	}

	// ─── 3. Inicializar o client HTTP do stock-service ─────────────────────
	stockClient := client.NewStockClient()

	// ─── 4. Montar as dependências (Dependency Injection) ──────────────────
	invoiceRepo := infraRepo.NewGormInvoiceRepository(database.DB)
	invoiceService := service.NewInvoiceService(invoiceRepo, stockClient)
	invoiceHandler := handler.NewInvoiceHandler(invoiceService)

	// ─── 5. Iniciar o servidor ─────────────────────────────────────────────
	r := router.Setup(invoiceHandler)

	port := getEnv("PORT", "8082")
	log.Printf("🚀 billing-service rodando na porta %s", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Erro ao iniciar servidor: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
