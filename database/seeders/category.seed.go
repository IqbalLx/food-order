package main

import (
	"context"
	"fmt"

	"github.com/IqbalLx/food-order/src/shared/entities"
	"github.com/jackc/pgx/v5"
	"github.com/leporo/sqlf"
)

func seedCategories(ctx context.Context, db *pgx.Conn) {
	categories := []entities.Category{
		{ID: "d969b527-24d7-4a9f-9a77-b66df38ed4cf", Name: "Appetizers"},
		{ID: "61a0fc9d-511a-4cd3-83e5-961ec3d18622", Name: "Main Courses"},
		{ID: "0c3fe8ac-ccf9-4d02-830a-76e7d8c77d5a", Name: "Desserts"},
		{ID: "7c1e17c5-cc36-43d4-a02a-6ad7ac660b12", Name: "Beverages"},
		{ID: "af16cbbd-c21d-4092-bbd9-1b67e587dd7e", Name: "Salads"},
		{ID: "e9c80b42-3094-4c87-b327-9a749e2aa3f5", Name: "Snacks"},
		{ID: "d0e8a59a-2498-4825-912a-dcb3fc12f3f2", Name: "Breakfast"},
		{ID: "36e9e128-0cfc-4b69-bc4a-94c953b1eb8e", Name: "Drinks"},
		{ID: "2a9cd8f8-1f35-48e7-ba6c-bdd0d48b966c", Name: "Soups"},
		{ID: "2f49954c-0dab-4a36-89a5-748b49de8312", Name: "Dairy Products"},
	}
	
	sqlf.SetDialect(sqlf.PostgreSQL)
	query := sqlf.InsertInto("categories")

	for _, category := range categories {
		query.NewRow().
			Set("id", category.ID).
			Set("name", category.Name)
	}

	fmt.Println("Seeding table categories...")
	sql, args := query.String(), query.Args()
	if _, err := db.Exec(context.Background(), sql, args...); err != nil {
		fmt.Println(err.Error())
		fmt.Println("Seeding table categories failed")
	} else {
		fmt.Println("Seeding table categories done!")
	}
}