// product.model.ts
// Interface TypeScript que representa a entidade Product do stock-service.
export interface Product {
  id: number;
  code: string;        // Código único do produto
  description: string; // Nome/descrição do produto
  balance: number;     // Saldo disponível em estoque
  created_at: string;
  updated_at: string;
}

// DTO para criação de produto (corresponde ao CreateProductRequest do Go)
export interface CreateProductDto {
  code: string;
  description: string;
  balance: number;
}

// DTO para atualização de produto
export interface UpdateProductDto {
  code?: string;
  description?: string;
  balance?: number;
}
