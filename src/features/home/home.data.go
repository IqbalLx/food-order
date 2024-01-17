package home

import (
	"context"

	"github.com/IqbalLx/food-order/src/shared/entities"
	"github.com/jackc/pgx/v5"
	"github.com/leporo/sqlf"
)

func getStores(ctx context.Context, db *pgx.Conn, size int) ([]entities.StoreWithCategories, error) {
	sqlf.SetDialect(sqlf.PostgreSQL)
	query := sqlf.
		From("stores as s").
		Join("store_categories as sc", "sc.store_id = s.id").
		Join("categories as c", "sc.category_id = c.id").
		GroupBy("s.id").
		Select("s.id, s.name, s.slug, s.image, s.short_desc, s.rating, s.secondary_id").
		Select("array_agg(c.name) as categories").
		OrderBy("s.secondary_id ASC").
		Limit(size)

	sql, args := query.String(), query.Args()
	rows, err := db.Query(ctx, sql, args...); if err != nil {
		return []entities.StoreWithCategories{}, err
	}
	stores, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (entities.StoreWithCategories, error) {
		var store entities.StoreWithCategories
		err = row.Scan(&store.ID, &store.Name, &store.Slug, &store.Image, 
			&store.ShortDesc, &store.Rating, &store.SecondaryID, &store.Categories)
		if err != nil {
			return store, err
		}

		return store, err
	})

	if err != nil{
		return []entities.StoreWithCategories{}, err
	}

	return stores, nil
}

func isStoresScrollable(ctx context.Context, db *pgx.Conn, size int) (bool, error) {
	sqlf.SetDialect(sqlf.PostgreSQL)
	query := sqlf.
		From("stores as s").
		Select("1 as one").
		Where("s.secondary_id > ?", size - 1).
		OrderBy("s.secondary_id ASC").
		Limit(1)

	sql, args := query.String(), query.Args()
	row := db.QueryRow(ctx, sql, args...)
	
	var one string
	err := row.Scan(&one); if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return false, nil
		}

		return false, err
	}

	return true, nil
}