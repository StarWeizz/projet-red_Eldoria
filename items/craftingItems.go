package items

type CraftingItem struct {
	Name        string
	Description string
	// DropRate int
}

var CraftingItems = map[string]CraftingItem{
	"Bâton": {
		Name:        "Bâton",
		Description: "Un simple morceau de bois, utile pour fabriquer une arme basique.",
		// DropRate: 60,
	},
	"Pierre": {
		Name:        "Pierre",
		Description: "Une pierre robuste, peut servir comme matériau pour une arme basique.",
		// DropRate: 30,
	},
	"Ecaille d'Azador": {
		Name:        "Ecaille d'Azador",
		Description: "Une écaille rare d’un ancien dragon, ingrédient précieux pour une arme basique.",
		// DropRate: 40,
	},
}
