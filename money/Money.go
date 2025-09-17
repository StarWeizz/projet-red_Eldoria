package money

type Money struct {
	Amount int
}

func NewMoney(start int) *Money {
	return &Money{Amount: start}
}

func (m *Money) Add(amount int) {
	m.Amount += amount
}

func (m *Money) Remove(amount int) bool {
	if m.Amount < amount {
		return false
	}
	m.Amount -= amount
	return true
}

func (m *Money) Get() int {
	return m.Amount
}
