package main

import (
	"net/http"
	"encoding/json"
	"log"
	"errors"
	"os"
	"strings"
)

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

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

func validateInfos(request InferenceDataRequest, w http.ResponseWriter) error{
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

func toString(v interface{}) string {
	s, _ := v.(string)
	return s
}

func toBool(v interface{}) bool {
	b, _ := v.(bool)
	return b
}

var accentReplacer = strings.NewReplacer(
	"á", "a", "à", "a", "ã", "a", "â", "a",
	"é", "e", "ê", "e", "í", "i",
	"ó", "o", "õ", "o", "ô", "o",
	"ú", "u", "ç", "c",
)

func normalize(s string) string {
	return accentReplacer.Replace(strings.ToLower(strings.TrimSpace(s)))
}