package main

import (
	"Chirpy/database"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	fileserverHits int
	db             *database.DB
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Could not get current directory")
		return
	}

	db, err := database.NewDB(cwd, *dbg)
	if err != nil {
		log.Fatal("Could not open database connection")
		return
	}

	cfg := apiConfig{
		fileserverHits: 0,
		db:             db,
	}

	fs := http.FileServer(http.Dir(filepathRoot))
	fsHandler := cfg.middlewareMetricsInc(http.StripPrefix("/app", fs))

	r := chi.NewRouter()
	r.Handle("/app", fsHandler)
	r.Handle("/app/*", fsHandler)

	api := chi.NewRouter()
	api.Get("/healthz", handlerReadiness)
	api.Get("/reset", cfg.handlerReset)
	api.Get("/chirps", cfg.getChirps)
	api.Post("/chirps", cfg.postChirp)
	api.Get("/chirps/{id}", cfg.getChirp)
	api.Post("/users", cfg.createUser)

	r.Mount("/api", api)

	admin := chi.NewRouter()
	admin.Get("/metrics", cfg.handlerMetrics)
	r.Mount("/admin", admin)

	corsMux := middlewareCors(r)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
