// Pacote entity define as entidades puras do domínio de Faturamento.
// Invoice = Nota Fiscal; InvoiceItem = produto dentro da NF.
package entity

import (
	"errors"
	"time"
)

// InvoiceStatus define os possíveis estados de uma Nota Fiscal.
type InvoiceStatus string

const (
	StatusOpen   InvoiceStatus = "Aberta"  // NF criada, ainda não impressa
	StatusClosed InvoiceStatus = "Fechada" // NF impressa e encerrada
)

// Invoice representa uma Nota Fiscal no sistema.
type Invoice struct {
	ID        uint          `json:"id" gorm:"primaryKey;autoIncrement"`
	Number    int           `json:"number" gorm:"uniqueIndex;not null"`  // Numeração sequencial automática
	Status    InvoiceStatus `json:"status" gorm:"not null;default:'Aberta'"`   // Aberta ou Fechada
	Items     []InvoiceItem `json:"items" gorm:"foreignKey:InvoiceID"`   // Produtos incluídos na NF
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// InvoiceItem representa um produto e sua quantidade dentro de uma NF.
type InvoiceItem struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	InvoiceID   uint      `json:"invoice_id" gorm:"not null;index"`
	ProductID   uint      `json:"product_id" gorm:"not null"`
	ProductCode string    `json:"product_code" gorm:"not null"`
	ProductName string    `json:"product_name" gorm:"not null"`
	Quantity    int       `json:"quantity" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
}

// CanPrint verifica se a NF pode ser impressa (apenas status Aberta).
func (inv *Invoice) CanPrint() error {
	if inv.Status != StatusOpen {
		return errors.New("apenas notas fiscais com status 'Aberta' podem ser impressas")
	}
	if len(inv.Items) == 0 {
		return errors.New("não é possível imprimir uma nota fiscal sem produtos")
	}
	return nil
}

// Close fecha a NF após a impressão bem-sucedida.
func (inv *Invoice) Close() {
	inv.Status = StatusClosed
}

// Validate verifica os campos obrigatórios da NF.
func (inv *Invoice) Validate() error {
	if len(inv.Items) == 0 {
		return errors.New("a nota fiscal deve conter ao menos um produto")
	}
	for _, item := range inv.Items {
		if item.Quantity <= 0 {
			return errors.New("a quantidade de cada produto deve ser maior que zero")
		}
	}
	return nil
}
