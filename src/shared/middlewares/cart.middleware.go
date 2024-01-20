package middlewares

import (
	"context"

	"github.com/IqbalLx/food-order/src/shared/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leporo/sqlf"
)

// data
func createNewCart(ctx context.Context, db *pgxpool.Pool) (string, error) {
	cartID := uuid.New().String()

	sqlf.SetDialect(sqlf.PostgreSQL)
	query := sqlf.InsertInto("carts").NewRow().Set("id", cartID)
	sql, args := query.String(), query.Args()

	_, err := db.Exec(ctx, sql, args...); if err != nil {
		return cartID, err
	}

	return cartID, nil
}

func isCartExists(ctx context.Context, db *pgxpool.Pool, cartID string) (bool, error) {
	sqlf.SetDialect(sqlf.PostgreSQL)
	query := sqlf.
		From("carts").
		Select("1 as one").
		Where("id = ?", cartID).
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

func CreateCartForNewUser(c *fiber.Ctx) error {
	// if from HTMX skip cookie setting, only when user load page directly
	fromHTMX := c.Get("HX-Request", "false")
	if (fromHTMX == "true") {
		return c.Next()
	}

	appConfig := utils.GetLocal[*utils.AppConfig](c, "appConfig")
	cartCookieName := appConfig.Name + "__cart"
	
	cartID := c.Cookies(cartCookieName, "")
	if cartID != "" {
		return c.Next()
	}

	// if not set
	db := utils.GetLocal[*pgxpool.Pool](c, "db")
	newCartID, err := createNewCart(c.Context(), db); if err != nil {
		return err
	}

	newCartCookie := &fiber.Cookie{
		Name: cartCookieName,
		Value: newCartID,
	}
	c.Cookie(newCartCookie)

	return c.Next()
}

func ValidateCart(c *fiber.Ctx) error {
	appConfig := utils.GetLocal[*utils.AppConfig](c, "appConfig")
	cartCookieName := appConfig.Name + "__cart"
	
	cartID := c.Cookies(cartCookieName, "")
	if cartID == "" {
		return c.Redirect("/")
	}

	db := utils.GetLocal[*pgxpool.Pool](c, "db")
	isCartExists, err := isCartExists(c.Context(), db, cartID); if err != nil {
		return err
	}

	if !isCartExists {
		return c.Redirect("/")
	}

	return c.Next()
}