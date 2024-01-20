package cart

import (
	"strconv"

	"github.com/IqbalLx/food-order/src/shared/middlewares"
	"github.com/IqbalLx/food-order/src/shared/utils"
	sharedComponents "github.com/IqbalLx/food-order/src/shared/views/components"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewCartController(app *fiber.App) {
	carts := app.Group("/carts", middlewares.ValidateCart)
	carts.Get("/count", getCartCountHandler)
	carts.Put("/", upsertCartHandler)
}

func getCartCountHandler(c *fiber.Ctx) error {
	// if not from HTMX reject
	fromHTMX := c.Get("HX-Request", "false")
	if (fromHTMX == "false") {
		return &fiber.Error{Code: 400}
	}

	appConfig := utils.GetLocal[*utils.AppConfig](c, "appConfig")
	cartID := c.Cookies(appConfig.Name + "__cart")

	db := utils.GetLocal[*pgxpool.Pool](c, "db")
	newCount, err := doCountCartItems(c.Context(), db, cartID); if err != nil {
		return err
	}

	return adaptor.HTTPHandler(
		templ.Handler(
			sharedComponents.CartIcon(
				newCount,
			),
		),
	)(c) 
}

func upsertCartHandler(c *fiber.Ctx) error {
	// if not from HTMX reject
	fromHTMX := c.Get("HX-Request", "false")
	if (fromHTMX == "false") {
		return &fiber.Error{Code: 400}
	}

	appConfig := utils.GetLocal[*utils.AppConfig](c, "appConfig")
	cartID := c.Cookies(appConfig.Name + "__cart")

	menuID := c.FormValue("menu_id")
	triggerOrigin := c.FormValue("origin")
	quantity, err := strconv.Atoi(c.FormValue("quantity", "1")); if err != nil {
		return err
	}

	db := utils.GetLocal[*pgxpool.Pool](c, "db")
	finalQuantity, err := doUpsertMenuToCart(c.Context(), db, cartID, menuID, quantity); if err != nil {
		return err
	}

	c.Set("HX-Trigger", "cart-count-update")
	if finalQuantity == 0 {
		return adaptor.HTTPHandler(
			templ.Handler(
				sharedComponents.MenuInitialPlusButton(
					menuID,
					triggerOrigin,
				),
			),
		)(c) 
	}

	return adaptor.HTTPHandler(
		templ.Handler(
			sharedComponents.MenuCounter(
				menuID,
				finalQuantity,
				triggerOrigin,
			),
		),
	)(c) 
}