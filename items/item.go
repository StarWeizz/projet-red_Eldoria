package items

// Item est l’interface de base pour tous les objets
type Item interface {
	GetName() string
	GetDescription() string
	GetPrice() int
}
