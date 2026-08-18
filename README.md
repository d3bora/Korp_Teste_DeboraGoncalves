# Korp_Teste_DeboraGoncalves

Sistema de Emissão de Notas Fiscais — Teste Técnico Korp

## Stack

| Camada | Tecnologia |
|--------|-----------|
| Frontend | Angular 17 + Angular Material |
| Backend | Go 1.21 + Gin Framework |
| ORM | GORM |
| Banco | PostgreSQL 16 |
| Gateway | Nginx |
| Containers | Docker + Docker Compose |

## Arquitetura

```
Angular (Frontend :4200)
        │
        ▼
   Nginx :80 (API Gateway)
   ├── /api/stock/*   → stock-service :8081
   └── /api/billing/* → billing-service :8082
                              │
                    (chama stock-service
                     para debitar saldo)
        │                     │
        ▼                     ▼
  PostgreSQL :5432      PostgreSQL :5433
   (stock_db)            (billing_db)
```

## Microsserviços

### stock-service (Estoque)
- Gerencia cadastro de produtos e saldos em estoque
- Endpoint de débito de saldo (chamado pelo billing-service na impressão de NF)

### billing-service (Faturamento)
- Gerencia notas fiscais com numeração sequencial
- Orquestra a impressão: valida status, debita estoque, fecha a NF

## Como executar

### Pré-requisitos
- Docker Desktop instalado e em execução
- Node.js 18+ (para o frontend Angular)
- Go 1.21+ (para desenvolvimento local dos serviços)

### 1. Subir a infraestrutura (banco + serviços + gateway)

```bash
docker-compose up --build -d
```

### 2. Verificar se os serviços estão saudáveis

```bash
docker-compose ps
# Todos devem estar com status "Up"

curl http://localhost/health
# Resposta: {"status":"ok","gateway":"nginx"}
```

### 3. Rodar o frontend Angular

```bash
cd frontend
npm install
npm run start
# Acesse http://localhost:4200
```

### 4. Derrubar tudo

```bash
docker-compose down
# Para remover os volumes (dados do banco):
docker-compose down -v
```

## Estrutura de Pastas

```
Korp_Teste_DeboraGoncalves/
├── stock-service/           # Microsserviço de Estoque (Go/Gin)
│   └── internal/
│       ├── domain/          # Entidades + interfaces (regras puras)
│       ├── application/     # Casos de uso / services
│       ├── infrastructure/  # GORM + conexão DB
│       └── presentation/    # Handlers HTTP + DTOs + rotas
│
├── billing-service/         # Microsserviço de Faturamento (Go/Gin)
│   └── internal/
│       ├── domain/
│       ├── application/
│       ├── infrastructure/  # Inclui client HTTP p/ stock-service
│       └── presentation/
│
├── frontend/                # Angular 17
│   └── src/app/
│       ├── core/            # Interceptors, services globais
│       ├── shared/          # Componentes reutilizáveis
│       └── features/
│           ├── products/    # Feature: Cadastro de Produtos
│           └── invoices/    # Feature: Notas Fiscais + Impressão
│
├── docker-compose.yml       # Orquestração dos containers
└── nginx.conf               # API Gateway
```

## Detalhamento Técnico

### Angular — Ciclos de vida utilizados
- **`ngOnInit`** — carregamento de dados via HTTP ao iniciar cada componente de listagem
- **`ngOnDestroy`** — cancelamento de subscriptions para evitar memory leaks

### RxJS
- **`Observable`** — todos os serviços HTTP retornam Observables via `HttpClient`
- **`switchMap`** — encadeamento de chamadas (ex: criar NF → carregar produtos disponíveis)
- **`catchError`** — captura erros HTTP e repassa mensagem amigável ao usuário
- **`finalize`** — desativa o loading spinner ao terminar qualquer operação assíncrona
- **`BehaviorSubject`** — gerencia estado global de loading no `LoadingService`

### Angular Material (componentes visuais)
- `MatTable` — listagem de produtos e notas fiscais
- `MatDialog` — formulários modais de cadastro/edição
- `MatSnackBar` — notificações de sucesso e erro
- `MatProgressSpinner` — indicador de carregamento durante impressão

### Go / Gin Framework
- Roteamento RESTful com grupos de rotas
- Middleware CORS para comunicação com Angular
- Middleware de recovery para evitar crash em panics

### Tratamento de Erros no Backend (Go)
- Retorno de HTTP status semânticos: `400` (validação), `404` (não encontrado), `422` (regra de negócio), `500` (erro interno)
- Resposta padronizada: `{ "error": "descrição do erro" }`
- Tratamento de falha do stock-service no billing-service com mensagem clara ao usuário

### GORM
- `AutoMigrate` para criação automática de tabelas
- Transações para operações críticas (impressão de NF: debitar saldo + fechar nota atomicamente)

### Gerenciamento de Dependências (Go)
- `go.mod` e `go.sum` gerenciam todas as dependências
- Principais: `github.com/gin-gonic/gin`, `gorm.io/gorm`, `gorm.io/driver/postgres`
