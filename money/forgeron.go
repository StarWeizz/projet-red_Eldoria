package money

import (
	"eldoria/items"
	//	"eldoria/player"
	"fmt"
)

type Blacksmith struct {
	Name string
	Shop map[string]items.Item // il peut aussi vendre des armes
}

func NewBlacksmith(name string) *Blacksmith {
	return &Blacksmith{
		Name: name,
		Shop: map[string]items.Item{
			"Lame d’entrainement": items.WeaponList["Lame d’entrainement"],
			"Grimoire":            items.WeaponList["Grimoire"],
			"Couteaux de Chasse":  items.WeaponList["Couteaux de Chasse"],
		},
	}
}

// Améliorer une arme équipée
func (b *Blacksmith) UpgradeWeapon(p *player.Player, weaponName string) {
	weaponItem, ok := p.Inventory.Refs[weaponName]
	if !ok {
		fmt.Println("Tu n’as pas cette arme dans ton inventaire.")
		return
	}

	weapon, ok := weaponItem.(items.Weapon)
	if !ok {
		fmt.Println("Cet objet n’est pas une arme.")
		return
	}

	// Règles : prix + sacrifier 2 armes du niveau précédent
	upgradeCost := 10 * weapon.Level // exemple : lvl1→2 coûte 10, lvl2→3 coûte 20, etc.
	requiredLevel := weapon.Level

	if p.Gold.Get() < upgradeCost {
		fmt.Println("Pas assez d’or pour améliorer", weapon.Name)
		return
	}

	// Vérifier si joueur a 2 autres armes du niveau précédent
	count := 0
	for _, it := range p.Inventory.Refs {
		if w, ok := it.(items.Weapon); ok && w.Level == requiredLevel {
			count++
		}
	}
	if count < 2 {
		fmt.Printf("Il te faut au moins 2 armes de niveau %d pour améliorer %s\n", requiredLevel, weapon.Name)
		return
	}

	// Appliquer l’amélioration
	p.Gold.Remove(upgradeCost)
	weapon.Level++
	weapon.Damage += 5 // bonus dégâts par niveau
	p.Inventory.Refs[weaponName] = weapon

	fmt.Printf("%s a été améliorée au niveau %d (+5 dégâts) !\n", weapon.Name, weapon.Level)
}

// Acheter une arme
func (b *Blacksmith) Buy(p *player.Player, weaponName string) {
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
}
