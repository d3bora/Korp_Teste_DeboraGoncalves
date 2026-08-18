// api-response.model.ts
// Tipagem padrão das respostas de erro retornadas pelos microsserviços Go.
export interface ApiError {
  error: string;
}

// Resposta de impressão de NF retornada pelo billing-service
export interface PrintResponse {
  message: string;
  status: string;
}
