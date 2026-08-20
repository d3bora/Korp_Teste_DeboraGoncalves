// Pacote database gerencia a conexão com o PostgreSQL do billing-service.
package database

import (
	"fmt"
	"log"
	"os"

	"github.com/deboragoncalves/korp/billing-service/internal/domain/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB é a instância global da conexão com billing_db.
var DB *gorm.DB

// Connect estabelece conexão com o PostgreSQL do serviço de Faturamento.
func Connect() error {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	name := getEnv("DB_NAME", "billing_db")
	user := getEnv("DB_USER", "billing_user")
	pass := getEnv("DB_PASSWORD", "billing_pass")

	dsn := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable TimeZone=UTC",
		host, port, name, user, pass,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("falha ao conectar ao banco de dados: %w", err)
	}

	DB = db
	log.Println("✅ Conexão com PostgreSQL estabelecida (billing_db)")
	return nil
}

// Migrate cria/atualiza as tabelas de Invoice e InvoiceItem.
func Migrate() error {
	// A ordem importa: Invoice deve existir antes de InvoiceItem (FK)
	err := DB.AutoMigrate(&entity.Invoice{}, &entity.InvoiceItem{})
	if err != nil {
		return fmt.Errorf("falha ao executar migrations: %w", err)
	}
	log.Println("✅ Migrations executadas com sucesso")
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
