package entities

type StoreMenu struct {
	ID string
	SecondaryID int
	StoreID string
	Name string
	Image string
	Price int
	OrderedCount int
	PricePromo int
	IsAvailable bool
}

type StoreMenuWithTimestamp struct {
	StoreMenu
	TimestampField
}