package entities

type Category struct {
	ID string
	Name string
}

type StoreCategory struct {
	ID string
	StoreID string
	CategoryID string
}

type CategoryWithTimestampField struct {
	Category
	TimestampField
}