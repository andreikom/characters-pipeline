package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"charactersync/pkg/characters"
	"charactersync/providers"
	"charactersync/store"
)

const requestTimeout = 5 * time.Minute

func main() {
	storeTempDir, err := os.MkdirTemp("", "charactersync_store_*")
	if err != nil {
		log.Fatalf("failed to create temp dir: %v", err)
	}
	fileStore := store.NewFileStore(fmt.Sprintf("%s/", storeTempDir))
	log.Printf("Temp filestore dir: %s", storeTempDir)

	fetchers := []providers.Fetcher{
		providers.NewPokemonFetcher(),
		providers.NewSWAPIFetcher(),
		providers.NewRickAndMortyFetcher(),
	}
	charactersService := characters.NewCharactersService(fetchers, fileStore)

	mux := http.NewServeMux()

	charactersHandlerFunc := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
		defer cancel()

		q := r.URL.Query().Get("limitPerProvider")
		limitPerProvider := 0 // zero = no limit
		if q != "" {
			if v, err := strconv.Atoi(q); err == nil && v > 0 {
				limitPerProvider = v
			} else {
				http.Error(w, "invalid limit parameter", http.StatusBadRequest)
				return
			}
		}

		allCharacters, err := charactersService.FetchAll(ctx, &providers.FetchOptions{LimitPerProvider: limitPerProvider})
		if err != nil {
			http.Error(w, "sync error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(allCharacters)
		if err != nil {
			http.Error(w, "response processing error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	mux.HandleFunc("/characters", charactersHandlerFunc)

	addr := ":8080"
	log.Printf("Server listening on %s", addr)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("failed server listenAndServe: %v", err)
	}
}
