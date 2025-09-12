package items

type Potion struct {
	Name        string
	Description string
	Heal        int
	Mana        int
	Poison      int
	Money       int
	//DropRate    int
}

func (p Potion) GetName() string {
	return p.Name
}
func (p Potion) GetDescription() string {
	return p.Description
}
func (w Weapon) GetPrice() int {
	return w.Price
}

var PotionsList = map[string]Potion{
	"Heal potion": {
		Name:        "Heal potion",
		Description: "Une potion qui soigne 20 PV.",
		Heal:        20,
		Mana:        0,
		Poison:      0,
		Money:       30,
		//DropRate:    60,
	},
	"Poison potion": {
		Name:        "Poison potion",
		Description: "Une potion empoisonnée infligeant 10 dégâts sur la durée.",
		Heal:        0,
		Mana:        0,
		Poison:      10,
		Money:       40,
		//	DropRate:    30,
	},
}
