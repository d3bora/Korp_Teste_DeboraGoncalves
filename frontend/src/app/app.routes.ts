// app.routes.ts
// Roteamento principal da aplicação com lazy loading das features.
// Redireciona a raiz para /products por padrão.
import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: '',
    redirectTo: 'products',
    pathMatch: 'full',
  },
  {
    // Feature de Produtos — carregada sob demanda (lazy loading)
    path: 'products',
    loadChildren: () =>
      import('./features/products/products.routes').then((m) => m.PRODUCTS_ROUTES),
  },
  {
    // Feature de Notas Fiscais — carregada sob demanda (lazy loading)
    path: 'invoices',
    loadChildren: () =>
      import('./features/invoices/invoices.routes').then((m) => m.INVOICES_ROUTES),
  },
];
