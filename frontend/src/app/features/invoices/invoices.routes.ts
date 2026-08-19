// invoices.routes.ts
// Rotas lazy-loaded da feature de Notas Fiscais.
import { Routes } from '@angular/router';

export const INVOICES_ROUTES: Routes = [
  {
    path: '',
    loadComponent: () =>
      import('./presentation/pages/invoice-list/invoice-list.component').then(
        (m) => m.InvoiceListComponent
      ),
  },
];
