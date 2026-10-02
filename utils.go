package main

import (
	"net/http"
	"encoding/json"
	"log"
	"errors"
	"os"
	"strings"
)

// coleta variáveis de ambiente e adiciona um valor padrão caso não exista
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

// escreve uma resposta JSON para o cliente
func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("json encode error: %v", err)
	}
}

// valida os dados recebidos na requisição
func validateInfos(request InferenceDataRequest, w http.ResponseWriter) error {
	if request.SessionID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": "sessionId is required",
		})
		var err = errors.New("sessionId is required")
		return err
	}
	if request.Step <= 0 || request.Step > 5 {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": "the step must be a number from 1 to 5",
		})
		var err = errors.New("the step must be a number from 1 to 5")
		return err
	}
	if request.Data == nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": "data is required",
		})
		var err = errors.New("data is required")
		return err
	}

	return nil
}

// converte um valor para int, retor0nando 0 como padrão
func toInt(v interface{}) int {
	switch n := v.(type) {
	case int:
		return n
	case int32:
		return int(n)
	case int64:
		return int(n)
	case float64:
		return int(n)
	}
	return 0
}

// converte um valor para string, retornando "" como padrão
func toString(v interface{}) string {
	s, _ := v.(string)
	return s
}

// converte um valor para bool, retornando false como padrão
func toBool(v interface{}) bool {
	b, _ := v.(bool)
	return b
}

// remove acentos e espaços e normaliza uma string para minúsculas
var accentReplacer = strings.NewReplacer(
	"á", "a", "à", "a", "ã", "a", "â", "a",
	"é", "e", "ê", "e", "í", "i",
	"ó", "o", "õ", "o", "ô", "o",
	"ú", "u", "ç", "c",
)
func normalize(s string) string {
	return accentReplacer.Replace(strings.ToLower(strings.TrimSpace(s)))
}