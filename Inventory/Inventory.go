package inventory

import (
	"eldoria/items"
	"fmt"
)

type Inventory struct {
	Items map[string]int        // quantité
	Refs  map[string]items.Item // références vers les objets
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

func (inv *Inventory) GetInventoryString() string {
	if len(inv.Items) == 0 {
		return "🎒 Inventaire vide\n\nVous n'avez aucun objet dans votre sac à dos."
	}

	result := "🎒 Contenu de votre inventaire :\n\n"
	for name, qty := range inv.Items {
		result += fmt.Sprintf("• %s x%d - %s\n", name, qty, inv.Refs[name].GetDescription())
	}
	result += "\nAppuyez sur [i] pour fermer l'inventaire."

	return result
}

func (inv *Inventory) HasItem(itemName string, minQty int) bool {
	if qty, exists := inv.Items[itemName]; exists {
		return qty >= minQty
	}
	return false
}
