package main

import (
	"net/http"
	"log"
	"encoding/json"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Server) info(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{ "message": "Ymir is running",
												       	"routes": []string{
																			  "Info:              GET  = /",
																		      "Inference Engine:  GET  = /v1/inference",
																			  "Inference Session: GET  = /v1/inference/session",
																			  "Inference Data:    POST = /v1/inference/data",
												  					      },
												 })
}

func (s *Server) InferenceSession(w http.ResponseWriter, r *http.Request) {
	var response InferSessionResponse
	collection := s.db.Collection("inference_sessions")
	result, err := collection.InsertOne(
		r.Context(),
		bson.M{"active": true, "facts": bson.M{}},
	)

	if err != nil {
		log.Printf("mongo insert error: %v", err)

		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "failed to create session",
		})

		return
	}

	sessionID := result.InsertedID.(primitive.ObjectID)
	response.SessionID = sessionID.Hex()

	writeJSON(w, http.StatusOK, response)
}


func (s *Server) inferenceData(w http.ResponseWriter, r *http.Request) {
	var request InferenceDataRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": "invalid request",
		})
		return
	}

	if err := validateInfos(request, w); err != nil {
		return
	}

	sessionID, err := primitive.ObjectIDFromHex(request.SessionID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": "invalid sessionId",
		})
		return
	}

	collection := s.db.Collection("inference_sessions")

	var session InferenceSession

	err = collection.FindOne(
		r.Context(),
		bson.M{"_id": sessionID},
	).Decode(&session)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			writeJSON(w, http.StatusNotFound, map[string]any{
				"error": "session not found",
			})
			return
		}

		log.Printf("mongo find error: %v", err)

		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "failed to find session",
		})
		return
	}

	if session.Facts == nil {
		session.Facts = InferenceFacts{}
	}

	switch request.Step {

	case 1:
		people, ok := request.Data.(float64)

		if !ok {
			writeJSON(w, http.StatusBadRequest, map[string]any{
				"error": "people must be a number",
			})
			return
		}

		if people < 1 || people > 6 {
			writeJSON(w, http.StatusBadRequest, map[string]any{
				"error": "the number of people must be a number from 1 to 6",
			})
			return
		}

		session.Facts["people"] = int(people)

	case 2:
		budget, ok := request.Data.(float64)

		if !ok {
			writeJSON(w, http.StatusBadRequest, map[string]any{
				"error": "budget must be a number",
			})
			return
		}

		if budget < 1 || budget > 3 {
			writeJSON(w, http.StatusBadRequest, map[string]any{
				"error": "budget must be 1 (baixo), 2 (médio) or 3 (alto)",
			})
			return
		}

		session.Facts["budget"] = int(budget)

	case 3:
		focus, ok := request.Data.(string)

		if !ok {
			writeJSON(w, http.StatusBadRequest, map[string]any{
				"error": "focus must be a string",
			})
			return
		}

		session.Facts["focus"] = focus

	case 4:
		child, ok := request.Data.(bool)

		if !ok {
			writeJSON(w, http.StatusBadRequest, map[string]any{
				"error": "child must be a boolean",
			})
			return
		}

		session.Facts["child"] = child

	case 5:
		accommodation, ok := request.Data.(string)

		if !ok {
			writeJSON(w, http.StatusBadRequest, map[string]any{
				"error": "accommodation must be a string",
			})
			return
		}

		session.Facts["accommodation"] = accommodation
	}

	_, err = collection.UpdateOne(
		r.Context(),
		bson.M{"_id": sessionID},
		bson.M{
			"$set": bson.M{
				"facts": session.Facts,
			},
		},
	)

	if err != nil {
		log.Printf("mongo update error: %v", err)

		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "failed to update session",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"message": "data received successfully",
		"facts":   session.Facts,
	})
}


func (s *Server) inference(w http.ResponseWriter, r *http.Request) {
	request := r.URL.Query().Get("sessionId")
	if request == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": "sessionId is required",
		})
		return
	}

	sessionID, err := primitive.ObjectIDFromHex(request)

	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{
			"error": "invalid sessionId",
		})
		return
	}

	collection := s.db.Collection("inference_sessions")

	var session InferenceSession

	err = collection.FindOne(
		r.Context(),
		bson.M{
			"_id":    sessionID,
			"active": true,
		},
	).Decode(&session)

	if err != nil {

		if err == mongo.ErrNoDocuments {
			writeJSON(w, http.StatusNotFound, map[string]any{
				"error": "active session not found",
			})
			return
		}

		log.Printf("mongo find error: %v", err)

		writeJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "failed to find session",
		})
		return
	}

	result, rule := runInference(session.Facts)

	writeJSON(w, http.StatusOK, InferenceResponse{
		Inferred: result,
		Rule:     rule,
	})
}
