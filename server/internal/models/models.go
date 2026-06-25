package models

import (
	"time"
)

type Product struct {
	ID          string    `db:"id"           json:"id"`
	Name        string    `db:"name"         json:"name"`
	Description string    `db:"description"  json:"description"`
	PriceCents  int64     `db:"price_cents"  json:"price_cents"`
	Currency    string    `db:"currency"     json:"currency"`
	ImageURL    string    `db:"image_url"    json:"image_url"`
	Stock       int       `db:"stock"        json:"stock"`
	CreatedAt   time.Time `db:"created_at"   json:"created_at"`
}
