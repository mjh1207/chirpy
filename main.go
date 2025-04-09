package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mjh1207/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
	platform string
	jwtSecret string
	polkaKey string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET must be set")
	}

	pKey := os.Getenv("POLKA_KEY")
	if pKey == "" {
		log.Fatal("POLKA_KEY must be set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("unable to open sql connection: %v", err)
	}

	const port = "8080"
	apiCfg := apiConfig {
		fileserverHits: atomic.Int32{},
		db: database.New(dbConn),
		platform: os.Getenv("PLATFORM"),
		jwtSecret: secret,
		polkaKey: pKey,
	}


	mux := http.NewServeMux()
	handler := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsers)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUser)

	mux.HandleFunc("POST /api/chirps", apiCfg.handlerPostChirps)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirp)
	
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)

	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerWebhooks)

	server := &http.Server {
		Handler: mux,
		Addr: ":" + port,
	}


	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
