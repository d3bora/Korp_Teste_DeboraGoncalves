// app.config.ts
// Configuração global da aplicação Angular 17 (standalone).
// Registra providers globais: roteamento, HTTP, interceptors e animações.
import { ApplicationConfig } from '@angular/core';
import { provideRouter } from '@angular/router';
import { provideHttpClient, withInterceptorsFromDi } from '@angular/common/http';
import { provideAnimations } from '@angular/platform-browser/animations';
import { HTTP_INTERCEPTORS } from '@angular/common/http';

import { routes } from './app.routes';
import { ErrorInterceptor } from './core/interceptors/error.interceptor';
import { LoadingInterceptor } from './core/interceptors/loading.interceptor';

export const appConfig: ApplicationConfig = {
  providers: [
    // Roteamento com as rotas definidas em app.routes.ts
    provideRouter(routes),

    // HTTP Client com suporte a interceptors via DI
    provideHttpClient(withInterceptorsFromDi()),

    // Animações do Angular Material
    provideAnimations(),

    // Interceptor de erros: captura todos os erros HTTP e exibe snackbar
    {
      provide: HTTP_INTERCEPTORS,
      useClass: ErrorInterceptor,
      multi: true,
    },

    // Interceptor de loading: ativa spinner durante requisições HTTP
    {
      provide: HTTP_INTERCEPTORS,
      useClass: LoadingInterceptor,
      multi: true,
    },
  ],
};
