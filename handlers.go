package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/mjh1207/chirpy/internal/database"
)

func handlerReadiness(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Unable to complete reset", nil)
		return
	}
	cfg.fileserverHits.Store(0)
	cfg.db.DeleteUsers(req.Context())
	w.WriteHeader(http.StatusOK)
}

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.CreateUser(req.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create new user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	})
}

func (cfg *apiConfig) handlerPostChirps(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
		User_Id uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	chirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{
		Body: replaceProfanity(params.Body),
		UserID: params.User_Id,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create new chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp {
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		User_Id: chirp.UserID.String(),
	})

}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.db.GetAllChirps(req.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps from database", err)
		return
	}

	var chirpsSlice []Chirp
	for _, chirp := range chirps {
		chirpsSlice = append(chirpsSlice, Chirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			User_Id: chirp.UserID.String(),
		})
	}
	respondWithJSON(w, http.StatusOK, chirpsSlice)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("chirpID")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Not a valid uuid", err)
		return
	}
	chirp, err := cfg.db.GetChirp(req.Context(), parsedId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Unable to get chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		User_Id: chirp.UserID.String(),
	})
}