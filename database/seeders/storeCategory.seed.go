package main

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/IqbalLx/food-order/src/shared/entities"
	"github.com/jackc/pgx/v5"
	"github.com/leporo/sqlf"
)

func seedStoreCategories(ctx context.Context, db *pgx.Conn) {
	stores := [...]entities.Store{
		{ID: "a6d86a41-58a4-4b69-9b72-2b8cf5c8f9f1", Name: "Chinese Restaurant", Slug: "chinese-restaurant", ShortDesc: "Authentic Chinese Cuisine", Desc: "A restaurant serving delicious Chinese dishes.", Rating: 0},
		{ID: "d2a5d3d7-8b6f-4e3c-a1b1-9e7808963b25", Name: "Pizza Place", Slug: "pizza-place", ShortDesc: "Best Pizza in Town", Desc: "Enjoy a variety of mouth-watering pizzas at our place.", Rating: 45},
		{ID: "7e25e8a0-16d9-42a8-a64a-7c8d1c8b991c", Name: "Coffee Shop", Slug: "coffee-shop", ShortDesc: "Premium Coffee Selection", Desc: "Indulge in the finest coffee blends at our cozy cafe.", Rating: 35},
		{ID: "8b9c7a2f-9263-4b53-81c9-40b953b4e918", Name: "Bakery Delight", Slug: "bakery-delight", ShortDesc: "Freshly Baked Goods", Desc: "Satisfy your sweet tooth with our delicious baked treats.", Rating: 40},
		{ID: "3f5f9a24-16ab-4e22-b4d0-dab51f09ee0c", Name: "Sushi Haven", Slug: "sushi-haven", ShortDesc: "Exquisite Japanese Sushi", Desc: "Experience the art of sushi-making at our authentic Japanese restaurant.", Rating: 45},
		{ID: "ea5e2bc0-9d8a-4eb7-9f91-3c13a7dbbc94", Name: "Taco Time", Slug: "taco-time", ShortDesc: "Tantalizing Tacos", Desc: "Delight in the flavors of Mexico with our mouthwatering tacos.", Rating: 40},
		{ID: "32a0b1fe-4e9b-4d9b-9e84-cc6b583e81da", Name: "Smoothie Paradise", Slug: "smoothie-paradise", ShortDesc: "Fresh and Healthy Smoothies", Desc: "Quench your thirst with our refreshing and nutritious smoothies.", Rating: 35},
		{ID: "de7445c4-3c10-4cb2-a55b-7753f60e2aa6", Name: "Burger Joint", Slug: "burger-joint", ShortDesc: "Gourmet Burgers", Desc: "Savor the juiciest gourmet burgers in town at our joint.", Rating: 40},
		{ID: "b1048131-670a-4183-93ea-3f9330ea6e33", Name: "Mediterranean Delight", Slug: "mediterranean-delight", ShortDesc: "Authentic Mediterranean Cuisine", Desc: "Transport your taste buds to the Mediterranean with our flavorful dishes.", Rating: 45},
		{ID: "864c55bd-2da7-4f7a-a5c9-88944f6ee774", Name: "Ice Cream Oasis", Slug: "ice-cream-oasis", ShortDesc: "Irresistible Ice Creams", Desc: "Treat yourself to a delightful selection of premium ice creams at our oasis.", Rating: 40},
		{ID: "c4d2e8b1-0e63-4a8f-bc3a-ff9a4f63186d", Name: "Cozy Corner Cafe", Slug: "cozy-corner-cafe", ShortDesc: "A charming spot for coffee lovers", Desc: "Discover the warmth of our cozy cafe, serving the finest coffee in town.", Rating: 35},
		{ID: "eaa6c8f9-23d1-42f5-b8a7-7b5fb3f7ee91", Name: "The Spice House", Slug: "the-spice-house", ShortDesc: "Spice up your dining experience", Desc: "Embark on a culinary journey with our diverse menu of flavorful dishes.", Rating: 40},
		{ID: "b8a0d9f7-13ac-4f79-b4a6-ecf35f8d6902", Name: "Sunset Bistro", Slug: "sunset-bistro", ShortDesc: "Elegant dining with a view", Desc: "Experience exquisite dishes while enjoying a breathtaking sunset view at our bistro.", Rating: 45},
		{ID: "8f9a4b62-6e14-45e5-b7bd-ae5e7b1d2f88", Name: "Green Garden Grill", Slug: "green-garden-grill", ShortDesc: "Healthy and Delicious", Desc: "Indulge in a delightful array of healthy and delicious dishes at our grill.", Rating: 40},
		{ID: "d16a4f79-1b2d-4a7c-8e1d-3f5f9a2dc9e5", Name: "Blue Seafood Paradise", Slug: "blue-seafood-paradise", ShortDesc: "Ocean-inspired culinary delights", Desc: "Savor the freshness of the ocean with our delectable seafood offerings.", Rating: 45},
		{ID: "c9e584e2-b8f0-4a97-968c-2e03a3c17bf7", Name: "Golden Wok Express", Slug: "golden-wok-express", ShortDesc: "Quick and Tasty Asian Cuisine", Desc: "Enjoy the convenience of quick and tasty Asian cuisine at our express outlet.", Rating: 35},
		{ID: "a7b3c9e5-f9a2-43d6-b8c4-3e5d7f2c48a0", Name: "Sunny Side Breakfast", Slug: "sunny-side-breakfast", ShortDesc: "Start your day with a smile", Desc: "Delight in a sunny breakfast experience with our mouthwatering morning menu.", Rating: 40},
		{ID: "9d8a4b6f-7e25-4c1c-b2d3-1f09ee0c7e5a", Name: "Royal Tea Palace", Slug: "royal-tea-palace", ShortDesc: "Elegant Tea Selection", Desc: "Immerse yourself in the royalty of tea with our exquisite selection and ambiance.", Rating: 45},
	}

	categories := [...]entities.Category{
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
	query := sqlf.InsertInto("store_categories")

	for _, store := range stores {
		numCategories := 1 + rand.Intn(3)
		var prevCategoryID string
		for i := 0; i < numCategories; i++ {

			// make sure category unique for each store by loop rand until it found unique value
			var category entities.Category
			for {
				category = categories[rand.Intn(len(categories))]
				if (category.ID != prevCategoryID) { 
					prevCategoryID = category.ID
					break 
				}

				prevCategoryID = category.ID
			}

			query.
				NewRow().
					Set("store_id", store.ID).
					Set("category_id", category.ID)
		}
	}

	fmt.Println("Seeding table store_categories...")
	sql, args := query.String(), query.Args()

	if _, err := db.Exec(context.Background(), sql, args...); err != nil {
		fmt.Println(err.Error())
		fmt.Println("Seeding table store_categories failed")
	} else {
		fmt.Println("Seeding table store_categories done!")
	}
}