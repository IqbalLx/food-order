package entities

type StoreMenuCategory struct {
	ID string
	StoreID string
	Name string
}

type StoreMenuCategoryItems struct {
	ID string
	StoreID string
	StoreMenuCategoryID string
	StoreMenuID string
}

type StoreMenuCategoryWithTimestamp struct {
	StoreMenuCategory
	TimestampField
}