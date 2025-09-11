package items

type Recipe struct {
	Result string
	Needs  []string
}

func (c CraftingItem) GetName() string {
	return c.Name
}
func (c CraftingItem) GetDescription() string {
	return c.Description
}

var Recipes = []Recipe{
	{
		Result: "Epée simple",
		Needs:  []string{"Bâton", "Pierre", "Pierre"},
	},
	{
		Result: "Grimoire Simple",
		Needs:  []string{"Parchemin", "Parchemin"},
	},
	{
		Result: "Parchemin",
		Needs:  []string{"Papier", "Papier"},
	},
	{
		Result: "Couteaux de Chasse",
		Needs:  []string{"Bâton", "Ecaille d'Azador"},
	},
}
