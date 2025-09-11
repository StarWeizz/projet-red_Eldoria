package items

type weapon struct {
	Name        string
	Description string
	Damage      int
}

func (w weapon) GetName() string {
	return w.Name
}
func (w weapon) GetDescription() string {
	return w.Description
}

var WeaponList = map[string]weapon{
	"Epée simple": {
		Name:        "Epée simple",
		Description: "",
		Damage:      5,
		//	DropRate:    60,
	},
	"Double épée": {
		Name:        "Double épée",
		Description: "",
		Damage:      10,
		//	DropRate:    20,
	},
	"Epée Démoniaque": {
		Name:        "Epée Démoniaque",
		Description: "",
		Damage:      15,
		//	DropRate:    10,
	},

	"Grimoire Simple": {
		Name:        "Grimoire Simple",
		Description: "",
		Damage:      5,
		//	DropRate:    5,
	},
	"Livre des Mage": {
		Name:        "Livre des Mage",
		Description: "",
		Damage:      10,
		//	DropRate:    5,
	},
	"Livre des Ombre": {
		Name:        "Livre des Ombre",
		Description: "",
		Damage:      15,
		//	DropRate:    5,
	},

	"Couteaux de Chasse": {
		Name:        "Couteaux de Chasse",
		Description: "",
		Damage:      5,
		//	DropRate:    5,
	},
	"épée court runique": {
		Name:        "épée court runique",
		Description: "",
		Damage:      10,
		//	DropRate:    5,
	},
	"Dague de l’ombre": {
		Name:        "Dague de l’ombre",
		Description: "",
		Damage:      15,
		//	DropRate:    5,
	},
}
