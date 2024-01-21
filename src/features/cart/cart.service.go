package cart

import (
	"context"
	"errors"

	"github.com/IqbalLx/food-order/src/shared/entities"
	"github.com/jackc/pgx/v5/pgxpool"
)

func doUpsertMenuToCart(ctx context.Context, db *pgxpool.Pool, cartID string, menuID string, quantity int) (int, error) {
	isMenuExists, err := isMenuExists(ctx, db, menuID); if err != nil {
		return 0, err
	}
	if !isMenuExists {
		return 0, errors.New("menu not found")
	}

	menu, err := getMenuByID(ctx, db, menuID); if err != nil {
		return 0, err
	}

	if (quantity == 0) {
		err = deleteMenuFromCart(ctx, db, cartID, menuID)
		return 0, err
	}

	finalQuantity, err := upsertMenuToCart(ctx, db, cartID, quantity, menu); if err != nil {
		return 0, err
	}

	return finalQuantity, nil
}

func doCountCartItems(ctx context.Context, db *pgxpool.Pool, cartID string) (int, error) {
	cartCount, err := countCartItems(ctx, db, cartID); if err != nil {
		return 0, err
	}

	return cartCount, nil
}

func doDeleteAllMenusFromCartByStoreID(ctx context.Context, db *pgxpool.Pool, cartID string, storeID string) error {
	return deleteAllMenusFromCartByStoreID(ctx, db, cartID, storeID)
}

func doDeleteMenuFromCart(ctx context.Context, db *pgxpool.Pool, cartID string, storeID string, menuID string) (int, error) {
	err := deleteMenuFromCart(ctx, db, cartID, menuID); if err != nil {
		return 0, err
	}

	storeQuantityAfterDelete, err := countCartItemsByStoreID(ctx, db, cartID, storeID); if err != nil {
		return 0, err
	}

	return storeQuantityAfterDelete, nil
}

func doGetCart(ctx context.Context, db *pgxpool.Pool, cartID string) ([]entities.StoreWithCartMenus, int, int, error) {
	var data []entities.StoreWithCartMenus
	
	stores, err := getStoresByCartID(ctx, db, cartID); if err != nil {
		return data, 0, 0, err
	}

	for _, store := range stores {
		data = append(data, entities.StoreWithCartMenus{
			Store: entities.Store{
				ID: store.ID,
				Slug: store.Slug,
				Name: store.Name,
			},
		})
	}

	menus, err := getMenusInCart(ctx, db, cartID); if err != nil {
		return data, 0, 0, err
	}

	err = populateCartMenus(&data, menus); if err != nil {
		return data, 0, 0, err
	}

	countItems, err := countCartItems(ctx, db, cartID); if err != nil {
		return data, 0, 0, err
	}
	totalItems, err := sumCartTotal(ctx, db, cartID); if err != nil {
		return data, 0, 0, err
	}

	return data, countItems, totalItems, nil
}

type CartState struct {
	CountMenus int
	TotalItems int
	CountStores int
}
func doGetCartState(ctx context.Context, db *pgxpool.Pool, cartID string) (CartState, error) {
	var state CartState

	countMenus, err := countCartMenus(ctx, db, cartID); if err != nil {
		return state, err
	}
	totalItems, err := sumCartTotal(ctx, db, cartID); if err != nil {
		return state, err
	}
	countStores, err := countCartStores(ctx, db, cartID); if err != nil {
		return state, err
	}

	state.CountMenus = countMenus
	state.TotalItems = totalItems
	state.CountStores = countStores

	return state, nil
}

func doGetCartStateByStoreID(ctx context.Context, db *pgxpool.Pool, cartID string, storeID string) (int, int, error) {
	countItems, err := countCartItemsByStoreID(ctx, db, cartID, storeID); if err != nil {
		return 0, 0, err
	}

	sumItems, err := sumCartItemsSubtotalByStoreID(ctx, db, cartID, storeID); if err != nil {
		return 0, 0, err
	}

	return countItems, sumItems, nil
}