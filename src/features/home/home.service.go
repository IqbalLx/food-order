package home

import (
	"context"

	"github.com/IqbalLx/food-order/src/shared/entities"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HomeInitialData struct {
	Stores []entities.StoreWithCategories
	IsStoresScrollable bool
}

func doGetHomeInitialData(ctx context.Context, db *pgxpool.Pool, size int) (HomeInitialData, error) {
	var data HomeInitialData
	
	stores, err := getStores(ctx, db, size); if err != nil {
		return data, err
	}
	isStoresScrollable, err := isStoresScrollable(ctx, db, size); if err != nil {
		return data, err
	}

	data.Stores = stores
	data.IsStoresScrollable = isStoresScrollable

	return data, nil
}