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
}
