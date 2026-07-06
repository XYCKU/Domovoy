package storage

import (
	"Domovoy/realt"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, dsn string) (*Store, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return &Store{pool: pool}, nil
}

func (s *Store) Close() { s.pool.Close() }

func (s *Store) GetSeenUUIDs(ctx context.Context) (map[string]bool, error) {
	rows, err := s.pool.Query(ctx, "SELECT uuid FROM listings")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	seen := make(map[string]bool)
	for rows.Next() {
		var uuid string
		if err := rows.Scan(&uuid); err != nil {
			return nil, err
		}
		seen[uuid] = true
	}
	return seen, nil
}

func (s *Store) Upsert(ctx context.Context, o realt.Object) (isNew bool, err error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			log.Printf("failed to rollback transaction: %v", err)
		}
	}()

	var existingPrice *float64
	err = tx.QueryRow(ctx, `SELECT price FROM listings WHERE uuid=$1`, o.UUID).Scan(&existingPrice)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		isNew = true
	case err != nil:
		return false, err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO listings(
			uuid, source,code, category,
			price, price_currency, price_per_m2,
			rooms, area_total, area_living, area_kitchen,
			storey, storeys, building_year,
			town_name, district_name, street_name, house_number, address,
			title, description, images,
			source_created_at, source_updated_at, raise_date,
			last_seen_at
		) VALUES (
			$1, $2, $3, $4,
			$5, $6, $7,
			$8, $9, $10, $11,
			$12, $13, $14,
			$15, $16, $17, $18, $19,
			$20, $21, $22,
			$23, $24, $25,
			now()
		)
		ON CONFLICT (uuid) DO UPDATE SET
			price = EXCLUDED.price,
			price_currency = EXCLUDED.price_currency,
			price_per_m2 = EXCLUDED.price_per_m2,
			source_updated_at = EXCLUDED.source_updated_at,
			raise_date = EXCLUDED.raise_date,
			last_seen_at = now()
	`,
		o.UUID, "realt.by", o.Code, o.Category,
		o.Price, o.PriceCurrency, o.PricePerM2,
		o.Rooms, o.AreaTotal, o.AreaLiving, o.AreaKitchen,
		o.Storey, o.Storeys, o.BuildingYear,
		o.TownName, o.DistrictName, o.StreetName, o.HouseNumber, o.Address,
		o.Title, o.Description, o.Images,
		parseTime(o.CreatedAt), parseTime(o.UpdatedAt), parseTime(o.RaiseDate),
	)
	if err != nil {
		return false, fmt.Errorf("failed to upsert object: %w", err)
	}

	if isNew || (existingPrice != nil && *existingPrice != o.Price) {
		_, err = tx.Exec(ctx, `INSERT INTO price_history (uuid, price, currency) VALUES ($1, $2, $3)`, o.UUID, o.Price, o.PriceCurrency)
		if err != nil {
			return false, fmt.Errorf("failed to insert price history: %w", err)
		}
	}

	return isNew, tx.Commit(ctx)
}

func parseTime(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil
	}
	return &t
}
