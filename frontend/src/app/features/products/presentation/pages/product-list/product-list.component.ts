// product-list.component.ts
// Página de listagem e gerenciamento de produtos.
// Ciclos de vida usados: ngOnInit (carga inicial) e ngOnDestroy (cancelar subscriptions).
import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';

import { MatTableModule } from '@angular/material/table';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatChipsModule } from '@angular/material/chips';

import { ProductService } from '../../../domain/services/product.service';
import { Product } from '../../../domain/models/product.model';
import { ProductFormComponent } from '../../components/product-form/product-form.component';
import { ConfirmDialogComponent } from '../../../../../shared/components/confirm-dialog/confirm-dialog.component';

@Component({
  selector: 'app-product-list',
  standalone: true,
  imports: [
    CommonModule,
    MatTableModule,
    MatButtonModule,
    MatIconModule,
    MatDialogModule,
    MatSnackBarModule,
    MatTooltipModule,
    MatChipsModule,
  ],
  templateUrl: './product-list.component.html',
  styleUrls: ['./product-list.component.scss'],
})
export class ProductListComponent implements OnInit, OnDestroy {
  products: Product[] = [];
  // Colunas exibidas na MatTable
  displayedColumns = ['code', 'description', 'balance', 'actions'];

  // Subject usado para cancelar subscriptions no ngOnDestroy (evita memory leak)
  private destroy$ = new Subject<void>();

  constructor(
    private productService: ProductService,
    private dialog: MatDialog,
    private snackBar: MatSnackBar
  ) {}

  // ngOnInit: ciclo de vida — executado após a criação do componente
  ngOnInit(): void {
    this.loadProducts();
  }

  // ngOnDestroy: ciclo de vida — cancela todas as subscriptions ao destruir o componente
  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  /** Carrega todos os produtos do stock-service */
  loadProducts(): void {
    this.productService
      .getAll()
      .pipe(takeUntil(this.destroy$)) // Cancela a subscription quando o componente for destruído
      .subscribe({
        next: (products) => (this.products = products),
        error: () => {}, // ErrorInterceptor já exibe o snackbar de erro
      });
  }

  /** Abre modal de criação de produto */
  openCreateDialog(): void {
    const dialogRef = this.dialog.open(ProductFormComponent, {
      width: '500px',
      data: null, // null = modo criação
    });

    dialogRef.afterClosed().subscribe((created) => {
      if (created) {
        this.loadProducts(); // Recarrega a lista após criação
        this.showSuccess('Produto cadastrado com sucesso!');
      }
    });
  }

  /** Abre modal de edição de produto */
  openEditDialog(product: Product): void {
    const dialogRef = this.dialog.open(ProductFormComponent, {
      width: '500px',
      data: product, // passa o produto atual para edição
    });

    dialogRef.afterClosed().subscribe((updated) => {
      if (updated) {
        this.loadProducts();
        this.showSuccess('Produto atualizado com sucesso!');
      }
    });
  }

  /** Confirma e executa a exclusão de um produto */
  deleteProduct(product: Product): void {
    const dialogRef = this.dialog.open(ConfirmDialogComponent, {
      data: {
        title: 'Excluir Produto',
        message: `Deseja excluir o produto "${product.description}"?`,
        confirmLabel: 'Excluir',
      },
    });

    dialogRef.afterClosed().subscribe((confirmed) => {
      if (confirmed) {
        this.productService
          .delete(product.id)
          .pipe(takeUntil(this.destroy$))
          .subscribe({
            next: () => {
              this.loadProducts();
              this.showSuccess('Produto excluído com sucesso!');
            },
          });
      }
    });
  }

  /** Retorna classe CSS baseada no saldo (alerta visual para saldo baixo) */
  getBalanceClass(balance: number): string {
    if (balance === 0) return 'balance-zero';
    if (balance <= 5) return 'balance-low';
    return 'balance-ok';
  }

  private showSuccess(message: string): void {
    this.snackBar.open(message, 'Fechar', {
      duration: 3000,
      horizontalPosition: 'end',
      verticalPosition: 'top',
      panelClass: ['success-snackbar'],
    });
  }
}
