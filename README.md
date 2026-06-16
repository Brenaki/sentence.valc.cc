# Sentence — Frase do Dia

Aplicação web que exibe uma **frase do dia** (citação literária) e permite ao
usuário reagir a ela (curtir / descurtir). A cada dia uma nova frase é buscada de
uma API externa, persistida em MySQL e servida ao front-end.

🔗 **Sistema publicado:** https://sentence.valc.cc
🔗 **Repositório:** https://github.com/Brenaki/sentence.valc.cc

---

## Sobre o projeto acadêmico

Este é o **projeto final individual** da disciplina, conforme solicitado pela
professora **Doutoranda Gabrielly de Queiroz Pereira**.

> O projeto final será individual e poderá ser desenvolvido sobre qualquer tema,
> desde que seja uma aplicação web funcional e com princípios éticos.

### Requisitos solicitados e onde foram atendidos

| Requisito | Atendido em |
|-----------|-------------|
| **HTML, CSS e JavaScript** | Front-end React + TypeScript (compila para HTML/CSS/JS) em `frontend/` |
| **Layout bonito e organizado** | UI com Tailwind CSS + componentes shadcn/ui (`frontend/src/components`) |
| **Uso de flexbox ou grid** | Layout responsivo via utilitários flex/grid do Tailwind |
| **Manipulação do DOM** | React renderiza e atualiza o DOM (cards, estados de carregamento) |
| **Eventos e funções em JavaScript** | Handlers de clique nas reações (`frontend/src/api/reactions.ts`, `App.tsx`) |
| **Persistência de dados em MySQL** | Tabela `frases` em MySQL 8.4 (`backend/migrations/001_create_frases.sql`) |
| **Consumo de API** | Back-end consome a API externa [api-ninjas](https://api-ninjas.com) (`backend/internal/infra/provider/ninja`) e o front-end consome a API própria |

### Princípios éticos

- Citações obtidas de fonte pública (api-ninjas) com atribuição de autor e obra.
- Nenhum dado pessoal do usuário é coletado; as reações são anônimas e agregadas.
- Chaves de API ficam apenas no servidor / variáveis de ambiente, fora do versionamento.

### Entrega

A entrega contém: link do sistema publicado, link do repositório no GitHub,
export SQL do banco e instruções de execução (todos abaixo). Na apresentação,
o sistema é demonstrado funcionando, com explicação do código e principais
funcionalidades.

---

## Arquitetura

```
sentence.valc.cc/
├── backend/     API em Go (Clean Architecture, SOLID, TDD) — ver backend/README.md
└── frontend/    SPA em React + Vite + TypeScript + Tailwind — ver frontend/README.md
```

- **Back-end (Go):** expõe a API, busca a frase do dia na api-ninjas quando ainda
  não existe registro do dia, persiste em MySQL e registra reações.
- **Front-end (React):** consome a API do back-end, exibe a frase e envia as reações.
- **Banco (MySQL 8.4):** sobe via Docker Compose; o schema é carregado
  automaticamente a partir de `backend/migrations/`.

### Endpoints da API

- `GET /quote-of-the-day` — retorna a frase do dia (busca/persiste se necessário).
- `POST /quotes/{id}/reactions` — corpo `{"reaction": 0|1}` (0 = descurtir, 1 = curtir).
- `GET /healthz` — verificação de saúde.

---

## Banco de dados / Export SQL

O schema fica em [`backend/migrations/001_create_frases.sql`](backend/migrations/001_create_frases.sql)
e é aplicado automaticamente na primeira subida do container MySQL
(`/docker-entrypoint-initdb.d`).

Tabela `frases`: `id`, `quote`, `author`, `work`, `categories` (JSON),
`like_quantity`, `deslike_quantity`, `created_at`, `updated_at`.

Para exportar o banco manualmente:

```bash
docker exec sentence-mysql mysqldump -uroot -proot sentence > export.sql
```

---

## Instruções básicas de execução

Pré-requisitos: **Docker**, **Go 1.25+** e **Node.js 20+** (com `npm`).

### 1. Banco de dados (MySQL via Docker)

```bash
cd backend
docker compose up -d        # sobe MySQL 8.4 e aplica as migrations
```

### 2. Back-end (API em Go)

```bash
cd backend
# .env deve conter API_KEY_NINJA=<sua_chave_api_ninjas>
set -a; source .env; set +a
go run ./cmd/api            # sobe a API em http://localhost:8080
```

Variáveis de ambiente (com defaults em `internal/config`):
`HTTP_ADDR` (`:8080`), `MYSQL_DSN`, `API_KEY_NINJA` (obrigatória), `ALLOW_ORIGINS` (`*`).

### 3. Front-end (React + Vite)

```bash
cd frontend
npm install
npm run dev                 # sobe o front em http://localhost:5173
```

O front-end lê de `frontend/.env`:
`VITE_API_URL` (URL da API, ex. `http://localhost:8080`) e `VITE_API_KEY`.

### Testes

```bash
cd backend && go test ./...   # casos de uso, handlers, provider e repositório
```

---

## Stack

- **Back-end:** Go 1.25, MySQL 8.4, [api-ninjas](https://api-ninjas.com),
  Clean Architecture, TDD (go-sqlmock / httptest).
- **Front-end:** React 19, TypeScript, Vite, Tailwind CSS 4, shadcn/ui,
  Phosphor Icons.
