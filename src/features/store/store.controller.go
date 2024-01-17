package store

import (
	"net/http"
	"strconv"

	components "github.com/IqbalLx/food-order/src/features/store/views/components"
	layouts "github.com/IqbalLx/food-order/src/features/store/views/layouts"
	"github.com/IqbalLx/food-order/src/shared/utils"
	sharedComponents "github.com/IqbalLx/food-order/src/shared/views/components"
	sharedLayouts "github.com/IqbalLx/food-order/src/shared/views/layouts"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/jackc/pgx/v5"
)

func NewStoreController(app *fiber.App) {
	stores := app.Group("/stores")
	stores.Get("/", getStoresHandler)
	stores.Get("/:slug", getStoreBySlugHandler)
	stores.Get("/:id/menus", getStoreMenusByStoreIDHandler)
}

func getStoresHandler(c *fiber.Ctx) error {
	// if not from HTMX redirect to home
	fromHTMX := c.Get("HX-Request", "false")
	if (fromHTMX == "false") {
		return c.Redirect("/")
	}

	size, err := strconv.Atoi(c.Query("size", "5")); if err != nil {
		return err
	}
	lastStoreSecondaryID, err := strconv.Atoi(c.Query("last_store_secondary_id", "0")); if err != nil {
		return err
	}

	db := utils.GetLocal[*pgx.Conn](c, "db")

	stores, isScrollable, err := doGetStores(c.Context(), db, size, lastStoreSecondaryID); if err != nil {
		return err
	}

	return adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			for i, store := range stores {
				sharedComponents.StoreCard(store, size, i == len(stores) - 1, isScrollable).Render(c.Context(), w)
			}

			if isScrollable {
				sharedComponents.GenericCardSkeleton("store-last-card").Render(c.Context(), w)
			}
		},
	)(c)
}

func getStoreBySlugHandler(c *fiber.Ctx) error {
	storeSlug := c.Params("slug")

	initialMenuSize := 5
	
	db := utils.GetLocal[*pgx.Conn](c, "db")
	data, err := doGetInitialStoreDetail(c.Context(), db, storeSlug, initialMenuSize); if err != nil {
		return err
	}

	return adaptor.HTTPHandler(
		templ.Handler(
			sharedLayouts.RootWithTitle(
				layouts.Store(data.Store, data.MenuCategories, data.Menus, initialMenuSize, data.IsMenusScrollable), 
				data.Store.Name,
			),
		),
	)(c)
}

func getStoreMenusByStoreIDHandler(c *fiber.Ctx) error {
	storeID := c.Params("id")
	size, err := strconv.Atoi(c.Query("size", "5")); if err != nil {
		return err
	}
	lastMenuSecondaryID, err := strconv.Atoi(c.Query("last_menu_secondary_id", "0")); if err != nil {
		return err
	}

	menuCategoryID := c.Query("menu_category_id", "")
	isWithCategory := menuCategoryID != ""

	db := utils.GetLocal[*pgx.Conn](c, "db")

	data, err := doGetMenus(c.Context(), db, storeID, size, lastMenuSecondaryID, isWithCategory, menuCategoryID); if err != nil {
		return err
	}

	return adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			for i, menu := range data.Menus {
				components.MenuCard(data.Store, menu, size, i == len(data.Menus) - 1, data.IsMenusScrollable, isWithCategory, menuCategoryID).Render(c.Context(), w)
			}

			if data.IsMenusScrollable {
				sharedComponents.GenericCardSkeleton("menu-last-card").Render(c.Context(), w)
			}
		},
	)(c)
}