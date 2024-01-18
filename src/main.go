package main

import (
	"context"
	"log"

	"github.com/IqbalLx/food-order/src/features/cart"
	"github.com/IqbalLx/food-order/src/features/home"
	"github.com/IqbalLx/food-order/src/features/store"
	"github.com/IqbalLx/food-order/src/shared/middlewares"
	"github.com/IqbalLx/food-order/src/shared/utils"
	"github.com/gofiber/fiber/v2"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	env := utils.NewEnv()

	// initialize database
	pgConfig := utils.NewPostgresConfig(env)
	var connOrDSNString string
	if (pgConfig.ConnStringAvailable) {
		connOrDSNString = pgConfig.ConnString
	} else {
		connOrDSNString = pgConfig.DSNString()
	}

	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, connOrDSNString)     // pgx.Connect(ctx, connOrDSNString)
	if err != nil {
		log.Fatalln("open: %w", err)
	}
	defer dbpool.Close()

	// initialize app
	appConfig := utils.NewAppConfig(env, pgConfig)
    
	app := fiber.New()

	app.Static("/static", "./src/static") 
	app.Static("/shoelace", "./node_modules/@shoelace-style/shoelace") 

	app.Use(func(c *fiber.Ctx) error {
		utils.SetLocal[*pgxpool.Pool](c, "db", dbpool)
		utils.SetLocal[*utils.AppConfig](c, "appConfig", appConfig)

		return c.Next()
	})
	app.Use(middlewares.CreateCartForNewUser)

	home.NewHomeController(app)
	store.NewStoreController(app)
	cart.NewCartController(app)

    app.Listen(appConfig.Address)
}