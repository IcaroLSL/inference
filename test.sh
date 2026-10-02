#!/bin/bash

BASE=http://localhost:8080

# 1) Criar sessão e guardar o ID
SESSION=$(curl -s $BASE/v1/inference/session | sed -E 's/.*"sessionId":"([^"]+)".*/\1/')
echo "Sessão: $SESSION"

# 2) Enviar os 5 passos
curl -s -X POST $BASE/v1/inference/data -H "Content-Type: application/json" \
  -d "{\"sessionId\":\"$SESSION\",\"step\":1,\"data\":2}"

curl -s -X POST $BASE/v1/inference/data -H "Content-Type: application/json" \
  -d "{\"sessionId\":\"$SESSION\",\"step\":2,\"data\":3}"

curl -s -X POST $BASE/v1/inference/data -H "Content-Type: application/json" \
  -d "{\"sessionId\":\"$SESSION\",\"step\":3,\"data\":\"Privacidade\"}"

curl -s -X POST $BASE/v1/inference/data -H "Content-Type: application/json" \
  -d "{\"sessionId\":\"$SESSION\",\"step\":4,\"data\":false}"

curl -s -X POST $BASE/v1/inference/data -H "Content-Type: application/json" \
  -d "{\"sessionId\":\"$SESSION\",\"step\":5,\"data\":\"Lua de mel\"}"

# 3) Pedir a inferência
curl -s "$BASE/v1/inference?sessionId=$SESSION" \ -H "Content-Type: application/json"