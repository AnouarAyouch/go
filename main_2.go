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
	// ALL THE  handlers
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.Handle("/app", fsHandler)
	mux.Handle("/app/assets", func(w http.ResponseWriter, r *http.Request) {
		apiCfg.middlewareMetricsInc(http.StripPrefix("/app/assets", http.FileServer(w, r, "./logo.png")))
	})

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET /api/metrics", apiCfg.handlerMetrics)
	//serving and listning to the server
	err := Serv.ListenAndServe()
	fmt.Println("Starting the server <<<>>>")
	if err != nil {
		log.Fatalf("can not start the server : %s", err)
	}
}

// themetricesfuncandhandthewrapper
func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
}
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

// thereadinesshandlefunc
func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// theresethandlefunction
func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hits restored to 0 "))
}

//func handlerArticlesCreate(w http.ResponseWriter, r *http.Request) {
//
//}
