package repositories

import (
	"context"
	"fmt"

	"github.com/VaraPrasad27/shop/server/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultPageSize = 25
	maxPageSize     = 100
)

func GetAllProducts(ctx context.Context, dbpool *pgxpool.Pool, limit, offset int) ([]models.Product, error) {
	if limit <= 0 {
		limit = defaultPageSize
	}
	if limit > maxPageSize {
		limit = maxPageSize
	}
	if offset < 0 {
		offset = 0
	}

	rows, err := dbpool.Query(ctx,
		"SELECT id, name, description, price_cents, currency, image_url, stock, created_at FROM products ORDER BY id LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("querying products: %w", err)
	}
	defer rows.Close()

	products, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Product])
	if err != nil {
		return nil, fmt.Errorf("scanning products: %w", err)
	}
	if products == nil {
		products = []models.Product{}
	}
	return products, nil
}
