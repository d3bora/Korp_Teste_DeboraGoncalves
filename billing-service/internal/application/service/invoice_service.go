// Pacote service contém os casos de uso do microsserviço de Faturamento.
// PrintInvoice é o caso de uso principal: valida, debita estoque e fecha a NF.
package service

import (
	"fmt"

	"github.com/deboragoncalves/korp/billing-service/internal/domain/entity"
	"github.com/deboragoncalves/korp/billing-service/internal/domain/repository"
	"github.com/deboragoncalves/korp/billing-service/internal/infrastructure/client"
)

// InvoiceService orquestra os casos de uso de Notas Fiscais.
type InvoiceService struct {
	repo        repository.InvoiceRepository
	stockClient *client.StockClient // Comunicação com o microsserviço de Estoque
}

// NewInvoiceService cria o service via Dependency Injection.
func NewInvoiceService(repo repository.InvoiceRepository, stockClient *client.StockClient) *InvoiceService {
	return &InvoiceService{repo: repo, stockClient: stockClient}
}

// GetAll retorna todas as notas fiscais.
func (s *InvoiceService) GetAll() ([]entity.Invoice, error) {
	return s.repo.FindAll()
}

// GetByID retorna uma nota fiscal pelo ID.
func (s *InvoiceService) GetByID(id uint) (*entity.Invoice, error) {
	return s.repo.FindByID(id)
}

// Create cria uma nova NF com numeração sequencial automática.
func (s *InvoiceService) Create(invoice *entity.Invoice) error {
	// Valida os itens antes de prosseguir
	if err := invoice.Validate(); err != nil {
		return err
	}

	// Obtém o próximo número sequencial
	nextNum, err := s.repo.NextNumber()
	if err != nil {
		return fmt.Errorf("erro ao gerar número da nota: %w", err)
	}

	invoice.Number = nextNum
	invoice.Status = entity.StatusOpen // Nova NF sempre começa Aberta

	return s.repo.Create(invoice)
}

// PrintInvoice é o caso de uso de impressão de NF.
// Fluxo:
//  1. Verifica se a NF pode ser impressa (status Aberta)
//  2. Para cada item, debita o saldo no stock-service
//  3. Fecha a NF (muda status para Fechada)
//
// Se o stock-service falhar, retorna erro detalhado para o frontend.
func (s *InvoiceService) PrintInvoice(id uint) (*entity.Invoice, error) {
	invoice, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Regra de negócio: só imprime NF com status Aberta
	if err := invoice.CanPrint(); err != nil {
		return nil, err
	}

	// Debita o saldo de cada produto no stock-service
	// ⚠️ Cenário de falha: se o stock-service cair aqui, retornamos erro claro
	for _, item := range invoice.Items {
		if err := s.stockClient.DebitProductBalance(item.ProductID, item.Quantity); err != nil {
			return nil, fmt.Errorf("falha ao debitar estoque do produto '%s': %w", item.ProductCode, err)
		}
	}

	// Fecha a NF após débito bem-sucedido de todos os itens
	invoice.Close()
	if err := s.repo.Update(invoice); err != nil {
		return nil, fmt.Errorf("erro ao fechar nota fiscal: %w", err)
	}

	return invoice, nil
}
