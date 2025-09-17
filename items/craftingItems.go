package items

type CraftingItem struct {
	Name        string
	Description string
	Price       int
	// DropRate int
}

func (c CraftingItem) GetPrice() int {
	return c.Price
}

func (c CraftingItem) GetName() string {
	return c.Name
}

func (c CraftingItem) GetDescription() string {
	return c.Description
}

var CraftingItems = map[string]CraftingItem{
	"Bâton": {
		Name:        "Bâton",
		Description: "Un simple morceau de bois, utile pour fabriquer une arme basique.",
		Price:       3,
		// DropRate: 60,
	},
	"Pierre": {
		Name:        "Pierre",
		Description: "Une pierre robuste, peut servir comme matériau pour une arme basique.",
		Price:       3,
		// DropRate: 30,
	},
	"Ecaille d'Azador": {
		Name:        "Ecaille d'Azador",
		Description: "Une écaille rare d’un ancien dragon, ingrédient précieux pour une arme basique.",
		Price:       15,
		// DropRate: 40,
	},
	"Papier": {
		Name:        "Papier",
		Description: "Du Papier simple pour des craft simple",
		Price:       5,
		// DropRate: 40,
	},
	"Parchemin": {
		Name:        "Parchemin",
		Description: "Un Parchemin pour fabriquer des arme basique",
		Price:       10,
		// DropRate: 40,
	},
}
