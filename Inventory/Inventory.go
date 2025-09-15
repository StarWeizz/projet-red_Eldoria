package inventory

import (
	items "eldoria/items"
	"fmt"
)

type Inventory struct {
	Items map[string]int
	Refs  map[string]items.Item
}

func NewInventory() *Inventory {
	return &Inventory{
		Items: make(map[string]int),
		Refs:  make(map[string]items.Item),
	}
}

func (inv *Inventory) Add(item items.Item, qty int) {
	name := item.GetName()
	inv.Items[name] += qty
	inv.Refs[name] = item
}

func (inv *Inventory) Remove(item items.Item, qty int) bool {
	name := item.GetName()
	if inv.Items[name] < qty {
		return false
	}
	inv.Items[name] -= qty
	if inv.Items[name] == 0 {
		delete(inv.Items, name)
		delete(inv.Refs, name)
	}
	return true
}

func (inv *Inventory) List() {
	for name, qty := range inv.Items {
		fmt.Printf("%s x%d - %s\n", name, qty, inv.Refs[name].GetDescription())
	}
}
