#!/usr/bin/env bash

BASE=${BASE:-http://localhost:8080}

PASS=0
FAIL=0

# Cria uma sessão e imprime o ID
new_session() {
  curl -s "$BASE/v1/inference/session" |
    sed -E 's/.*"sessionId":"([^"]+)".*/\1/'
}

# check <descrição> <método> <status_esperado> <trecho_esperado> <endpoint> <json>
check() {
  local desc=$1
  local method=$2
  local exp_status=$3
  local exp_body=$4
  local endpoint=$5
  local json=$6

  local resp status body

  if [ "$method" = "GET" ]; then
    resp=$(curl -s -w "\n%{http_code}" \
      -X GET "$BASE$endpoint" \
      -H "Content-Type: application/json")
  else
    resp=$(curl -s -w "\n%{http_code}" \
      -X "$method" "$BASE$endpoint" \
      -H "Content-Type: application/json" \
      -d "$json")
  fi

  status=$(echo "$resp" | tail -n1)
  body=$(echo "$resp" | sed '$d')

  if [ "$status" = "$exp_status" ] && echo "$body" | grep -q "$exp_body"; then
    echo "OK   $desc -> $status"
    PASS=$((PASS+1))
  else
    echo "FAIL $desc | esperado: $exp_status contendo '$exp_body' | recebido: $status $body"
    FAIL=$((FAIL+1))
  fi
}

SESSION=$(new_session)

if [ -z "$SESSION" ]; then
  echo "Não consegui criar sessão em $BASE (servidor fora do ar?)"
  exit 1
fi

echo "Sessão de teste: $SESSION"
echo

# ---------- /v1/inference ----------

check "Sessão inexistente" \
  GET \
  404 \
  "active session not found" \
  "/v1/inference?sessionId=000000000000000000000000" \
  ""

check "sessionId malformado" \
  GET \
  400 \
  "invalid sessionId" \
  "/v1/inference?sessionId=abc" \
  ""

check "sessionId ausente" \
  GET \
  400 \
  "sessionId is required" \
  "/v1/inference" \
  ""

# ---------- /v1/inference/data ----------

check "Orçamento fora da faixa (9)" \
  POST \
  400 \
  "budget must be" \
  /v1/inference/data \
  "{\"sessionId\":\"$SESSION\",\"step\":2,\"data\":9}"

check "Pessoas fora da faixa (10)" \
  POST \
  400 \
  "number of people" \
  /v1/inference/data \
  "{\"sessionId\":\"$SESSION\",\"step\":1,\"data\":10}"

check "Pessoas fora da faixa (0)" \
  POST \
  400 \
  "number of people" \
  /v1/inference/data \
  "{\"sessionId\":\"$SESSION\",\"step\":1,\"data\":0}"

check "Pessoas com tipo errado" \
  POST \
  400 \
  "people must be a number" \
  /v1/inference/data \
  "{\"sessionId\":\"$SESSION\",\"step\":1,\"data\":\"dois\"}"

check "Criança com tipo errado" \
  POST \
  400 \
  "child must be a boolean" \
  /v1/inference/data \
  "{\"sessionId\":\"$SESSION\",\"step\":4,\"data\":\"sim\"}"

check "Step inválido (6)" \
  POST \
  400 \
  "step must be" \
  /v1/inference/data \
  "{\"sessionId\":\"$SESSION\",\"step\":6,\"data\":1}"

check "Data ausente" \
  POST \
  400 \
  "data is required" \
  /v1/inference/data \
  "{\"sessionId\":\"$SESSION\",\"step\":1}"

check "Sessão inexistente (data)" \
  POST \
  404 \
  "session not found" \
  /v1/inference/data \
  '{"sessionId":"000000000000000000000000","step":1,"data":2}'

# ---------- Sem dados -> indefinido ----------

EMPTY=$(new_session)

check "Sessão sem dados" \
  GET \
  200 \
  '"inferred":"indefinido"' \
  "/v1/inference?sessionId=$EMPTY" \
  ""

echo

echo "Passou: $PASS | Falhou: $FAIL"

[ "$FAIL" -eq 0 ]
