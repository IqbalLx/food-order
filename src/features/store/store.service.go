package store

import (
	"context"
	"errors"

	"github.com/IqbalLx/food-order/src/shared/entities"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StoresWithMenus struct {
	Store entities.StoreWithCategories
	MenuCategories []entities.StoreMenuCategory
	Menus []entities.StoreMenuWithQuantity
	IsMenusScrollable bool
}

func doGetStores(ctx context.Context, db *pgxpool.Pool, size int, lastStoreSecondaryID int) ([]entities.StoreWithCategories, bool, error) {
	stores, err := getStores(ctx, db, size, lastStoreSecondaryID); if err != nil {
		return stores, false, err
	}
	isScrollable, err := isStoresScrollable(ctx, db, size, lastStoreSecondaryID); if err != nil {
		return stores, false, err
	}

	return stores, isScrollable, nil
}

func doGetInitialStoreDetail(ctx context.Context, db *pgxpool.Pool, cartdID string, slug string, menuSize int) (StoresWithMenus, error) {
	var data StoresWithMenus

	isExists, err := isStoreExistsBySlug(ctx, db, slug); if err != nil {
		return data, err
	}

	if !isExists {
		return data, errors.New("store not found")
	}

	store, err := getStoreBySlug(ctx, db, slug); if err != nil {
		return data, err
	}
	menuCategories, err := getStoreMenuCategories(ctx, db, store.ID); if err != nil {
		return data, err
	}
	menus, err := getStoreMenusByStoreID(ctx, db, cartdID, store.ID, menuSize, 0, false, ""); if err != nil {
		return data, err
	}
	isMenusScrollable, err := isMenusScrollable(ctx, db, store.ID, menuSize, 0, false, ""); if err != nil {
		return data, err
	}

	data.MenuCategories = menuCategories
	data.Store = store
	data.Menus = menus
	data.IsMenusScrollable = isMenusScrollable

	return data, nil
}

func doGetMenus(ctx context.Context, db *pgxpool.Pool, cartID string, storeID string, menuSize int, lastMenuSecondaryID int,
	isWithCategory bool, menuCategoryID string) (StoresWithMenus, error) {
	var data StoresWithMenus
	
	isStoreExists, err := isStoreExistsByID(ctx, db, storeID); if err != nil {
		return data, err
	}
	if !isStoreExists {
		return data, err
	}

	store, err := getStoreByID(ctx, db, storeID); if err != nil {
		return data, err
	}

	menus, err := getStoreMenusByStoreID(ctx, db, cartID, storeID, menuSize, lastMenuSecondaryID, isWithCategory, menuCategoryID); if err != nil {
		return data, err
	}
	isMenusScrollable, err := isMenusScrollable(ctx, db, storeID, menuSize, lastMenuSecondaryID, isWithCategory, menuCategoryID); if err != nil {
		return data, err
	}

	data.Menus = menus
	data.Store = store
	data.IsMenusScrollable = isMenusScrollable

	return data, nil
}

func doCheckIfReadyToCheckout(ctx context.Context, db *pgxpool.Pool, cartID string) (bool, error) {
	cartCount, err := countCartItems(ctx, db, cartID); if err != nil {
		return false, err
	}

	return cartCount > 0, nil
}