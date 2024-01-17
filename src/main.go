package main

import (
	"context"
	"log"

	"github.com/IqbalLx/food-order/src/features/home"
	"github.com/IqbalLx/food-order/src/features/store"
	"github.com/IqbalLx/food-order/src/shared/utils"
	"github.com/gofiber/fiber/v2"

	"github.com/jackc/pgx/v5"
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
	db, err := pgx.Connect(ctx, connOrDSNString)
	if err != nil {
		log.Fatalln("open: %w", err)
	}
	defer db.Close(ctx)

	// initialize app
	appConfig := utils.NewAppConfig(env, pgConfig)
    
	app := fiber.New()

	app.Static("/static", "./src/static") 
	app.Static("/shoelace", "./node_modules/@shoelace-style/shoelace") 

	app.Use(func(c *fiber.Ctx) error {
		utils.SetLocal[*pgx.Conn](c, "db", db)
		return c.Next()
	})

	home.NewHomeController(app)
	store.NewStoreController(app)

    app.Listen(appConfig.Address)
}