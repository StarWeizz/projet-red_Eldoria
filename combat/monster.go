package combat

import (
	"math/rand"
	"time"
)

// Monster déjà défini
type Monster struct {
	Name    string
	HP      int
	Attack  int
	Defense int
}

func (m *Monster) IsAlive() bool {
	return m.HP > 0
}

func (m *Monster) TakeDamage(dmg int) {
	m.HP -= dmg
	if m.HP < 0 {
		m.HP = 0
	}
}

// NewRandomMonster crée un monstre aléatoire
func NewRandomMonster() *Monster {
	rand.Seed(time.Now().UnixNano())

	monsters := []*Monster{
		{Name: "Gobelin", HP: 20, Attack: 5, Defense: 2},
		{Name: "Loup", HP: 15, Attack: 7, Defense: 1},
		{Name: "Troll", HP: 30, Attack: 8, Defense: 3},
	}

	return monsters[rand.Intn(len(monsters))]
}
