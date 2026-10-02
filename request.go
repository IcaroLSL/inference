package main

type InferenceDataRequest struct {
    SessionID string `json:"sessionId"`
    Data      any    `json:"data"`
	Step      int    `json:"step"`
}


