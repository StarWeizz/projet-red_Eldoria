package forgeron
<<<<<<< HEAD
=======
package forgeron
>>>>>>> f8fb55b (Refactoring files)

import (
	"eldoria/items"
	createcharacter "eldoria/player"
	"fmt"
)

type Blacksmith struct {
	Name string
	Shop map[string]items.Item
}

func NewBlacksmith(name string) *Blacksmith {
	return &Blacksmith{
		Name: name,
		Shop: map[string]items.Item{
			"Lame rouillé":       items.WeaponList["Lame rouillé"],
			"Grimoire":           items.WeaponList["Grimoire"],
			"Couteaux de Chasse": items.WeaponList["Couteaux de Chasse"],
		},
	}
}

func (b *Blacksmith) ShowStock() {
	fmt.Printf("=== Forge de %s ===\n", b.Name)
	for _, it := range b.Shop {
		fmt.Printf("- %s : %d or\n", it.GetName(), it.GetPrice())
	}
}

func (b *Blacksmith) Buy(p *createcharacter.Character, weaponName string) {
	weapon, ok := b.Shop[weaponName]
	if !ok {
		fmt.Println("Cette arme n’est pas en vente.")
		return
	}

	if p.Gold.Get() < weapon.GetPrice() {
		fmt.Println("Pas assez d’or.")
		return
	}

	p.Gold.Remove(weapon.GetPrice())
	p.Inventory.Add(weapon, 1)
	fmt.Printf("Tu as acheté %s pour %d or.\n", weapon.GetName(), weapon.GetPrice())
<<<<<<< HEAD
}
=======
}
>>>>>>> f8fb55b (Refactoring files)
