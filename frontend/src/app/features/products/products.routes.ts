// products.routes.ts
// Rotas lazy-loaded da feature de Produtos.
// Carregado apenas quando o usuário acessa /products (melhora o tempo de boot da app).
import { Routes } from '@angular/router';

export const PRODUCTS_ROUTES: Routes = [
  {
    path: '',
    loadComponent: () =>
      import('./presentation/pages/product-list/product-list.component').then(
        (m) => m.ProductListComponent
      ),
  },
];
