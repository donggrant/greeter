package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
)

// Global stats for the server
var (
	globalStats Stats
	statsMutex  sync.RWMutex
)

type GreetingRequest struct {
	Name     string `json:"name"`
	Language string `json:"language"`
}

type GreetingResponse struct {
	Greeting string `json:"greeting"`
	Stats    *Stats `json:"stats,omitempty"`
}

// updateGlobalStats updates the global stats with the latest translation stats
func updateGlobalStats(newStats *Stats) {
	if newStats == nil {
		return
	}

	statsMutex.Lock()
	defer statsMutex.Unlock()

	globalStats.APICalls += newStats.APICalls
	globalStats.CharsSent += newStats.CharsSent
	globalStats.CostEstimate += newStats.CostEstimate
	globalStats.CacheHits += newStats.CacheHits
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func RunServer() {
	// Serve static files from the frontend/dist directory
	fs := http.FileServer(http.Dir("frontend/dist"))
	http.Handle("/", fs)

	// API endpoint for greetings
	http.HandleFunc("/api/greet", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		language := r.URL.Query().Get("language")

		if name == "" || language == "" {
			http.Error(w, "Missing name or language parameter", http.StatusBadRequest)
			return
		}

		greeter, err := NewGreeter(name)
		if err != nil {
			http.Error(w, "Failed to create greeter: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer greeter.Close()

		greeter.SetLanguage(Language(language))
		greeting, err := greeter.Greet()
		if err != nil {
			http.Error(w, "Failed to get greeting: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Update global stats with the latest translation stats
		updateGlobalStats(greeter.stats)

		response := GreetingResponse{
			Greeting: greeting,
		}

		// Only include stats if there was an API call or cache hit
		if greeter.stats.APICalls > 0 || greeter.stats.CacheHits > 0 {
			response.Stats = greeter.stats
			log.Printf("Stats: calls=%d, chars=%d, cost=%.5f, hits=%d",
				greeter.stats.APICalls, greeter.stats.CharsSent, greeter.stats.CostEstimate, greeter.stats.CacheHits)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Wrap all handlers with CORS middleware
	handler := corsMiddleware(http.DefaultServeMux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}
