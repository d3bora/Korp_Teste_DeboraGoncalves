// Pacote repository define a interface (Port) do repositório de produtos.
// Seguindo Clean Architecture: o domínio define O QUE precisa, a infra implementa O COMO.
package repository

import "github.com/deboragoncalves/korp/stock-service/internal/domain/entity"

// ProductRepository é a interface que abstrai o acesso ao banco de dados.
// A camada de domínio depende apenas desta interface, nunca da implementação concreta.
type ProductRepository interface {
	// FindAll retorna todos os produtos cadastrados
	FindAll() ([]entity.Product, error)

	// FindByID busca um produto pelo ID; retorna erro se não encontrado
	FindByID(id uint) (*entity.Product, error)

	// FindByCode busca um produto pelo código único
	FindByCode(code string) (*entity.Product, error)

	// Create persiste um novo produto no banco
	Create(product *entity.Product) error

	// Update atualiza os dados de um produto existente
	Update(product *entity.Product) error

	// Delete remove um produto pelo ID
	Delete(id uint) error
}
