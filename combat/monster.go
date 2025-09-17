package combat

import (
	createcharacter "eldoria/player"
	"fmt"
	"math/rand"
	"time"
)

type Boss struct {
	Monster
	SpecialCooldown int // compteur pour attaque spéciale
}

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

// Crée le Boss Maximor
func NewMaximor() *Boss {
	return &Boss{
		Monster: Monster{
			Name:    "Maximor",
			HP:      150,
			Attack:  10,
			Defense: 3,
		},
		SpecialCooldown: 0,
	}
}

func (b *Boss) AttackHero(h *createcharacter.Character) int {
	rand.Seed(time.Now().UnixNano())

	// Vérifier si attaque spéciale prête
	if b.SpecialCooldown <= 0 && rand.Intn(100) < 30 { // 30% de chance
		b.SpecialCooldown = 3    // cooldown 3 tours
		damage := b.Attack*2 - 2 // dégâts spéciaux, moins la défense du héros
		if damage < 0 {
			damage = 0
		}
		fmt.Printf("💥 %s utilise Frappe Dévastatrice et inflige %d dégâts !\n", b.Name, damage)
		return damage
	}

	// Attaque normale
	damage := b.Attack - 3 // remplacer 3 par h.Defense si déf finie
	if damage < 0 {
		damage = 0
	}

	// Réduire le cooldown si nécessaire
	if b.SpecialCooldown > 0 {
		b.SpecialCooldown--
	}

	return damage
}

func NewRandomMonster() *Monster {
	rand.Seed(time.Now().UnixNano())

	monsters := []*Monster{
		{Name: "Apprenti Azador", HP: 25, Attack: 5, Defense: 2},
		{Name: "Azador", HP: 35, Attack: 7, Defense: 3},
		{Name: "Azador Chevalier", HP: 50, Attack: 9, Defense: 4},
	}

	return monsters[rand.Intn(len(monsters))]
}
