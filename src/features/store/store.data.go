package store

import (
	"context"

	"github.com/IqbalLx/food-order/src/shared/entities"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leporo/sqlf"
)

func getStores(ctx context.Context, db *pgxpool.Pool, size int, lastStoreSecondaryID int) ([]entities.StoreWithCategories, error) {
	sqlf.SetDialect(sqlf.PostgreSQL)
	query := sqlf.
		From("stores as s").
		Join("store_categories as sc", "sc.store_id = s.id").
		Join("categories as c", "sc.category_id = c.id").
		GroupBy("s.id").
		Select("s.id, s.name, s.slug, s.image, s.short_desc, s.rating, s.secondary_id").
		Select("array_agg(c.name) as categories").
		Where("s.secondary_id > ?", lastStoreSecondaryID).
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

func isStoresScrollable(ctx context.Context, db *pgxpool.Pool, size int, lastStoreSecondaryID int) (bool, error) {
	sqlf.SetDialect(sqlf.PostgreSQL)
	query := sqlf.
		From("stores as s").
		Select("1 as one").
		Where("s.secondary_id > ?", lastStoreSecondaryID).
		Offset(size).
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

func isStoreExistsBySlug(ctx context.Context, db *pgxpool.Pool, slug string) (bool, error) {
	sqlf.SetDialect(sqlf.PostgreSQL)
	query := sqlf.
		From("stores as s").
		Select("1 as one").
		Where("s.slug = ?", slug).
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

func isStoreExistsByID(ctx context.Context, db *pgxpool.Pool, id string) (bool, error) {
	sqlf.SetDialect(sqlf.PostgreSQL)
	query := sqlf.
		From("stores as s").
		Select("1 as one").
		Where("s.id = ?", id).
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

func getStoreBySlug(ctx context.Context, db *pgxpool.Pool, slug string) (entities.StoreWithCategories, error) {
	sqlf.SetDialect(sqlf.PostgreSQL)
	query := sqlf.
		From("stores as s").
		Where("s.slug = ?", slug).
		Join("store_categories as sc", "sc.store_id = s.id").
		Join("categories as c", "sc.category_id = c.id").
		GroupBy("s.id").
		Select("s.id, s.name, s.image, s.short_desc, s.desc, s.rating").
		Select("array_agg(c.name) as categories")

	sql, args := query.String(), query.Args()
	row := db.QueryRow(ctx, sql, args...)

	var store entities.StoreWithCategories
	err := row.Scan(&store.ID, &store.Name, &store.Image, 
		&store.ShortDesc, &store.Desc, &store.Rating, &store.Categories)
	if err != nil {
		return store, err
	}

	return store, nil
}

func getStoreByID(ctx context.Context, db *pgxpool.Pool, id string) (entities.StoreWithCategories, error) {
	sqlf.SetDialect(sqlf.PostgreSQL)
	query := sqlf.
		From("stores as s").
		Where("s.id = ?", id).
		Join("store_categories as sc", "sc.store_id = s.id").
		Join("categories as c", "sc.category_id = c.id").
		GroupBy("s.id").
		Select("s.id, s.name, s.image, s.short_desc, s.desc, s.rating").
		Select("array_agg(c.name) as categories")

	sql, args := query.String(), query.Args()
	row := db.QueryRow(ctx, sql, args...)

	var store entities.StoreWithCategories
	err := row.Scan(&store.ID, &store.Name, &store.Image, 
		&store.ShortDesc, &store.Desc, &store.Rating, &store.Categories)
	if err != nil {
		return store, err
	}

	return store, nil
}

func getStoreMenusByStoreID(ctx context.Context, db *pgxpool.Pool, cartID string, storeID string, size int, lastMenuSecondaryID int,
	isWithCategory bool, menuCategoryID string, withSearchQuery bool, searchQuery string) ([]entities.StoreMenuWithQuantity, error) {
	sqlf.SetDialect(sqlf.PostgreSQL)
	query := sqlf.
		From("store_menus as sm").
		Clause(
			`LEFT JOIN cart_items as ci ON 
				ci.store_menu_id = sm.id AND 
				ci.store_id = sm.store_id AND`,
		).
		Expr("ci.cart_id = ?", cartID).
		Where("sm.store_id = ?", storeID).
		Where("sm.secondary_id > ?", lastMenuSecondaryID).
		OrderBy("sm.secondary_id ASC").
		Limit(size).
		Select(
			`sm.id,
			 sm.secondary_id,
			 sm.name,
			 sm.image,
			 sm.price,
			 sm.ordered_count,
			 sm.price_promo,
			 sm.is_available`,
		).
		Select("COALESCE(ci.quantity, 0) as quantity")

	if (isWithCategory) {
		query.
			Join("store_menu_category_items as smci", "smci.store_menu_id = sm.id").
			Where("smci.store_menu_category_id = ?", menuCategoryID)
	}

	if (withSearchQuery) {
		query.
			Where("sm.name ILIKE ?", "%" + searchQuery + "%")
	}

	sql, args := query.String(), query.Args()
	rows, err := db.Query(ctx, sql, args...); if err != nil {
		return []entities.StoreMenuWithQuantity{}, err
	}
	menus, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (entities.StoreMenuWithQuantity, error) {
		var menu entities.StoreMenuWithQuantity
		err = row.Scan(
			&menu.ID,
			&menu.SecondaryID,
			&menu.Name,
			&menu.Image,
			&menu.Price,
			&menu.OrderedCount,
			&menu.PricePromo,
			&menu.IsAvailable,
			&menu.Quantity,
		)
		if err != nil {
			return menu, err
		}

		return menu, err
	})

	if err != nil{
		return []entities.StoreMenuWithQuantity{}, err
	}

	return menus, nil
}

func isMenusScrollable(ctx context.Context, db *pgxpool.Pool, storeID string, size int, lastMenuSecondaryID int,
	isWithCategory bool, menuCategoryID string, withSearchQuery bool, searchQuery string) (bool, error) {
	sqlf.SetDialect(sqlf.PostgreSQL)
	query := sqlf.
		From("store_menus as sm").
		Select("1 as one").
		Where("sm.store_id = ?", storeID).
		Where("sm.secondary_id > ?", lastMenuSecondaryID).
		Offset(size).
		OrderBy("sm.secondary_id ASC").
		Limit(1)

	if (isWithCategory) {
		query.
			Join("store_menu_category_items as smci", "smci.store_menu_id = sm.id").
			Where("smci.store_menu_category_id = ?", menuCategoryID)
	}

	if (withSearchQuery) {
		query.
			Where("sm.name ILIKE ?", "%" + searchQuery + "%")
	}

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

func getStoreMenuCategories(ctx context.Context, db *pgxpool.Pool, storeID string, 
	withSearchQuery bool, searchQuery string) ([]entities.StoreMenuCategory, error) {
	sqlf.SetDialect(sqlf.PostgreSQL)
	query := sqlf.
		From("store_menu_categories as smc")

	if (withSearchQuery) {
		query.
			Clause(
				`INNER JOIN store_menu_category_items as smci ON 
					smci.store_menu_category_id = smc.id AND`,
			).
			Expr("smci.store_id = ?", storeID).
			Clause(
				`INNER JOIN store_menus as sm ON 
					sm.id = smci.store_menu_id AND`,
			).
			Expr("sm.name ILIKE ?", "%" + searchQuery + "%")
	}

	query.
		Where("smc.store_id = ?", storeID).
		OrderBy("smc.name ASC").
		Select("smc.id, smc.name")

	sql, args := query.String(), query.Args()
	rows, err := db.Query(ctx, sql, args...); if err != nil {
		return []entities.StoreMenuCategory{}, err
	}
	storeMenuCategories, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (entities.StoreMenuCategory, error) {
		var storeMenuCategory entities.StoreMenuCategory
		err = row.Scan(
			&storeMenuCategory.ID,
			&storeMenuCategory.Name,
		)
		if err != nil {
			return storeMenuCategory, err
		}

		return storeMenuCategory, err
	})

	if err != nil{
		return []entities.StoreMenuCategory{}, err
	}

	return storeMenuCategories, nil
}

func countCartItems(ctx context.Context, db *pgxpool.Pool, cartID string) (int, error) {
	sqlf.SetDialect(sqlf.PostgreSQL)
	query := sqlf.
		From("cart_items").
		Where("cart_id = ?", cartID).
		Select("COALESCE(SUM(quantity), 0) as quantity")
	
	sql, args := query.String(), query.Args()
	row := db.QueryRow(ctx, sql, args...)

	var quantity int
	err := row.Scan(&quantity); if err != nil {
		return quantity, err
	}

	return quantity, nil
}