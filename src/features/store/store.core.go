package store

import (
	"github.com/IqbalLx/food-order/src/shared/entities"
)

func populateMatcingMenus(data *[]entities.StoreWithMatchingMenu, menus []entities.StoreMenuWithQuantity) error {
	menusGroupedByStore := make(map[string][]entities.StoreMenuWithQuantity)
	for _, menu := range menus {
		_, isExists := menusGroupedByStore[menu.StoreID]; if !isExists {
			menusGroupedByStore[menu.StoreID] = []entities.StoreMenuWithQuantity{menu}
			continue
		}

		menusGroupedByStore[menu.StoreID] = append(menusGroupedByStore[menu.StoreID], menu)
	}

	for i, store := range (*data) {
		if menus, ok := menusGroupedByStore[store.ID]; ok {
			(*data)[i].Menus = menus
		}
	}

	return nil
}