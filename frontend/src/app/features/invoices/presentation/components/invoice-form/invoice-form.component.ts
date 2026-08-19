// invoice-form.component.ts
// Modal de criação de Notas Fiscais.
// Permite buscar produtos disponíveis e adicionar múltiplos itens com quantidade.
import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, FormArray, Validators, ReactiveFormsModule } from '@angular/forms';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';

import { MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatSelectModule } from '@angular/material/select';
import { MatDividerModule } from '@angular/material/divider';
import { MatTableModule } from '@angular/material/table';

import { InvoiceService } from '../../../domain/services/invoice.service';
import { ProductService } from '../../../../products/domain/services/product.service';
import { Product } from '../../../../products/domain/models/product.model';
import { CreateInvoiceItemDto } from '../../../domain/models/invoice.model';

@Component({
  selector: 'app-invoice-form',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatDialogModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatIconModule,
    MatSelectModule,
    MatDividerModule,
    MatTableModule,
  ],
  templateUrl: './invoice-form.component.html',
  styleUrls: ['./invoice-form.component.scss'],
})
export class InvoiceFormComponent implements OnInit, OnDestroy {
  form!: FormGroup;
  availableProducts: Product[] = [];
  isSubmitting = false;
  displayedItemColumns = ['product', 'quantity', 'remove'];

  private destroy$ = new Subject<void>();

  constructor(
    private fb: FormBuilder,
    private invoiceService: InvoiceService,
    private productService: ProductService,
    private dialogRef: MatDialogRef<InvoiceFormComponent>
  ) {}

  // ngOnInit: carrega produtos disponíveis e inicializa o formulário com 1 item vazio
  ngOnInit(): void {
    this.form = this.fb.group({
      items: this.fb.array([], Validators.required),
    });

    // Carrega produtos do stock-service para o select
    this.productService
      .getAll()
      .pipe(takeUntil(this.destroy$))
      .subscribe({
        next: (products) => {
          // Filtra produtos sem saldo para não permitir seleção
          this.availableProducts = products.filter((p) => p.balance > 0);
          this.addItem(); // Adiciona primeiro item vazio por padrão
        },
      });
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  /** FormArray de itens da NF */
  get items(): FormArray {
    return this.form.get('items') as FormArray;
  }

  /** Adiciona uma nova linha de produto à NF */
  addItem(): void {
    this.items.push(
      this.fb.group({
        product: [null, Validators.required],
        quantity: [1, [Validators.required, Validators.min(1)]],
      })
    );
  }

  /** Remove uma linha de produto da NF */
  removeItem(index: number): void {
    this.items.removeAt(index);
  }

  /** Valida se a quantidade não excede o saldo disponível */
  getMaxQuantity(index: number): number {
    const product: Product = this.items.at(index).get('product')?.value;
    return product?.balance ?? 999;
  }

  /** Submete o formulário criando a NF no billing-service */
  onSubmit(): void {
    if (this.form.invalid) return;
    this.isSubmitting = true;

    const itemDtos: CreateInvoiceItemDto[] = this.items.value.map(
      (item: { product: Product; quantity: number }) => ({
        product_id: item.product.id,
        product_code: item.product.code,
        product_name: item.product.description,
        quantity: item.quantity,
      })
    );

    this.invoiceService
      .create({ items: itemDtos })
      .pipe(takeUntil(this.destroy$))
      .subscribe({
        next: () => this.dialogRef.close(true),
        error: () => (this.isSubmitting = false),
      });
  }

  onCancel(): void {
    this.dialogRef.close(false);
  }
}
