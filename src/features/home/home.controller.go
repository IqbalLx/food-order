package home

import (
	layouts "github.com/IqbalLx/food-order/src/features/home/views/layouts"
	"github.com/IqbalLx/food-order/src/shared/utils"
	sharedLayouts "github.com/IqbalLx/food-order/src/shared/views/layouts"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewHomeController(app *fiber.App) {
	app.Get("/", homeHandler)
}

func homeHandler(c *fiber.Ctx) error {
	db := utils.GetLocal[*pgxpool.Pool](c, "db")

	initialSize := 5
	data, err := doGetHomeInitialData(c.Context(), db, initialSize); if err != nil {
		return err
	}

	return adaptor.HTTPHandler(
		templ.Handler(
			sharedLayouts.Root(
				layouts.Home(
					data.Stores, 
					initialSize, 
					len(data.Stores), 
					data.IsStoresScrollable,
				),
			),
		),
	)(c)

}