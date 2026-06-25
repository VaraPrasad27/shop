package services

import (
	"context"

	"github.com/VaraPrasad27/shop/server/internal/models"
	"github.com/VaraPrasad27/shop/server/internal/repositories"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetAllProducts(ctx context.Context, pool *pgxpool.Pool, limit, offset int) ([]models.Product, error) {
	return repositories.GetAllProducts(ctx, pool, limit, offset)
}
