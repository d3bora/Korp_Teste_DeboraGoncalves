// product.service.ts
// Serviço Angular responsável por todas as chamadas HTTP ao stock-service.
// Usa HttpClient que retorna Observables — padrão RxJS no Angular.
import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Product, CreateProductDto, UpdateProductDto } from '../models/product.model';

@Injectable({
  providedIn: 'root',
})
export class ProductService {
  // URL base do stock-service via Nginx gateway
  private readonly apiUrl = '/api/stock/products';

  constructor(private http: HttpClient) {}

  /**
   * Retorna todos os produtos cadastrados.
   * Observable: o componente se inscreve e reage quando os dados chegam.
   */
  getAll(): Observable<Product[]> {
    return this.http.get<Product[]>(this.apiUrl);
  }

  /**
   * Busca um produto pelo ID.
   */
  getById(id: number): Observable<Product> {
    return this.http.get<Product>(`${this.apiUrl}/${id}`);
  }

  /**
   * Cria um novo produto.
   */
  create(dto: CreateProductDto): Observable<Product> {
    return this.http.post<Product>(this.apiUrl, dto);
  }

  /**
   * Atualiza um produto existente.
   */
  update(id: number, dto: UpdateProductDto): Observable<Product> {
    return this.http.put<Product>(`${this.apiUrl}/${id}`, dto);
  }

  /**
   * Remove um produto pelo ID.
   */
  delete(id: number): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/${id}`);
  }
}
