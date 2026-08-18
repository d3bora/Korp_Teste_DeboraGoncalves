// Implementação concreta do InvoiceRepository usando GORM + PostgreSQL.
package repository

import (
	"errors"
	"fmt"

	"github.com/deboragoncalves/korp/billing-service/internal/domain/entity"
	domainRepo "github.com/deboragoncalves/korp/billing-service/internal/domain/repository"
	"gorm.io/gorm"
)

// GormInvoiceRepository é o Adapter concreto para notas fiscais.
type GormInvoiceRepository struct {
	db *gorm.DB
}

// NewGormInvoiceRepository cria o repositório via Dependency Injection.
func NewGormInvoiceRepository(db *gorm.DB) domainRepo.InvoiceRepository {
	return &GormInvoiceRepository{db: db}
}

// FindAll retorna todas as NFs com seus itens (Preload para eager loading).
func (r *GormInvoiceRepository) FindAll() ([]entity.Invoice, error) {
	var invoices []entity.Invoice
	result := r.db.Preload("Items").Order("number DESC").Find(&invoices)
	if result.Error != nil {
		return nil, fmt.Errorf("erro ao buscar notas fiscais: %w", result.Error)
	}
	return invoices, nil
}

// FindByID busca uma NF pelo ID incluindo seus itens.
func (r *GormInvoiceRepository) FindByID(id uint) (*entity.Invoice, error) {
	var invoice entity.Invoice
	result := r.db.Preload("Items").First(&invoice, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("nota fiscal com ID %d não encontrada", id)
		}
		return nil, fmt.Errorf("erro ao buscar nota fiscal: %w", result.Error)
	}
	return &invoice, nil
}

// Create persiste uma nova NF com seus itens em uma transação.
// Garante atomicidade: ou tudo é salvo ou nada é salvo.
func (r *GormInvoiceRepository) Create(invoice *entity.Invoice) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(invoice).Error; err != nil {
			return fmt.Errorf("erro ao criar nota fiscal: %w", err)
		}
		return nil
	})
}

// Update salva as alterações de uma NF (ex: mudança de status para Fechada).
func (r *GormInvoiceRepository) Update(invoice *entity.Invoice) error {
	result := r.db.Save(invoice)
	if result.Error != nil {
		return fmt.Errorf("erro ao atualizar nota fiscal: %w", result.Error)
	}
	return nil
}

// NextNumber retorna o próximo número sequencial de NF.
// Usa MAX(number) + 1 para garantir sequência crescente.
func (r *GormInvoiceRepository) NextNumber() (int, error) {
	var maxNumber int
	result := r.db.Model(&entity.Invoice{}).Select("COALESCE(MAX(number), 0)").Scan(&maxNumber)
	if result.Error != nil {
		return 0, fmt.Errorf("erro ao calcular próximo número: %w", result.Error)
	}
	return maxNumber + 1, nil
}
