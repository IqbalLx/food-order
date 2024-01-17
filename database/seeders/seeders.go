package main

import (
	"context"
	"log"

	"github.com/IqbalLx/food-order/src/shared/utils"
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

	seedStores(ctx, db)
	seedCategories(ctx, db)
	seedStoreCategories(ctx, db)
	seedStoreMenus(ctx, db)
	seedStoreMenuCategory(ctx, db)
}