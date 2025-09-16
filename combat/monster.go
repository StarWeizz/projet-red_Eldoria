package combat

import (
	"math/rand"
	"time"
)

// Structure Monster
type Monster struct {
	Name    string
	HP      int
	Attack  int
	Defense int
}

// Vérifie si le monstre est encore vivant
func (m *Monster) IsAlive() bool {
	return m.HP > 0
}

// Inflige des dégâts au monstre et retourne la valeur infligée
func (m *Monster) TakeDamage(damage int) int {
	if damage < 0 {
		damage = 0
	}
	m.HP -= damage
	if m.HP < 0 {
		m.HP = 0
	}
	return damage
}

// ---------- BESTIAIRE ----------

// Liste des monstres de base
var Bestiary = []*Monster{
	{Name: "Gobelin", HP: 30, Attack: 6, Defense: 2},
	{Name: "Orc", HP: 50, Attack: 10, Defense: 4},
	{Name: "Loup", HP: 40, Attack: 8, Defense: 3},
	{Name: "Troll", HP: 70, Attack: 12, Defense: 5},
}

// Tire un monstre aléatoire dans le bestiaire
func GetRandomMonster() *Monster {
	rand.Seed(time.Now().UnixNano())
	monster := Bestiary[rand.Intn(len(Bestiary))]
	// Retourne une copie pour ne pas modifier l'original
	return &Monster{
		Name:    monster.Name,
		HP:      monster.HP,
		Attack:  monster.Attack,
		Defense: monster.Defense,
	}
}
