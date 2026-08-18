// Pacote database gerencia a conexão com o PostgreSQL via GORM.
// Utiliza variáveis de ambiente para não expor credenciais no código.
package database

import (
	"fmt"
	"log"
	"os"

	"github.com/deboragoncalves/korp/stock-service/internal/domain/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB é a instância global da conexão com o banco de dados.
var DB *gorm.DB

// Connect lê as variáveis de ambiente e estabelece a conexão com o PostgreSQL.
// Retorna erro se não conseguir conectar após as tentativas do Docker healthcheck.
func Connect() error {
	// Lê configurações do ambiente (definidas no docker-compose.yml)
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	name := getEnv("DB_NAME", "stock_db")
	user := getEnv("DB_USER", "stock_user")
	pass := getEnv("DB_PASSWORD", "stock_pass")

	dsn := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable TimeZone=America/Sao_Paulo",
		host, port, name, user, pass,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Mostra SQL no log apenas em desenvolvimento
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("falha ao conectar ao banco de dados: %w", err)
	}

	DB = db
	log.Println("✅ Conexão com PostgreSQL estabelecida (stock_db)")
	return nil
}

// Migrate executa as migrations automáticas das entidades do domínio.
// GORM cria/atualiza as tabelas conforme as structs definidas.
func Migrate() error {
	err := DB.AutoMigrate(&entity.Product{})
	if err != nil {
		return fmt.Errorf("falha ao executar migrations: %w", err)
	}
	log.Println("✅ Migrations executadas com sucesso")
	return nil
}

// getEnv retorna o valor da variável de ambiente ou o valor padrão.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
