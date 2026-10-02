package main

type InferSessionResponse struct {
    SessionID string `json:"sessionId"`
}

type InferenceResponse struct {
	Inferred string `json:"inferred"`
	Rule     string `json:"rule"`
}