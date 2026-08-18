// error.interceptor.ts
// Interceptor HTTP global que captura todos os erros de resposta.
// Exibe notificação ao usuário com a mensagem do backend e relança o erro.
import { Injectable } from '@angular/core';
import {
  HttpRequest,
  HttpHandler,
  HttpEvent,
  HttpInterceptor,
  HttpErrorResponse,
} from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError } from 'rxjs/operators';
import { MatSnackBar } from '@angular/material/snack-bar';

@Injectable()
export class ErrorInterceptor implements HttpInterceptor {
  constructor(private snackBar: MatSnackBar) {}

  intercept(
    request: HttpRequest<unknown>,
    next: HttpHandler
  ): Observable<HttpEvent<unknown>> {
    return next.handle(request).pipe(
      // catchError: captura erros HTTP e exibe mensagem amigável
      catchError((error: HttpErrorResponse) => {
        const message = this.extractErrorMessage(error);
        this.showErrorNotification(message);
        // Relança o erro para que o componente possa reagir se necessário
        return throwError(() => error);
      })
    );
  }

  /** Extrai a mensagem de erro do corpo da resposta ou usa mensagens padrão */
  private extractErrorMessage(error: HttpErrorResponse): string {
    // Mensagem retornada pelo backend Go: { "error": "..." }
    if (error.error?.error) {
      return error.error.error;
    }

    // Mensagens padrão por status HTTP
    switch (error.status) {
      case 0:
        return 'Sem conexão com o servidor. Verifique se os serviços estão em execução.';
      case 400:
        return 'Requisição inválida. Verifique os dados informados.';
      case 404:
        return 'Registro não encontrado.';
      case 422:
        return error.error?.error || 'Operação não permitida pelas regras de negócio.';
      case 503:
        return 'Serviço de estoque temporariamente indisponível. Tente novamente.';
      default:
        return `Erro inesperado (${error.status}). Tente novamente.`;
    }
  }

  private showErrorNotification(message: string): void {
    this.snackBar.open(message, 'Fechar', {
      duration: 5000,
      horizontalPosition: 'end',
      verticalPosition: 'top',
      panelClass: ['error-snackbar'],
    });
  }
}
