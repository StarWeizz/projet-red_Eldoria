package combat

import (
	createcharacter "eldoria/player"
	"math/rand"
)

// Retourne : victoire, tableau dégâts joueur->monstre, tableau dégâts monstre->joueur, bool si fuite
func StartCombat(h *createcharacter.Character, m *Monster, playerChoice func(h *createcharacter.Character, m *Monster) string) (bool, []int, []int, bool) {
	var playerDamages []int
	var monsterDamages []int
	for h.CurrentHP > 0 && m.HP > 0 {
		// Choix du joueur
		action := playerChoice(h, m)
		if action == "flee" {
			// Le joueur fuit le combat
			return false, playerDamages, monsterDamages, true
		}
		if action == "heal" {
			// Vérifier et consommer la potion de Heal
			if h.Inventory.HasItem("Heal potion", 1) {
				// Consommer la potion
				potion := h.Inventory.Refs["Heal potion"]
				if h.Inventory.Remove(potion, 1) {
					healAmount := 45
					h.CurrentHP += healAmount
					if h.CurrentHP > h.MaxHP {
						h.CurrentHP = h.MaxHP
					}
				}
				continue
				// Pas d'attaque ce tour, on passe au tour du monstre
			} // Si pas de potion, rien ne se passe (déjà géré côté UI)
		} else if action == "attack" {
			// Tour du héros
			baseDamage := 10
			maxWeaponDamage := 0
			for name, qty := range h.Inventory.Items {
				if qty > 0 {
					if weapon, ok := h.Inventory.Refs[name].(interface{ GetDamage() int }); ok {
						dmg := weapon.GetDamage()
						if dmg > maxWeaponDamage {
							maxWeaponDamage = dmg
						}
					}
				}
			}
			damageToMonster := baseDamage + maxWeaponDamage - m.Defense

			if damageToMonster < 0 {
				damageToMonster = 0
			}
			m.TakeDamage(damageToMonster)
			playerDamages = append(playerDamages, damageToMonster)

			if !m.IsAlive() {
				h.Gold.Add(20)
				return true, playerDamages, monsterDamages, false
			}
		}

		// Tour du monstre
		damageToHero := m.Attack - 3 // remplacer 3 par h.Defense si ajouté
		special := false
		// Si c'est Maximor, il a une chance d'attaque spéciale
		if m.Name == "Maximor" {
			// 45% de chance d'attaque spéciale
			if rand.Intn(100) < 45 {
				damageToHero += 35 // dégâts bonus
				special = true
			}
			if rand.Intn(100) < 5 {
				damageToHero += 75 // dégâts bonus
				special = true
			}
		}
		if damageToHero < 0 {
			damageToHero = 0
		}
		h.CurrentHP -= damageToHero
		// On encode l'attaque spéciale dans le signe du dégât (pour l'affichage)
		if special {
			monsterDamages = append(monsterDamages, -damageToHero) // négatif = attaque spéciale
		} else {
			monsterDamages = append(monsterDamages, damageToHero)
		}

		if h.CurrentHP <= 0 {
			return false, playerDamages, monsterDamages, false
		}
	}
	return false, playerDamages, monsterDamages, false
}
