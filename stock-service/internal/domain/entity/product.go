// Pacote entity define as entidades puras do domínio de Estoque.
// Esta camada NÃO depende de banco de dados, frameworks ou bibliotecas externas.
package entity

import (
	"errors"
	"time"
)

// Product representa um produto cadastrado no sistema.
// É a entidade central do microsserviço de Estoque.
type Product struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Code        string    `json:"code" gorm:"uniqueIndex;not null"`        // Código único do produto
	Description string    `json:"description" gorm:"not null"`             // Nome/descrição do produto
	Balance     int       `json:"balance" gorm:"not null;default:0"`       // Saldo disponível em estoque
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Validate verifica se os campos obrigatórios do produto estão preenchidos.
// Retorna erro descritivo se alguma regra de negócio for violada.
func (p *Product) Validate() error {
	if p.Code == "" {
		return errors.New("código do produto é obrigatório")
	}
	if p.Description == "" {
		return errors.New("descrição do produto é obrigatória")
	}
	if p.Balance < 0 {
		return errors.New("saldo do produto não pode ser negativo")
	}
	return nil
}

// DebitBalance desconta a quantidade utilizada do saldo.
// Retorna erro se o saldo for insuficiente (regra de negócio).
func (p *Product) DebitBalance(quantity int) error {
	if quantity <= 0 {
		return errors.New("quantidade para débito deve ser maior que zero")
	}
	if p.Balance < quantity {
		return errors.New("saldo insuficiente para o produto: " + p.Code)
	}
	p.Balance -= quantity
	return nil
}
