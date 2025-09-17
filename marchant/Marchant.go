package marchant
<<<<<<< HEAD
=======
package marchant
>>>>>>> f8fb55b (Refactoring files)

import (
	items "eldoria/items"
	createcharacter "eldoria/player"
	"fmt"
)

type Merchant struct {
	Name  string
	Stock map[string]items.Item
}

func NewMerchant(name string) *Merchant {
	return &Merchant{
		Name: name,
		Stock: map[string]items.Item{
			"Poison potion":    items.PotionsList["Poison potion"],
			"Heal potion":      items.PotionsList["Heal potion"],
			"Bâton":            items.CraftingItems["Bâton"],
			"Pierre":           items.CraftingItems["Pierre"],
			"Papier":           items.CraftingItems["Papier"],
			"Parchemin":        items.CraftingItems["Parchemin"],
			"Ecaille d'Azador": items.CraftingItems["Ecaille d'Azador"],
		},
	}
}

func (m *Merchant) ShowStock() {
	fmt.Printf("=== Boutique de %s ===\n", m.Name)
<<<<<<< HEAD
	fmt.Println("Articles disponibles :")

	i := 1
	for _, it := range m.Stock {
		fmt.Printf("%d. %s - %d or\n", i, it.GetName(), it.GetPrice())
		i++
	}
	fmt.Println()
=======
	for _, it := range m.Stock {
		fmt.Printf("- %s : %d or\n", it.GetName(), it.GetPrice())
	}
>>>>>>> f8fb55b (Refactoring files)
}

func (m *Merchant) Buy(p *createcharacter.Character, itemName string) {
	it, ok := m.Stock[itemName]
	if !ok {
		fmt.Println("Cet objet n’est pas en vente.")
		return
	}

	if p.Gold.Get() < it.GetPrice() {
		fmt.Println("Pas assez d’or pour acheter", itemName)
		return
	}

	p.Gold.Remove(it.GetPrice())
	p.Inventory.Add(it, 1)
	fmt.Printf("Achat réussi : %s ajouté à ton inventaire !\n", itemName)
<<<<<<< HEAD
}
=======
}
>>>>>>> f8fb55b (Refactoring files)
