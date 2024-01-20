package store

import (
	"fmt"
	"net/http"
	"strconv"

	components "github.com/IqbalLx/food-order/src/features/store/views/components"
	layouts "github.com/IqbalLx/food-order/src/features/store/views/layouts"
	"github.com/IqbalLx/food-order/src/shared/middlewares"
	"github.com/IqbalLx/food-order/src/shared/utils"
	sharedComponents "github.com/IqbalLx/food-order/src/shared/views/components"
	sharedLayouts "github.com/IqbalLx/food-order/src/shared/views/layouts"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewStoreController(app *fiber.App) {
	stores := app.Group("/stores")
	stores.Get("/", getStoresHandler)

	stores.Get("/search", middlewares.ValidateCard, searchStoresByMenuNameHandler)
	stores.Post("/search", middlewares.ValidateCard, searchStoresByMenuNameHandler)

	stores.Get("/:slug", middlewares.ValidateCard, getStoreBySlugHandler)
	stores.Get("/:id/menus", middlewares.ValidateCard, getStoreMenusByStoreIDHandler)
	
	stores.Get("/:id/menus", middlewares.ValidateCard, getStoreMenusByStoreIDHandler) 

	// POST used to get FormData required for searching
	stores.Post("/:id/menus", middlewares.ValidateCard, getStoreMenusByStoreIDHandler)
	stores.Post("/:id/menus/categories", getMenuCategoriesHandler)

	// states
	storeStates := stores.Group("/states")
	storeStates.Get("/checkout", middlewares.ValidateCard, getCheckoutStatehandler)
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

	db := utils.GetLocal[*pgxpool.Pool](c, "db")

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

	appConfig := utils.GetLocal[*utils.AppConfig](c, "appConfig")
	cartID := c.Cookies(appConfig.Name + "__cart")

	db := utils.GetLocal[*pgxpool.Pool](c, "db")

	searchQuery := c.Query("query", "")
	isWithSearchQuery := searchQuery != ""

	data, err := doGetInitialStoreDetail(c.Context(), db, cartID, storeSlug, initialMenuSize, isWithSearchQuery, searchQuery); if err != nil {
		return err
	}

	return adaptor.HTTPHandler(
		templ.Handler(
			sharedLayouts.RootWithTitle(
				layouts.Store(
					data.Store, data.MenuCategories, data.Menus, 
					initialMenuSize, data.IsMenusScrollable,
					isWithSearchQuery, searchQuery,
					), 
				data.Store.Name,
			),
		),
	)(c)
}

func getStoreMenusByStoreIDHandler(c *fiber.Ctx) error {
	// if not from HTMX redirect to home
	fromHTMX := c.Get("HX-Request", "false")
	if (fromHTMX == "false") {
		return c.Redirect("/")
	}

	storeID := c.Params("id")
	size, err := strconv.Atoi(c.Query("size", "5")); if err != nil {
		return err
	}
	lastMenuSecondaryID, err := strconv.Atoi(c.Query("last_menu_secondary_id", "0")); if err != nil {
		return err
	}

	searchQuery := ""
	isWithSearchQuery := false

	if c.Route().Method == "POST" {
		searchQuery = c.FormValue("query", "")
		isWithSearchQuery = searchQuery != ""

		c.Set("HX-Trigger", "store-menu-category-update")
	} else if c.Route().Method == "GET" {
		searchQuery = c.Query("query", "")
		isWithSearchQuery = searchQuery != ""
	}

	menuCategoryID := c.Query("menu_category_id", "")
	isWithCategory := menuCategoryID != ""

	appConfig := utils.GetLocal[*utils.AppConfig](c, "appConfig")
	cartID := c.Cookies(appConfig.Name + "__cart")

	db := utils.GetLocal[*pgxpool.Pool](c, "db")

	data, err := doGetMenus(c.Context(), db, cartID, storeID, size, lastMenuSecondaryID, isWithCategory, menuCategoryID, 
	isWithSearchQuery, searchQuery); if err != nil {
		return err
	}

	if (isWithSearchQuery && c.Route().Method == "POST") {
		c.Set("HX-Push-Url", fmt.Sprintf("/stores/%s?query=%s", data.Store.Slug, utils.EncodeQuerystring(searchQuery)))
	} else if (!isWithSearchQuery && c.Route().Method == "POST") {
		c.Set("HX-Push-Url", fmt.Sprintf("/stores/%s", data.Store.Slug))
	}
	return adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			for i, menu := range data.Menus {
				components.MenuCard(data.Store, menu, size, i == len(data.Menus) - 1, data.IsMenusScrollable, isWithCategory, 
				menuCategoryID, isWithSearchQuery, searchQuery).Render(c.Context(), w)
			}

			if data.IsMenusScrollable {
				sharedComponents.GenericCardSkeleton("menu-last-card").Render(c.Context(), w)
			}
		},
	)(c)
}

func getMenuCategoriesHandler(c *fiber.Ctx) error {
	// if not from HTMX redirect to home
	fromHTMX := c.Get("HX-Request", "false")
	if (fromHTMX == "false") {
		return c.Redirect("/")
	}

	storeID := c.Params("id")
	searchQuery := c.FormValue("query", "")
	isWithSearchQuery := searchQuery != ""

	db := utils.GetLocal[*pgxpool.Pool](c, "db")

	data, err := doGetMenuCategories(c.Context(), db, storeID, isWithSearchQuery, searchQuery); if err != nil {
		return err
	}

	return adaptor.HTTPHandler(
		templ.Handler(
			layouts.StoreFooter(data.Store, data.MenuCategories),
		),
	)(c)
}

func searchStoresByMenuNameHandler(c *fiber.Ctx) error {
	var searchQuery string
	var page, pageSize int

	if c.Route().Method == "POST" {
		searchQuery = c.FormValue("query", "")

		// POST request indicates initial search action, so set default page and size
		page = 1
		pageSize = 10
	} else if c.Route().Method == "GET" {
		searchQuery = c.Query("query", "")

		pageStr := c.Query("page", "1")
		pageSizeStr := c.Query("size", "10")

		var err error
		page, err = strconv.Atoi(pageStr); if err != nil {
			return err
		}
		pageSize, err = strconv.Atoi(pageSizeStr); if err != nil {
			return err
		}
	}

	if (searchQuery == "") {
		return adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for i := 0; i < 2; i++ {
				sharedComponents.StoreCardWithMenuIndicator().Render(c.Context(), w)
			}
		})(c)
	}

	appConfig := utils.GetLocal[*utils.AppConfig](c, "appConfig")
	cartID := c.Cookies(appConfig.Name + "__cart")

	db := utils.GetLocal[*pgxpool.Pool](c, "db")
	stores, maxPage, err := doSearchStoresByMenuName(c.Context(), db, cartID, searchQuery, page, pageSize); if err != nil {
		return err
	}

	nextPage := page + 1
	isNextPageAvailable := nextPage <= maxPage

	c.Set("HX-Push-Url", fmt.Sprintf("/search?query=%s", utils.EncodeQuerystring(searchQuery)))
	return adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for i, store := range stores {
			sharedComponents.StoreCardWithMenu(
				store, 
				searchQuery, 
				pageSize, 
				isNextPageAvailable, 
				nextPage, 
				i == len(stores) - 1,
			 ).Render(c.Context(), w)
		}

		if isNextPageAvailable {
			sharedComponents.GenericCardSkeleton("store-with-menu-last-card").Render(c.Context(), w)
		}

	})(c)
}

// Store States
func getCheckoutStatehandler(c *fiber.Ctx) error {
	appConfig := utils.GetLocal[*utils.AppConfig](c, "appConfig")
	cartID := c.Cookies(appConfig.Name + "__cart")

	db := utils.GetLocal[*pgxpool.Pool](c, "db")
	readyToCheckout, err := doCheckIfReadyToCheckout(c.Context(), db, cartID); if err != nil {
		return err
	}

	if !readyToCheckout {
		return c.Status(200).SendString("")
	}

	return adaptor.HTTPHandler(
		templ.Handler(components.CheckoutButton()),
	)(c)
}