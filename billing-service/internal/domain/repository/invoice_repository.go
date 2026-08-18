// Interface do repositório de Notas Fiscais (Port da Clean Architecture).
package repository

import "github.com/deboragoncalves/korp/billing-service/internal/domain/entity"

// InvoiceRepository abstrai o acesso ao banco de dados para Notas Fiscais.
type InvoiceRepository interface {
	// FindAll retorna todas as notas fiscais com seus itens
	FindAll() ([]entity.Invoice, error)

	// FindByID busca uma NF pelo ID, incluindo seus itens
	FindByID(id uint) (*entity.Invoice, error)

	// Create persiste uma nova NF com seus itens
	Create(invoice *entity.Invoice) error

	// Update salva o estado atual da NF (ex: mudar status para Fechada)
	Update(invoice *entity.Invoice) error

	// NextNumber retorna o próximo número sequencial disponível para NF
	NextNumber() (int, error)
}
