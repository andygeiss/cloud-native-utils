package api

import (
	"cloud-native/services/key-value/internal/app/adapters/common/config"
	"net/http"
)

func Route(cfg *config.Config) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/v1/store", Delete(cfg))
	mux.HandleFunc("GET /api/v1/store", Get(cfg))
	mux.HandleFunc("PUT /api/v1/store", Put(cfg))
	return mux
}
