package combat

import (
	createcharacter "eldoria/player"
)

func StartCombat(h *createcharacter.Character, m *Monster) bool {
	for h.CurrentHP > 0 && m.HP > 0 {
		// Tour du héros
		damageToMonster := 10 - m.Defense // temporaire, peut utiliser l’arme du joueur
		if damageToMonster < 0 {
			damageToMonster = 0
		}
		m.TakeDamage(damageToMonster)

		if !m.IsAlive() {
			h.Gold.Add(20)
			return true
		}

		// Tour du monstre
		damageToHero := m.Attack - 3 // remplacer 3 par h.Defense si ajouté
		if damageToHero < 0 {
			damageToHero = 0
		}
		h.CurrentHP -= damageToHero

		if h.CurrentHP <= 0 {
			return false
		}
	}
	return false
}
