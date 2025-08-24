# MBA Challenge 1

Este projeto é um desafio do MBA, desenvolvido em Go, com foco em arquitetura de software, uso de Redis, controle de taxa (rate limiting) e boas práticas de desenvolvimento.

## Estrutura do Projeto

- `cmd/`: Ponto de entrada da aplicação (`main.go`).
- `configs/`: Configurações e logger.
- `internal/`: Código interno da aplicação.
  - `infra/database/`: Implementação de repositórios e integração com Redis.
  - `infra/internal_error/`: Definições de erros internos.
  - `infra/middleware/`: Middleware de rate limiting.
  - `infra/mocks/`: Mocks para testes.
  - `usecase/limiter/`: Lógica de negócio para o limitador de taxa.

## Como rodar o projeto

### Pré-requisitos
- [Go](https://golang.org/doc/install) 1.18+
- [Docker](https://www.docker.com/get-started) e [Docker Compose](https://docs.docker.com/compose/)

### Subindo o ambiente com Docker Compose

```sh
docker-compose up --build
```

### Rodando localmente

1. Instale as dependências:
   ```sh
   go mod download
   ```
2. Execute a aplicação:
   ```sh
   go run cmd/main.go
   ```

## Testes

Para rodar os testes:

```sh
go test ./...
```

## Funcionalidades
- Middleware de rate limiting usando Redis
- Estrutura modularizada seguindo boas práticas
- Testes automatizados

## Autor
- Julio Sandroni

---

Sinta-se à vontade para contribuir ou sugerir melhorias!
