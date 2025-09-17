package items

type Potion struct {
	Name        string
	Description string
	Heal        int
	Mana        int
	Poison      int
	Price       int
	//DropRate    int
}

func (p Potion) GetName() string {
	return p.Name
}
func (p Potion) GetDescription() string {
	return p.Description
}
func (p Potion) GetPrice() int {
	return p.Price
}

var PotionsList = map[string]Potion{
	"Heal potion": {
		Name:        "Heal potion",
		Description: "Une potion qui soigne 45 PV.",
		Heal:        45,
		Mana:        0,
		Poison:      0,
		Price:       20,
		//DropRate:    60,
	},
	"Poison potion": {
		Name:        "Poison potion",
		Description: "Une potion empoisonnée infligeant 10 dégâts sur la durée.",
		Heal:        0,
		Mana:        0,
		Poison:      10,
		Price:       40,
		//	DropRate:    30,
	},
}
