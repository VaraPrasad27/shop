package routes

import (
	"net/http"

	"github.com/VaraPrasad27/shop/server/internal/handlers"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(pool *pgxpool.Pool) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		_, _ = w.Write([]byte("shop api"))
	})

	r.Get("/products", handlers.GetAllProductsHandler(pool))

	return r
}
