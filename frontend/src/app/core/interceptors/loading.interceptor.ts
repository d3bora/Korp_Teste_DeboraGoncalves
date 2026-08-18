// loading.interceptor.ts
// Interceptor que ativa o LoadingService no início de cada requisição
// e o desativa quando a resposta chega (usando finalize do RxJS).
import { Injectable } from '@angular/core';
import {
  HttpRequest,
  HttpHandler,
  HttpEvent,
  HttpInterceptor,
} from '@angular/common/http';
import { Observable } from 'rxjs';
import { finalize } from 'rxjs/operators';
import { LoadingService } from '../services/loading.service';

@Injectable()
export class LoadingInterceptor implements HttpInterceptor {
  // Contador de requisições ativas (para múltiplas requisições simultâneas)
  private activeRequests = 0;

  constructor(private loadingService: LoadingService) {}

  intercept(
    request: HttpRequest<unknown>,
    next: HttpHandler
  ): Observable<HttpEvent<unknown>> {
    // Ativa loading na primeira requisição
    if (this.activeRequests === 0) {
      this.loadingService.show();
    }
    this.activeRequests++;

    return next.handle(request).pipe(
      // finalize: executado SEMPRE ao terminar (sucesso ou erro)
      finalize(() => {
        this.activeRequests--;
        // Desativa loading apenas quando todas as requisições terminarem
        if (this.activeRequests === 0) {
          this.loadingService.hide();
        }
      })
    );
  }
}
