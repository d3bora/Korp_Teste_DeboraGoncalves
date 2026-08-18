// Pacote repository contém a implementação concreta (Adapter) do repositório de produtos.
// Implementa a interface domain/repository.ProductRepository usando GORM + PostgreSQL.
package repository

import (
	"errors"
	"fmt"

	"github.com/deboragoncalves/korp/stock-service/internal/domain/entity"
	domainRepo "github.com/deboragoncalves/korp/stock-service/internal/domain/repository"
	"gorm.io/gorm"
)

// GormProductRepository é o Adapter concreto que usa GORM para persistência.
type GormProductRepository struct {
	db *gorm.DB
}

// NewGormProductRepository cria uma nova instância via Dependency Injection.
func NewGormProductRepository(db *gorm.DB) domainRepo.ProductRepository {
	return &GormProductRepository{db: db}
}

// FindAll retorna todos os produtos ordenados por código.
func (r *GormProductRepository) FindAll() ([]entity.Product, error) {
	var products []entity.Product
	result := r.db.Order("code ASC").Find(&products)
	if result.Error != nil {
		return nil, fmt.Errorf("erro ao buscar produtos: %w", result.Error)
	}
	return products, nil
}

// FindByID busca um produto pelo ID primário.
// Retorna erro descritivo se o produto não existir.
func (r *GormProductRepository) FindByID(id uint) (*entity.Product, error) {
	var product entity.Product
	result := r.db.First(&product, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("produto com ID %d não encontrado", id)
		}
		return nil, fmt.Errorf("erro ao buscar produto: %w", result.Error)
	}
	return &product, nil
}

// FindByCode busca um produto pelo código único.
func (r *GormProductRepository) FindByCode(code string) (*entity.Product, error) {
	var product entity.Product
	result := r.db.Where("code = ?", code).First(&product)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("produto com código '%s' não encontrado", code)
		}
		return nil, fmt.Errorf("erro ao buscar produto por código: %w", result.Error)
	}
	return &product, nil
}

// Create persiste um novo produto; valida unicidade do código.
func (r *GormProductRepository) Create(product *entity.Product) error {
	// Verifica se o código já existe antes de inserir
	existing, _ := r.FindByCode(product.Code)
	if existing != nil {
		return fmt.Errorf("já existe um produto com o código '%s'", product.Code)
	}

	result := r.db.Create(product)
	if result.Error != nil {
		return fmt.Errorf("erro ao criar produto: %w", result.Error)
	}
	return nil
}

// Update salva as alterações de um produto existente.
func (r *GormProductRepository) Update(product *entity.Product) error {
	result := r.db.Save(product)
	if result.Error != nil {
		return fmt.Errorf("erro ao atualizar produto: %w", result.Error)
	}
	return nil
}

// Delete remove um produto pelo ID.
func (r *GormProductRepository) Delete(id uint) error {
	result := r.db.Delete(&entity.Product{}, id)
	if result.Error != nil {
		return fmt.Errorf("erro ao deletar produto: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("produto com ID %d não encontrado", id)
	}
	return nil
}
