package items

type Weapon struct {
	Name        string
	Description string
	Damage      int
	Price       int
	Level       int
	Category    string
}

func (w Weapon) GetName() string {
	return w.Name
}
func (w Weapon) GetDescription() string {
	return w.Description
}
func (w Weapon) GetPrice() int {
	return w.Price
}
func (w Weapon) GetDamage() int {
	return w.Damage
}

// Liste des armes disponibles
var WeaponList = map[string]Weapon{
	"Lame rouillée": {
		Name:        "Lame rouillée",
		Description: "Une vieille lame usée par le temps.",
		Damage:      5,
		Price:       35,
		Level:       1,
		Category:    "Sword",
	},
	"Épée de chevalier": {
		Name:        "Épée de chevalier",
		Description: "Une épée bien forgée, robuste.",
		Damage:      10,
		Level:       2,
		Category:    "Sword",
	},
	"Épée démoniaque": {
		Name:        "Épée démoniaque",
		Description: "Une arme maudite aux pouvoirs obscurs.",
		Damage:      15,
		Level:       3,
		Category:    "Sword",
	},
	"Grimoire": {
		Name:        "Grimoire",
		Description: "",
		Damage:      5,
		Price:       35,
		Level:       1,
		Category:    "Magie",
		//	DropRate:    5,
	},
	"Livre de Magie": {
		Name:        "Livre de Magie",
		Description: "",
		Damage:      10,
		Level:       2,
		Category:    "Sword",
		//	DropRate:    5,
	},
	"Livre des Ombre": {
		Name:        "Livre des Ombre",
		Description: "",
		Damage:      15,
		Level:       3,
		Category:    "Sword",
		//	DropRate:    5,
	},

	"Couteaux de Chasse": {
		Name:        "Couteaux de Chasse",
		Description: "",
		Damage:      5,
		Price:       35,
		Level:       1,
		Category:    "Dague",
		//	DropRate:    5,
	},
	"épée court runique": {
		Name:        "épée court runique",
		Description: "",
		Damage:      10,
		Level:       2,
		Category:    "Dague",

		//	DropRate:    5,
	},
	"Dague de l’ombre": {
		Name:        "Dague de l’ombre",
		Description: "",
		Damage:      15,
		Level:       3,
		Category:    "Dague",
		//	DropRate:    5,
	},
}
