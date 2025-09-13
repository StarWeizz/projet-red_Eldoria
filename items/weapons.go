package items

//import (
//money "eldoria/money"
//)

type weapon struct {
	Name        string
	Description string
	Damage      int
	Price       int
}

func (w weapon) GetName() string {
	return w.Name
}
func (w weapon) GetDescription() string {
	return w.Description
}
func (w weapon) GetPrice() int {
	return w.Price
}

var WeaponList = map[string]weapon{
	"Lame d’entrainement": {
		Name:        "Lame d’entrainement",
		Description: "",
		Damage:      5,
		Price:       35,
		//	DropRate:    60,
	},
	"épée de chevalier": {
		Name:        "épée de chevalier",
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

	"Grimoire": {
		Name:        "Grimoire",
		Description: "",
		Damage:      5,
		Price:       35,
		//	DropRate:    5,
	},
	"Livre de Magie": {
		Name:        "Livre de Magie",
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
		Price:       35,
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
