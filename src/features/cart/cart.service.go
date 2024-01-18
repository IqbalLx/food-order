package cart

import (
	"context"
	"errors"

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
		err = deleteMenuFromCart(ctx, db, cartID, menu)
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