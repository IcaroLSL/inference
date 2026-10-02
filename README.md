# API de Inferência de Hotel

API em Go para criar sessões de inferência e recomendar o tipo de acomodação com base em dados do cliente.

## Visão geral

A API gera uma sessão, coleta dados em etapas e retorna a inferência final com a regra que disparou a recomendação.

### Stack

- Go
- Chi Router
- MongoDB
- CORS habilitado

## Requisitos

Antes de iniciar o serviço, configure as variáveis de ambiente:

```bash
export PORT=8080
export MONGO_URI="mongodb://localhost:27017"
export MONGO_DB="hotel_inference"
export MONGO_TLS_INSECURE="false"
```

Em seguida, rode:

```bash
go run .
```

## Endpoints

### 1) Informações da API

```http
GET /
```

Retorna um resumo dos endpoints disponíveis.

Resposta exemplo:

```json
{
  "message": "Ymir API is running",
  "routes": [
    "Info:              GET  = /",
    "Inference Engine:  GET  = /v1/inference",
    "Inference Session: GET  = /v1/inference/session",
    "Inference Data:    POST = /v1/inference/data"
  ]
}
```

---

### 2) Criar sessão

```http
GET /v1/inference/session
```

Cria uma nova sessão ativa no MongoDB e retorna um identificador único.

Resposta exemplo:

```json
{
  "sessionId": "64f0d4a7d7f01a1c91b6d5fe"
}
```

---

### 3) Enviar dados da inferência

```http
POST /v1/inference/data
Content-Type: application/json
```

Body:

```json
{
  "sessionId": "64f0d4a7d7f01a1c91b6d5fe",
  "step": 1,
  "data": 2
}
```

#### Campos

- `sessionId` (string): ID da sessão criada anteriormente.
- `step` (number): etapa da coleta de dados. Deve estar entre `1` e `5`.
- `data` (any): valor conforme a etapa.

#### Etapas permitidas

| Step | Campo | Tipo | Validação |
| --- | --- | --- | --- |
| 1 | people | number | entre 1 e 6 |
| 2 | budget | number | 1 = baixo, 2 = médio, 3 = alto |
| 3 | focus | string | exemplos: "economia", "conforto", "privacidade", "espaco" |
| 4 | child | boolean | true / false |
| 5 | accommodation | string | exemplos: "lua de mel", "normal" |

Exemplo de envio completo:

```bash
curl -X POST http://localhost:8080/v1/inference/data \
  -H "Content-Type: application/json" \
  -d '{
    "sessionId": "64f0d4a7d7f01a1c91b6d5fe",
    "step": 1,
    "data": 2
  }'
```

Resposta de sucesso:

```json
{
  "message": "data received successfully",
  "facts": {
    "people": 2
  }
}
```

---

### 4) Executar inferência

```http
GET /v1/inference?sessionId={sessionId}
```

Busca a sessão ativa e aplica a regra de inferência com os fatos coletados.

Exemplo:

```bash
curl "http://localhost:8080/v1/inference?sessionId=64f0d4a7d7f01a1c91b6d5fe"
```

Resposta exemplo:

```json
{
  "inferred": "Quarto Luxo",
  "rule": "R03"
}
```

---

## Regras de inferência

A API compara os dados da sessão com um conjunto de regras. Os resultados possíveis são:

- `Quarto Standard`
- `Quarto Luxo`
- `Suíte Premium`
- `Quarto Família`
- `indefinido`

### Exemplos de regras importantes

- 1 pessoa + orçamento baixo => `Quarto Standard`
- 2 pessoas + orçamento médio + foco em conforto => `Quarto Luxo`
- 2 pessoas + orçamento alto + foco em privacidade => `Suíte Premium`
- 3 pessoas + criança => `Quarto Família`
- 2 pessoas + estadia `lua de mel` => `Suíte Premium`

---

## Códigos de resposta

| Código | Descrição |
| --- | --- |
| 200 | Requisição concluída com sucesso |
| 400 | Dados inválidos ou ausentes |
| 404 | Sessão não encontrada |
| 500 | Erro interno do servidor |

## Tratamento de erros

Erros comuns:

```json
{
  "error": "sessionId is required"
}
```

```json
{
  "error": "the number of people must be a number from 1 to 6"
}
```

```json
{
  "error": "budget must be 1 (baixo), 2 (médio) or 3 (alto)"
}
```

```json
{
  "error": "active session not found"
}
```

---

## Fluxo recomendado da API

1. Crie uma sessão com `GET /v1/inference/session`
2. Envie `step` 1 até 5 em requisições `POST /v1/inference/data`
3. Chame a inferência com `GET /v1/inference?sessionId=<id>`
4. Interprete o campo `inferred` e a `rule` que disparou a resposta

## Exemplo completo

```bash
SESSION=$(curl -s http://localhost:8080/v1/inference/session | sed -E 's/.*"sessionId":"([^"]+)".*/\1/')

curl -s -X POST http://localhost:8080/v1/inference/data \
  -H "Content-Type: application/json" \
  -d '{"sessionId":"'"$SESSION"'","step":1,"data":2}'

curl -s -X POST http://localhost:8080/v1/inference/data \
  -H "Content-Type: application/json" \
  -d '{"sessionId":"'"$SESSION"'","step":2,"data":3}'

curl -s -X POST http://localhost:8080/v1/inference/data \
  -H "Content-Type: application/json" \
  -d '{"sessionId":"'"$SESSION"'","step":3,"data":"Privacidade"}'

curl -s -X POST http://localhost:8080/v1/inference/data \
  -H "Content-Type: application/json" \
  -d '{"sessionId":"'"$SESSION"'","step":4,"data":false}'

curl -s -X POST http://localhost:8080/v1/inference/data \
  -H "Content-Type: application/json" \
  -d '{"sessionId":"'"$SESSION"'","step":5,"data":"Lua de mel"}'

curl -s "http://localhost:8080/v1/inference?sessionId=$SESSION"
```

## Observação

A API persiste cada sessão em MongoDB na coleção `inference_sessions` e guarda os fatos em `facts` para posterior cálculo da recomendação.
