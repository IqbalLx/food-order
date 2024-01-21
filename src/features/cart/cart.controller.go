package cart

import (
	"strconv"

	components "github.com/IqbalLx/food-order/src/features/cart/views/components"
	layouts "github.com/IqbalLx/food-order/src/features/cart/views/layouts"
	"github.com/IqbalLx/food-order/src/shared/middlewares"
	"github.com/IqbalLx/food-order/src/shared/utils"
	sharedComponents "github.com/IqbalLx/food-order/src/shared/views/components"
	sharedLayouts "github.com/IqbalLx/food-order/src/shared/views/layouts"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewCartController(app *fiber.App) {
	carts := app.Group("/carts", middlewares.ValidateCart)
	carts.Get("/", getCartViewHandler)
	carts.Put("/", upsertCartHandler)
	carts.Get("/items", getCartItemsHandler)
	carts.Get("/count", getCartCountHandler)

	carts.Delete("/stores/:store_id", deleteStoreFromCartHandler)
	carts.Delete("/stores/:store_id/menus/:menu_id", deleteMenuFromCartHandler)

	states := carts.Group("/states")
	states.Get("/", getCartStateHandler)
	states.Get("/stores/:store_id", getStoreCartStateHandler)
}

func getCartViewHandler(c *fiber.Ctx) error {
	return adaptor.HTTPHandler(
		templ.Handler(
			sharedLayouts.RootWithTitle(
				layouts.Cart(),
				"Keranjang",
			),
		),
	)(c)
}

func getCartItemsHandler(c *fiber.Ctx) error {
	appConfig := utils.GetLocal[*utils.AppConfig](c, "appConfig")
	cartID := c.Cookies(appConfig.Name + "__cart")

	db := utils.GetLocal[*pgxpool.Pool](c, "db")

	stores, countItems, totalItems, err := doGetCart(c.Context(), db, cartID); if err != nil {
		return err
	}

	return adaptor.HTTPHandler(
		templ.Handler(
			layouts.CartItems(stores, countItems, totalItems),
		),
	)(c)
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

	if (triggerOrigin == "cart") {
		c.Set("HX-Trigger", "cart-state-update")
		return adaptor.HTTPHandler(
			templ.Handler(components.CartMenuCounter(menuID, finalQuantity, "cart")),
		)(c)
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

func getCartStateHandler(c *fiber.Ctx) error {
	appConfig := utils.GetLocal[*utils.AppConfig](c, "appConfig")
	cartID := c.Cookies(appConfig.Name + "__cart")

	db := utils.GetLocal[*pgxpool.Pool](c, "db")
	state, err := doGetCartState(c.Context(), db, cartID); if err != nil {
		return err
	}

	if state.CountMenus == 0 {
		return adaptor.HTTPHandler(
			templ.Handler(layouts.CartFooterEmpty()),
		)(c)
	}

	return adaptor.HTTPHandler(
		templ.Handler(layouts.CartFooter(state.CountStores, state.CountMenus, state.TotalItems)),
	)(c)
}

func deleteStoreFromCartHandler(c *fiber.Ctx) error {
	storeID := c.Params("store_id")

	appConfig := utils.GetLocal[*utils.AppConfig](c, "appConfig")
	cartID := c.Cookies(appConfig.Name + "__cart")

	db := utils.GetLocal[*pgxpool.Pool](c, "db")

	err := doDeleteAllMenusFromCartByStoreID(c.Context(), db, cartID, storeID); if err != nil {
		return err
	}

	c.Set("HX-Trigger", "cart-state-update")

	// empty string indicates the item removed from view
	return c.Status(200).SendString("")
}

func deleteMenuFromCartHandler(c *fiber.Ctx) error {
	storeID := c.Params("store_id")
	menuID := c.Params("menu_id")

	appConfig := utils.GetLocal[*utils.AppConfig](c, "appConfig")
	cartID := c.Cookies(appConfig.Name + "__cart")

	db := utils.GetLocal[*pgxpool.Pool](c, "db")

	storeQuantityAfterDelete, err := doDeleteMenuFromCart(c.Context(), db, cartID, storeID, menuID); if err != nil {
		return err
	}
	
	triggerEvent := "cart-state-update"
	if (storeQuantityAfterDelete == 0) {
		triggerEvent += ", cart-refresh"
	}

	c.Set("HX-Trigger", triggerEvent)

	// empty string indicates the item removed from view
	return c.Status(200).SendString("")
}

func getStoreCartStateHandler(c *fiber.Ctx) error {
	storeID := c.Params("store_id")

	appConfig := utils.GetLocal[*utils.AppConfig](c, "appConfig")
	cartID := c.Cookies(appConfig.Name + "__cart")

	db := utils.GetLocal[*pgxpool.Pool](c, "db")

	countItems, subtotalItems, err := doGetCartStateByStoreID(c.Context(), db, cartID, storeID); if err != nil {
		return err
	}

	return adaptor.HTTPHandler(
		templ.Handler(
			components.CartItemStateInfo(countItems, subtotalItems),
		),
	)(c)
}