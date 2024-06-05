package api

import (
	"awesome/config"
	"awesome/db"
	"log"
	"net/http"
	"time"
)

type apiService struct {
	mux *http.ServeMux
}

func NewAPIService() *apiService {
	return &apiService{mux: http.NewServeMux()}
}

func (s *apiService) AddEndpoints() {
	s.mux.HandleFunc("POST /api/v1/signin", handleSignin)
	s.mux.HandleFunc("POST /api/v1/signup", handleSignup)
}

func (s *apiService) Run() {
	server := &http.Server{
		Addr:         config.Config.Api.Address,
		Handler:      s.mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Print("Started API Service! 🚀")
	server.ListenAndServe()
}

func handleSignin(w http.ResponseWriter, r *http.Request) {

	var dbConfig = config.GetDatabaseConfiguration()
	mysqlDb, err := db.NewMySqlDb(dbConfig.ConnectionString)

	if err != nil {
		log.Print(err)
		http.Error(w, "Database Service Unavailable!", http.StatusServiceUnavailable)
		return
	}

	defer mysqlDb.Dispose()
}

func handleSignup(w http.ResponseWriter, r *http.Request) {
}
