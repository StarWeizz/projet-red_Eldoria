package items

type Recipe struct {
	Result string
	Needs  []string
}

var Recipes = []Recipe{
	{
		Result: "Lame rouillé",
		Needs:  []string{"Bâton", "Pierre", "Pierre"},
	},
	{
		Result: "Grimoire",
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
