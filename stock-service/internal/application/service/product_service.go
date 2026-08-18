// Pacote service contém os casos de uso do microsserviço de Estoque.
// Orquestra as regras de negócio usando as interfaces do domínio (sem dependência de infra).
package service

import (
	"github.com/deboragoncalves/korp/stock-service/internal/domain/entity"
	"github.com/deboragoncalves/korp/stock-service/internal/domain/repository"
)

// ProductService contém a lógica de negócio para gerenciamento de produtos.
type ProductService struct {
	repo repository.ProductRepository // Depende da interface, não da implementação concreta
}

// NewProductService cria o service via Dependency Injection.
func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

// GetAll retorna todos os produtos cadastrados.
func (s *ProductService) GetAll() ([]entity.Product, error) {
	return s.repo.FindAll()
}

// GetByID busca um produto pelo ID.
func (s *ProductService) GetByID(id uint) (*entity.Product, error) {
	return s.repo.FindByID(id)
}

// Create valida e persiste um novo produto.
func (s *ProductService) Create(product *entity.Product) error {
	// Delega validação à entidade (regra de negócio pura)
	if err := product.Validate(); err != nil {
		return err
	}
	return s.repo.Create(product)
}

// Update valida e atualiza um produto existente.
func (s *ProductService) Update(id uint, updated *entity.Product) (*entity.Product, error) {
	// Verifica se o produto existe antes de atualizar
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Atualiza apenas os campos editáveis
	existing.Description = updated.Description
	existing.Balance = updated.Balance

	// Valida o código apenas se foi alterado
	if updated.Code != "" && updated.Code != existing.Code {
		existing.Code = updated.Code
	}

	if err := existing.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// Delete remove um produto pelo ID.
func (s *ProductService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// DebitBalance debita a quantidade do saldo do produto.
// Chamado pelo billing-service durante a impressão de uma nota fiscal.
func (s *ProductService) DebitBalance(id uint, quantity int) (*entity.Product, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Delega a regra de saldo insuficiente à entidade
	if err := product.DebitBalance(quantity); err != nil {
		return nil, err
	}

	if err := s.repo.Update(product); err != nil {
		return nil, err
	}
	return product, nil
}
