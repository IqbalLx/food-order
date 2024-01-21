package cart

import (
	"github.com/IqbalLx/food-order/src/shared/entities"
)

func populateCartMenus(data *[]entities.StoreWithCartMenus, menus []entities.StoreMenuWithQuantityAndSubtotal) error {
	menusGroupedByStore := make(map[string][]entities.StoreMenuWithQuantityAndSubtotal)
	menusSubtotalByStore := make(map[string]int)
	
	for _, menu := range menus {
		_, isExists := menusGroupedByStore[menu.StoreID]; if !isExists {
			menusGroupedByStore[menu.StoreID] = []entities.StoreMenuWithQuantityAndSubtotal{menu}
			menusSubtotalByStore[menu.StoreID] = menu.Subtotal
			continue
		}

		menusGroupedByStore[menu.StoreID] = append(menusGroupedByStore[menu.StoreID], menu)
		menusSubtotalByStore[menu.StoreID] += menu.Subtotal
	}

	for i, store := range (*data) {
		if menus, ok := menusGroupedByStore[store.ID]; ok {
			(*data)[i].Menus = menus
		}

		if subtotal, ok := menusSubtotalByStore[store.ID]; ok {
			(*data)[i].Subtotal = subtotal
		}
	}

	return nil
}