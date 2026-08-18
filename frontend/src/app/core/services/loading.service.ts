// loading.service.ts
// Serviço global de loading usando BehaviorSubject (RxJS).
// Permite que qualquer componente ou interceptor ative/desative o spinner global.
import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';

@Injectable({
  providedIn: 'root', // Singleton disponível em toda a aplicação
})
export class LoadingService {
  // BehaviorSubject mantém o estado atual e emite para novos subscribers
  private loadingSubject = new BehaviorSubject<boolean>(false);

  // Observable público — componentes se inscrevem aqui para reagir ao loading
  readonly isLoading$: Observable<boolean> = this.loadingSubject.asObservable();

  /** Ativa o indicador de carregamento */
  show(): void {
    this.loadingSubject.next(true);
  }

  /** Desativa o indicador de carregamento */
  hide(): void {
    this.loadingSubject.next(false);
  }
}
