# Mini Edge API Platform

API desenvolvida em Go com foco em otimização de requisições, uso de cache e integração com serviços externos.

O projeto simula o comportamento de uma camada intermediária (gateway/proxy), responsável por reduzir latência e melhorar a eficiência no acesso a dados.

---

## Tecnologias

- Go (Golang)
- Gin
- net/http

---

## Funcionalidades

### Health Check

GET /health

Verifica se a aplicação está em funcionamento.

---

### Data Endpoint

GET /data?type=products

- Retorna dados simulados
- Suporte a query parameters
- Cache em memória com TTL

---

### External API Proxy

GET /external-data

- Consome uma API externa (JSONPlaceholder)
- Atua como proxy
- Implementa cache para reduzir chamadas externas
- Retorna os dados diretamente ao cliente

---

## Cache

O projeto utiliza um cache em memória com expiração (TTL).

Características:

- Armazenamento por chave
- Expiração automática
- Redução de chamadas externas

Exemplos de chave:

external:posts  
data:type=products

---

## Headers HTTP

### X-Cache

Indica a origem da resposta:

X-Cache: MISS → chamada externa  
X-Cache: HIT → resposta do cache  

---

### Cache-Control

Cache-Control: public, max-age=30

Permite que clientes armazenem a resposta por um período definido.

---

## Fluxo do endpoint /external-data

Client  
↓  
API (cache)  
↓  
API externa (em caso de MISS)  
↓  
Cache  
↓  
Client  

---

## Rate Limiting

A aplicação possui middleware para controle de requisições por cliente.

---

## Estrutura do projeto

cmd/server  
internal/api  
internal/cache  
internal/middleware  
internal/routes  

---

## Como rodar

go mod tidy  
go run ./cmd/server  

A aplicação estará disponível em:

http://localhost:8080

---

## Exemplo de uso

PowerShell:

iwr http://localhost:8080/external-data -UseBasicParsing

---

## Próximos passos

- Cache baseado em query parameters
- Timeout e retry em chamadas externas
- Logs estruturados
- Cache distribuído (ex: Redis)

---

## Autor

Fernanda Ishida