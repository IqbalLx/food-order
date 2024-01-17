package home

import (
	"context"

	"github.com/IqbalLx/food-order/src/shared/entities"
	"github.com/jackc/pgx/v5"
)

func doGetStores(ctx context.Context, db *pgx.Conn, size int) ([]entities.StoreWithCategories, bool, error) {
	stores, err := getStores(ctx, db, size); if err != nil {
		return stores, false, err
	}
	isScrollable, err := isStoresScrollable(ctx, db, size); if err != nil {
		return stores, false, err
	}

	return stores, isScrollable, nil
}