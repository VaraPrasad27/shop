package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/VaraPrasad27/shop/server/internal/services"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetAllProductsHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, offset := paginationFromQuery(r)

		products, err := services.GetAllProducts(r.Context(), pool, limit, offset)
		if err != nil {
			log.Printf("get all products: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(products); err != nil {
			log.Printf("encode products: %v", err)
		}
	}
}

func paginationFromQuery(r *http.Request) (limit, offset int) {
	limit, _ = strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ = strconv.Atoi(r.URL.Query().Get("offset"))
	return
}
