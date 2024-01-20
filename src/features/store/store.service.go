package store

import (
	"context"
	"errors"

	"github.com/IqbalLx/food-order/src/shared/entities"
	"github.com/gofiber/fiber/v2/log"
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

func doGetInitialStoreDetail(ctx context.Context, db *pgxpool.Pool, cartdID string, slug string, menuSize int,
	isWithSearchQuery bool, searchQuery string) (StoresWithMenus, error) {
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
	menuCategories, err := getStoreMenuCategories(ctx, db, store.ID, isWithSearchQuery, searchQuery); if err != nil {
		return data, err
	}
	menus, err := getStoreMenusByStoreID(ctx, db, cartdID, store.ID, menuSize, 0, false, "", isWithSearchQuery, searchQuery); if err != nil {
		return data, err
	}
	isMenusScrollable, err := isMenusScrollable(ctx, db, store.ID, menuSize, 0, false, "", isWithSearchQuery, searchQuery); if err != nil {
		return data, err
	}

	data.MenuCategories = menuCategories
	data.Store = store
	data.Menus = menus
	data.IsMenusScrollable = isMenusScrollable

	return data, nil
}

func doGetMenus(ctx context.Context, db *pgxpool.Pool, cartID string, storeID string, menuSize int, lastMenuSecondaryID int,
	isWithCategory bool, menuCategoryID string, isWithSearchQuery bool, searchQuery string) (StoresWithMenus, error) {
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

	menus, err := getStoreMenusByStoreID(ctx, db, cartID, storeID, menuSize, lastMenuSecondaryID, isWithCategory, 
		menuCategoryID, isWithSearchQuery, searchQuery); if err != nil {
		return data, err
	}
	isMenusScrollable, err := isMenusScrollable(ctx, db, storeID, menuSize, lastMenuSecondaryID, isWithCategory, 
		menuCategoryID, isWithSearchQuery, searchQuery); if err != nil {
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

func doGetMenuCategories(ctx context.Context, db *pgxpool.Pool, storeID string, isWithSearchQuery bool, searchQuery string) (StoresWithMenus, error) {
	var data StoresWithMenus
	store, err := getStoreByID(ctx, db, storeID); if err != nil {
		return data, err
	}
	menuCategories, err := getStoreMenuCategories(ctx, db, storeID, isWithSearchQuery, searchQuery); if err != nil {
		return data, err
	}

	data.Store = store
	data.MenuCategories = menuCategories

	return data, nil
}

func doSearchStoresByMenuName(ctx context.Context, db *pgxpool.Pool, cartID string, menuName string, page int, pageSize int) ([]entities.StoreWithMatchingMenu, int, error) {
	stores, maxPage, err := searchStoresByMenuName(ctx, db, menuName, page, pageSize); if err != nil {
		return []entities.StoreWithMatchingMenu{}, 0, err
	}

	if len(stores) == 0 {
		return []entities.StoreWithMatchingMenu{}, 0, nil
	}

	storeIDS := make([]interface{}, len(stores))
	for i, store := range stores {
		storeIDS[i] = store.ID
	}

	topMatchCount := 3
	menus, err := getTopMacthingMenuFromStores(ctx, db, cartID, menuName, storeIDS, topMatchCount); if err != nil {
		log.Debug(err.Error())
		return []entities.StoreWithMatchingMenu{}, 0, err
	}

	if err = populateMatcingMenus(&stores, menus); err != nil {
		return []entities.StoreWithMatchingMenu{}, 0, err
	}

	return stores, maxPage, nil
}

func doGetMenuByID(ctx context.Context, db *pgxpool.Pool, cartID string, storeID string, menuID string) (entities.StoreMenuWithQuantity, error) {
	isExists, err := isMenuExistsInStore(ctx, db, storeID, menuID); if err != nil {
		return entities.StoreMenuWithQuantity{}, err
	}

	if (!isExists) {
		return entities.StoreMenuWithQuantity{}, errors.New("menu doesnt exists")
	}

	menu, err := getStoreMenuByID(ctx, db, cartID, menuID); if err != nil {
		return entities.StoreMenuWithQuantity{}, err
	}

	return menu, nil
}