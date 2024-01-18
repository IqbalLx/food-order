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
			Set("quantity", quantity)

	if menu.PricePromo > 0 {
		query.Set("subtotal", quantity * menu.PricePromo)
	} else {
		query.Set("subtotal", quantity * menu.Price)
	}

	query.
		Clause("ON CONFLICT ON CONSTRAINT cart_items_unique_constraint DO UPDATE SET").
			Expr("quantity = EXCLUDED.quantity").
			Expr("subtotal = EXCLUDED.subtotal").
	Returning("quantity")

	sql, args := query.String(), query.Args()
	row := db.QueryRow(ctx, sql, args...)

	var finalQuantity int
	err := row.Scan(&finalQuantity); if err != nil {
		return finalQuantity, err
	}

	return finalQuantity, nil
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