#!/bin/bash

BASE=${BASE:-http://localhost:8080}
PASS=0; FAIL=0

post() {
  curl -s -o /dev/null -X POST "$BASE/v1/inference/data" \
    -H "Content-Type: application/json" -d "$1"
}

run() {
  local rule=$1 expected=$2 people=$3 budget=$4 focus=$5 child=$6 stay=$7

  local sid
  sid=$(curl -s "$BASE/v1/inference/session" | sed -E 's/.*"sessionId":"([^"]+)".*/\1/')

  post "{\"sessionId\":\"$sid\",\"step\":1,\"data\":$people}"
  post "{\"sessionId\":\"$sid\",\"step\":2,\"data\":$budget}"
  post "{\"sessionId\":\"$sid\",\"step\":3,\"data\":\"$focus\"}"
  post "{\"sessionId\":\"$sid\",\"step\":4,\"data\":$child}"
  post "{\"sessionId\":\"$sid\",\"step\":5,\"data\":\"$stay\"}"

  local out
  out=$(curl -s "$BASE/v1/inference?sessionId=$sid" -H "Content-Type: application/json")

  if echo "$out" | grep -q "\"inferred\":\"$expected\"" && echo "$out" | grep -q "\"rule\":\"$rule\""; then
    echo "OK   $rule -> $expected"
    PASS=$((PASS+1))
  else
    echo "FAIL $rule | esperado: $expected ($rule) | recebido: $out"
    FAIL=$((FAIL+1))
  fi
}

#  regra  resultado     pessoas orç focus       criança estadia
run R01 "Quarto Standard"     1 1 "Economia"     false "Normal"
run R02 "Quarto Standard"     2 1 "Economia"     false "Normal"
run R03 "Quarto Luxo"         2 2 "Conforto"     false "Normal"
run R04 "Suíte Premium"       2 3 "Privacidade"  false "Normal"
run R05 "Quarto Família"      3 2 "Conforto"     true  "Normal"
run R06 "Quarto Família"      4 2 "Espaço"       false "Normal"
run R07 "Suíte Premium"       3 3 "Conforto"     false "Normal"
run R08 "Suíte Premium"       2 2 "Economia"     false "Lua de mel"
run R09 "Quarto Família"      3 2 "Economia"     false "Normal"
run R10 "Quarto Luxo"         1 3 "Conforto"     false "Normal"

echo
echo "Passou: $PASS | Falhou: $FAIL"
