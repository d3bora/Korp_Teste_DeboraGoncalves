// invoice.service.ts
// Serviço Angular responsável por todas as chamadas HTTP ao billing-service.
import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Invoice, CreateInvoiceDto } from '../models/invoice.model';
import { PrintResponse } from '../../../../core/models/api-response.model';

@Injectable({
  providedIn: 'root',
})
export class InvoiceService {
  // URL base do billing-service via Nginx gateway
  private readonly apiUrl = '/api/billing/invoices';

  constructor(private http: HttpClient) {}

  /**
   * Retorna todas as notas fiscais com seus itens.
   */
  getAll(): Observable<Invoice[]> {
    return this.http.get<Invoice[]>(this.apiUrl);
  }

  /**
   * Busca uma nota fiscal pelo ID.
   */
  getById(id: number): Observable<Invoice> {
    return this.http.get<Invoice>(`${this.apiUrl}/${id}`);
  }

  /**
   * Cria uma nova NF. O backend define o número sequencial e status 'Aberta'.
   */
  create(dto: CreateInvoiceDto): Observable<Invoice> {
    return this.http.post<Invoice>(this.apiUrl, dto);
  }

  /**
   * Imprime a NF: debita o estoque e muda o status para 'Fechada'.
   * O LoadingInterceptor ativa automaticamente o spinner durante esta chamada.
   */
  print(id: number): Observable<PrintResponse> {
    return this.http.post<PrintResponse>(`${this.apiUrl}/${id}/print`, {});
  }
}
