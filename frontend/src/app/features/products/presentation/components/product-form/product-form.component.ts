// product-form.component.ts
// Modal de formulário para criação e edição de produtos.
// Recebe o produto via MAT_DIALOG_DATA (null = criação, objeto = edição).
import { Component, OnInit, Inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';

import { MatDialogModule, MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';

import { ProductService } from '../../../domain/services/product.service';
import { Product } from '../../../domain/models/product.model';

@Component({
  selector: 'app-product-form',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatDialogModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatIconModule,
  ],
  templateUrl: './product-form.component.html',
  styleUrls: ['./product-form.component.scss'],
})
export class ProductFormComponent implements OnInit {
  form!: FormGroup;
  isEditMode: boolean;
  isSubmitting = false;

  constructor(
    private fb: FormBuilder,
    private productService: ProductService,
    private dialogRef: MatDialogRef<ProductFormComponent>,
    // MAT_DIALOG_DATA: null = criar, Product = editar
    @Inject(MAT_DIALOG_DATA) public data: Product | null
  ) {
    this.isEditMode = data !== null;
  }

  // ngOnInit: inicializa o formulário preenchendo com dados existentes em modo edição
  ngOnInit(): void {
    this.form = this.fb.group({
      code: [
        { value: this.data?.code ?? '', disabled: this.isEditMode }, // Código não editável
        [Validators.required, Validators.maxLength(20)],
      ],
      description: [
        this.data?.description ?? '',
        [Validators.required, Validators.maxLength(100)],
      ],
      balance: [
        this.data?.balance ?? 0,
        [Validators.required, Validators.min(0)],
      ],
    });
  }

  /** Submete o formulário — cria ou atualiza dependendo do modo */
  onSubmit(): void {
    if (this.form.invalid) return;
    this.isSubmitting = true;

    const formValue = this.form.getRawValue(); // getRawValue inclui campos disabled

    const request$ = this.isEditMode
      ? this.productService.update(this.data!.id, formValue)
      : this.productService.create(formValue);

    request$.subscribe({
      next: () => this.dialogRef.close(true), // fecha o dialog e sinaliza sucesso
      error: () => (this.isSubmitting = false), // ErrorInterceptor exibe o erro
    });
  }

  onCancel(): void {
    this.dialogRef.close(false);
  }

  get title(): string {
    return this.isEditMode ? 'Editar Produto' : 'Novo Produto';
  }
}
