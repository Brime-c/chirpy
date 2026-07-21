package main

import (
	"chirpy/internal/database"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbqueries      database.Queries
	platform       string
}
type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}
type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	var apiCfg apiConfig
	apiCfg.dbqueries = *dbQueries
	apiCfg.platform = os.Getenv("PLATFORM")
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("."))
	prefixed := http.StripPrefix("/app", fileServer)
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(prefixed))
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("POST /api/users", func(w http.ResponseWriter, r *http.Request) {
		type Request struct {
			Email string `json:"email"`
		}
		decoder := json.NewDecoder(r.Body)
		req := Request{}
		err := decoder.Decode(&req)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		user, err := dbQueries.CreateUser(r.Context(), req.Email)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		newUser := User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		}
		respondWithJSON(w, 201, newUser)
	})
	mux.HandleFunc("POST /api/chirps", func(w http.ResponseWriter, r *http.Request) {
		type response struct {
			Body   string    `json:"body"`
			UserID uuid.UUID `json:"user_id"`
		}
		decoder := json.NewDecoder(r.Body)
		resp := response{}
		err := decoder.Decode(&resp)
		if err != nil {
			respondWithError(w, 500, "Something went wrong")
			return
		}
		if len(resp.Body) > 140 {
			respondWithError(w, 400, "Chirp is too long")
			return
		}
		cleanedBody := getCleanedBody(resp.Body)

		chirp, err := dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{Body: cleanedBody, UserID: resp.UserID})
		if err != nil {
			respondWithError(w, 500, "error trying to create chirp")
			return
		}
		newChirp := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
		respondWithJSON(w, 201, newChirp)
	})
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	myServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(myServer.ListenAndServe())
}
