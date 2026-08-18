// invoice.model.ts
// Interfaces TypeScript para as entidades do billing-service.

export type InvoiceStatus = 'Aberta' | 'Fechada';

// Item de uma Nota Fiscal (produto + quantidade)
export interface InvoiceItem {
  id: number;
  invoice_id: number;
  product_id: number;
  product_code: string;
  product_name: string;
  quantity: number;
  created_at: string;
}

// Nota Fiscal completa com seus itens
export interface Invoice {
  id: number;
  number: number;       // Numeração sequencial
  status: InvoiceStatus; // 'Aberta' ou 'Fechada'
  items: InvoiceItem[];
  created_at: string;
  updated_at: string;
}

// DTO para criar uma NF (corresponde ao CreateInvoiceRequest do Go)
export interface CreateInvoiceDto {
  items: CreateInvoiceItemDto[];
}

export interface CreateInvoiceItemDto {
  product_id: number;
  product_code: string;
  product_name: string;
  quantity: number;
}
