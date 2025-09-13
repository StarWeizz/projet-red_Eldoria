package items

type Item interface {
	GetName() string
	GetDescription() string
	GetPrice() int
}
