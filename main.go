package main

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/joho/godotenv"
)

func main() {
	// Variaveis de ambiente
	_ = godotenv.Load()
	mongoURI := getEnv("MONGO_URI", "")
	dbName := getEnv("MONGO_DB", "")
	port := getEnv("PORT", "")

    // Configuração do MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	tlsInsecure := getEnv("MONGO_TLS_INSECURE", "")
	clientOpts := options.Client().ApplyURI(mongoURI).SetServerSelectionTimeout(20 * time.Second)
	if tlsInsecure == "true" {
		clientOpts = clientOpts.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})
	}

	// Conexão com o MongoDB
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatalf("mongo connect error: %v", err)
	}

	// Ping do MongoDB para tentativas de reconexão
	var pingErr error
	for i := 0; i < 3; i++ {
		pingCtx, pingCancel := context.WithTimeout(context.Background(), 10*time.Second)
		pingErr = client.Ping(pingCtx, nil)
		pingCancel()
		if pingErr == nil {
			break
		}
		log.Printf("mongo ping attempt %d failed: %v", i+1, pingErr)
		time.Sleep(2 * time.Second)
	}

	if pingErr != nil {
		log.Fatalf("mongo ping failed after retries: %v", pingErr)
	}

	// Configuração do servidor HTTP
	server := &Server{
		db: client.Database(dbName),
	}

	// Configuração do roteador
	r := chi.NewRouter()

	// Configuração do CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"OPTIONS", "GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Rotas
	r.Get("/", server.info)

	r.Get("/v1/inference/session", server.InferenceSession)
	r.Post("/v1/inference/data", server.inferenceData)
	r.Get("/v1/inference", server.inference)
	log.Printf("listening on :%s", port)

	// Inicia o servidor HTTP
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}


// caso de uso e qual numero corresponde a situação
// quais os papeis existem