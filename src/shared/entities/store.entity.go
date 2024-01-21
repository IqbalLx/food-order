package entities

type Store struct {
	ID string
	SecondaryID int
	Name string
	Slug string
	Image string
	ShortDesc string
	Desc string
	Rating int
}

type StoreWithTimestampField struct {
	Store
	TimestampField
}

type StoreWithCategories struct {
	Store
	Categories []string
}

type StoreWithMenus struct {
	Store
	Menus []StoreMenu
}

type StoreWithMatchingMenu struct {
	StoreWithCategories
	MatchingMenuCount int
	Menus []StoreMenuWithQuantity
}

type StoreWithCartMenus struct {
	Store
	Subtotal int
	Menus []StoreMenuWithQuantityAndSubtotal
}