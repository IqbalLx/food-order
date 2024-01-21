package cart

import (
	"context"

	"github.com/IqbalLx/food-order/src/shared/entities"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leporo/sqlf"
)

func isMenuExists(ctx context.Context, db *pgxpool.Pool, menuID string) (bool, error) {
	sqlf.SetDialect(sqlf.PostgreSQL)
	query := sqlf.
		From("store_menus").
		Select("1 as one").
		Where("id = ?", menuID).
		Where("is_available = ?", true).
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

func getMenuByID(ctx context.Context, db *pgxpool.Pool, menuID string) (entities.StoreMenu, error) {
	sqlf.SetDialect(sqlf.PostgreSQL)
	query := sqlf.
		From("store_menus").
		Select("id, store_id, price, price_promo").
		Where("id = ?", menuID).
		Limit(1)

	sql, args := query.String(), query.Args()
	row := db.QueryRow(ctx, sql, args...)

	var menu entities.StoreMenu
	err := row.Scan(&menu.ID, &menu.StoreID, &menu.Price, &menu.PricePromo); if err != nil {
		return menu, err
	}

	return menu, nil
}

func upsertMenuToCart(ctx context.Context, db *pgxpool.Pool, cartID string, 
	quantity int, menu entities.StoreMenu) (int, error) {
	sqlf.SetDialect(sqlf.PostgreSQL)
	query := sqlf.
		InsertInto("cart_items").
		NewRow().
			Set("cart_id", cartID).
			Set("store_id", menu.StoreID).
			Set("store_menu_id", menu.ID).
			Set("quantity", quantity).
			SetExpr("updated_at", "NOW()")

	if menu.PricePromo > 0 {
		query.Set("subtotal", quantity * menu.PricePromo)
	} else {
		query.Set("subtotal", quantity * menu.Price)
	}

	query.
		Clause("ON CONFLICT ON CONSTRAINT cart_items_unique_constraint DO UPDATE SET").
			Expr("quantity = EXCLUDED.quantity").
			Expr("subtotal = EXCLUDED.subtotal").
			Expr("updated_at = EXCLUDED.updated_at").
	Returning("quantity")

	sql, args := query.String(), query.Args()
	row := db.QueryRow(ctx, sql, args...)

	var finalQuantity int
	err := row.Scan(&finalQuantity); if err != nil {
		return finalQuantity, err
	}

	return finalQuantity, nil
}

func deleteMenuFromCart(ctx context.Context, db *pgxpool.Pool, cartID string, menuID string) error {
	sqlf.SetDialect(sqlf.PostgreSQL)
	query := sqlf.
		DeleteFrom("cart_items").
		Where("cart_id = ?", cartID).
		Where("store_menu_id = ?", menuID)
	
	sql, args := query.String(), query.Args()
	_, err := db.Exec(ctx, sql, args...); if err != nil {
		return err
	}
	
	return nil
}

func countCartItems(ctx context.Context, db *pgxpool.Pool, cartID string) (int, error) {
	sqlf.SetDialect(sqlf.PostgreSQL)
	query := sqlf.
		From("cart_items").
		Where("cart_id = ?", cartID).
		Where("quantity > 0").
		Select("COALESCE(SUM(quantity), 0) as quantity")
	
	sql, args := query.String(), query.Args()
	row := db.QueryRow(ctx, sql, args...)

	var quantity int
	err := row.Scan(&quantity); if err != nil {
		return quantity, err
	}

	return quantity, nil
}

func countCartMenus(ctx context.Context, db *pgxpool.Pool, cartID string) (int, error) {
	sqlf.SetDialect(sqlf.PostgreSQL)
	query := sqlf.
		From("cart_items").
		Where("cart_id = ?", cartID).
		Where("quantity > 0").
		Select("COALESCE(COUNT(DISTINCT store_menu_id), 0) as quantity")
	
	sql, args := query.String(), query.Args()
	row := db.QueryRow(ctx, sql, args...)

	var quantity int
	err := row.Scan(&quantity); if err != nil {
		return quantity, err
	}

	return quantity, nil
}

func countCartStores(ctx context.Context, db *pgxpool.Pool, cartID string) (int, error) {
	sqlf.SetDialect(sqlf.PostgreSQL)
	query := sqlf.
		From("cart_items").
		Where("cart_id = ?", cartID).
		Where("quantity > 0").
		Select("COALESCE(COUNT(DISTINCT store_id), 0) as store_count")

	sql, args := query.String(), query.Args()
	row := db.QueryRow(ctx, sql, args...)

	var count int
	err := row.Scan(&count); if err != nil {
		return count, err
	}

	return count, nil
}

func sumCartTotal(ctx context.Context, db *pgxpool.Pool, cartID string) (int, error) {
	sqlf.SetDialect(sqlf.PostgreSQL)
	query := sqlf.
		From("cart_items").
		Where("cart_id = ?", cartID).
		Select("COALESCE(SUM(subtotal), 0) as subtotal")
	
	sql, args := query.String(), query.Args()
	row := db.QueryRow(ctx, sql, args...)

	var subtotal int
	err := row.Scan(&subtotal); if err != nil {
		return subtotal, err
	}

	return subtotal, nil
}

func countCartItemsByStoreID(ctx context.Context, db *pgxpool.Pool, cartID string, storeID string) (int, error) {
	sqlf.SetDialect(sqlf.PostgreSQL)
	query := sqlf.
		From("cart_items").
		Where("cart_id = ?", cartID).
		Where("store_id = ?", storeID).
		Where("quantity > 0").
		Select("COALESCE(COUNT(*), 0) as count")
	
	sql, args := query.String(), query.Args()
	row := db.QueryRow(ctx, sql, args...)

	var count int
	err := row.Scan(&count); if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return 0, nil
		}

		return count, err
	}

	return count, nil
}

func sumCartItemsSubtotalByStoreID(ctx context.Context, db *pgxpool.Pool, cartID string, storeID string) (int, error) {
	sqlf.SetDialect(sqlf.PostgreSQL)
	query := sqlf.
		From("cart_items").
		Where("cart_id = ?", cartID).
		Where("store_id = ?", storeID).
		Where("quantity > 0").
		GroupBy("store_id").
		Select("COALESCE(SUM(subtotal), 0) as subtotal")
	
	sql, args := query.String(), query.Args()
	row := db.QueryRow(ctx, sql, args...)

	var subtotal int
	err := row.Scan(&subtotal); if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return 0, nil
		}

		return subtotal, err
	}

	return subtotal, nil
}

func deleteAllMenusFromCartByStoreID(ctx context.Context, db *pgxpool.Pool, cartID string, storeID string) error {
	sqlf.SetDialect(sqlf.PostgreSQL)
	query := sqlf.
		DeleteFrom("cart_items").
		Where("cart_id = ?", cartID).
		Where("store_id = ?", storeID)
	
	sql, args := query.String(), query.Args()
	_, err := db.Exec(ctx, sql, args...); if err != nil {
		return err
	}

	return nil
}

func getStoresByCartID(ctx context.Context, db *pgxpool.Pool, cartID string, ) ([]entities.Store, error) {
	sqlf.SetDialect(sqlf.PostgreSQL)
	query := sqlf.
		Select("DISTINCT ON (s.id) s.id, s.slug, s.name").
		From("carts as c").
		Join("cart_items as ci", "ci.cart_id = c.id").
		Join("stores as s", "s.id = ci.store_id").
		Where("c.id = ?", cartID).
		Where("ci.quantity > ?", 0).
		OrderBy("s.id, ci.updated_at DESC")

	sql, args := query.String(), query.Args()
	rows, err := db.Query(ctx, sql, args...); if err != nil {
		return []entities.Store{}, err
	}
	stores, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (entities.Store, error) {
		var store entities.Store
		err = row.Scan(
			&store.ID,
			&store.Slug,
			&store.Name,
		)
		if err != nil {
			return store, err
		}

		return store, err
	})

	if err != nil{
		return []entities.Store{}, err
	}

	return stores, nil
}

func getMenusInCart(ctx context.Context, db *pgxpool.Pool, cartID string) ([]entities.StoreMenuWithQuantityAndSubtotal, error) {
	sqlf.SetDialect(sqlf.PostgreSQL)
	query := sqlf.
		Select("sm.id, sm.store_id, sm.name, sm.image, sm.price, sm.price_promo, ci.quantity, ci.subtotal").
		From("cart_items as ci").
		Clause(`JOIN store_menus as sm ON ci.store_menu_id = sm.id AND sm.is_available = true AND ci.store_id = sm.store_id AND ci.quantity > 0`).
		Where("ci.cart_id = ?", cartID).
		OrderBy("ci.updated_at DESC")

	sql, args := query.String(), query.Args()
	rows, err := db.Query(ctx, sql, args...); if err != nil {
		return []entities.StoreMenuWithQuantityAndSubtotal{}, err
	}
	menus, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (entities.StoreMenuWithQuantityAndSubtotal, error) {
		var menu entities.StoreMenuWithQuantityAndSubtotal
		err = row.Scan(
			&menu.ID,
			&menu.StoreID,
			&menu.Name,
			&menu.Image,
			&menu.Price,
			&menu.PricePromo,
			&menu.Quantity,
			&menu.Subtotal,
		)
		if err != nil {
			return menu, err
		}

		return menu, err
	})

	if err != nil{
		return []entities.StoreMenuWithQuantityAndSubtotal{}, err
	}

	return menus, nil
}