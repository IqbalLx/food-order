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

type StoreMenuWithQuantity struct {
	StoreMenu
	Quantity int
}

type StoreMenuWithQuantityAndSubtotal struct {
	StoreMenu
	Quantity int
	Subtotal int
}

type StoreMenuWithTimestamp struct {
	StoreMenu
	TimestampField
}