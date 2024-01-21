package entities

type Cart struct {
	ID string
	UserID string
}

type CartItem struct {
	ID string
	CartID string
	StoreID string
	StoreMenuID string
	Quantity int
	Subtotal int
}

type CartWithTimestamp struct {
	Cart
	TimestampField
}

type CartItemWithTimestamp struct {
	CartItem
	TimestampField
}