package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	mux := http.NewServeMux()
	Serv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	// Welcom to chairpy and image handler
	mux.Handle("/", http.FileServer(http.Dir(".")))
	mux.Handle("/assets", http.FileServer(http.Dir(".")))
	mux.HandleFunc("/ready", handlerReadiness)

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodPost)   // tells client only POST is allowed
			w.WriteHeader(http.StatusMethodNotAllowed) // sends 405
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(http.StatusOK)))
	})

	// the hits
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("/reset", apiCfg.handlerReset)
	//mux.HandleFunc("POST /articles", handlerArticlesCreate)
	err := Serv.ListenAndServe()

	fmt.Println("Starting the server <<<>>>")
	if err != nil {
		log.Fatalf("can not start the server : %s", err)
	}
}

// the main

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodPost)   // tells client only POST is allowed
		w.WriteHeader(http.StatusMethodNotAllowed) // sends 405
		return
	}
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
}
func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)   // tells client only POST is allowed
		w.WriteHeader(http.StatusMethodNotAllowed) // sends 405
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hits restored to 0 "))
}

//func handlerArticlesCreate(w http.ResponseWriter, r *http.Request) {
//
//}
