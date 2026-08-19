// invoice-list.component.ts
// Página de listagem de Notas Fiscais com botão de impressão.
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
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatChipsModule } from '@angular/material/chips';
import { MatBadgeModule } from '@angular/material/badge';

import { InvoiceService } from '../../../domain/services/invoice.service';
import { Invoice } from '../../../domain/models/invoice.model';
import { InvoiceFormComponent } from '../../components/invoice-form/invoice-form.component';
import { ConfirmDialogComponent } from '../../../../../shared/components/confirm-dialog/confirm-dialog.component';

@Component({
  selector: 'app-invoice-list',
  standalone: true,
  imports: [
    CommonModule,
    MatTableModule,
    MatButtonModule,
    MatIconModule,
    MatDialogModule,
    MatSnackBarModule,
    MatTooltipModule,
    MatProgressSpinnerModule,
    MatChipsModule,
    MatBadgeModule,
  ],
  templateUrl: './invoice-list.component.html',
  styleUrls: ['./invoice-list.component.scss'],
})
export class InvoiceListComponent implements OnInit, OnDestroy {
  invoices: Invoice[] = [];
  displayedColumns = ['number', 'status', 'items', 'created_at', 'actions'];

  // Controla o spinner de impressão por NF (key = id da NF)
  printingIds = new Set<number>();

  private destroy$ = new Subject<void>();

  constructor(
    private invoiceService: InvoiceService,
    private dialog: MatDialog,
    private snackBar: MatSnackBar
  ) {}

  // ngOnInit: carrega as NFs ao inicializar o componente
  ngOnInit(): void {
    this.loadInvoices();
  }

  // ngOnDestroy: cancela subscriptions para evitar memory leak
  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  loadInvoices(): void {
    this.invoiceService
      .getAll()
      .pipe(takeUntil(this.destroy$))
      .subscribe({
        next: (invoices) => (this.invoices = invoices),
        error: () => {},
      });
  }

  openCreateDialog(): void {
    const dialogRef = this.dialog.open(InvoiceFormComponent, {
      width: '700px',
      disableClose: true, // Evita fechar acidentalmente ao clicar fora
    });

    dialogRef.afterClosed().subscribe((created) => {
      if (created) {
        this.loadInvoices();
        this.showSuccess('Nota fiscal criada com sucesso!');
      }
    });
  }

  /**
   * Imprime a NF: exibe confirmação → spinner → debita estoque → fecha NF.
   * Usa switchMap internamente no service para encadear operações.
   */
  printInvoice(invoice: Invoice): void {
    if (invoice.status !== 'Aberta') {
      this.snackBar.open('Esta nota fiscal já foi impressa.', 'Ok', { duration: 3000 });
      return;
    }

    const dialogRef = this.dialog.open(ConfirmDialogComponent, {
      data: {
        title: 'Imprimir Nota Fiscal',
        message: `Confirma a impressão da NF nº ${invoice.number}? Esta ação irá debitar o estoque dos produtos e fechar a nota.`,
        confirmLabel: 'Imprimir',
        cancelLabel: 'Cancelar',
      },
    });

    dialogRef.afterClosed().subscribe((confirmed) => {
      if (!confirmed) return;

      // Ativa o spinner individual desta NF
      this.printingIds.add(invoice.id);

      this.invoiceService
        .print(invoice.id)
        .pipe(takeUntil(this.destroy$))
        .subscribe({
          next: () => {
            this.printingIds.delete(invoice.id);
            this.loadInvoices(); // Recarrega para exibir status Fechada
            this.showSuccess(`NF nº ${invoice.number} impressa com sucesso!`);
          },
          error: () => {
            this.printingIds.delete(invoice.id);
            // ErrorInterceptor já exibe o snackbar de erro
          },
        });
    });
  }

  isPrinting(invoiceId: number): boolean {
    return this.printingIds.has(invoiceId);
  }

  private showSuccess(message: string): void {
    this.snackBar.open(message, 'Fechar', {
      duration: 4000,
      horizontalPosition: 'end',
      verticalPosition: 'top',
      panelClass: ['success-snackbar'],
    });
  }
}
