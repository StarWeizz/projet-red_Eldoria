package combat

import (
	createcharacter "eldoria/player"
	"fmt"
)

func StartCombat(h *createcharacter.Character, m *Monster) {
	fmt.Printf("⚔️  Un combat commence ! %s (%s) affronte %s\n", h.Name, h.Class, m.Name)

	for h.CurrentHP > 0 && m.HP > 0 {
		// ----- Tour du héros -----
		fmt.Printf("\n--- Tour de %s ---\n", h.Name)
		damageToMonster := 10 - m.Defense // TODO: remplacer par h.Attack quand tu ajoutes l'attribut
		if damageToMonster < 0 {
			damageToMonster = 0
		}
		m.TakeDamage(damageToMonster)
		fmt.Printf("%s attaque %s et inflige %d dégâts ! (PV monstre: %d)\n", h.Name, m.Name, damageToMonster, m.HP)

		if !m.IsAlive() {
			fmt.Printf("%s est vaincu ! 🏆\n", m.Name)
			h.Gold.Add(20) // récompense en or
			fmt.Printf("Vous gagnez 20 or. Total: %d\n", h.Gold.Get())
			break
		}

		// ----- Tour du monstre -----
		fmt.Printf("\n--- Tour de %s ---\n", m.Name)
		damageToHero := m.Attack - 3 // TODO: remplacer 3 par h.Defense quand dispo
		if damageToHero < 0 {
			damageToHero = 0
		}
		h.CurrentHP -= damageToHero
		fmt.Printf("%s attaque %s et inflige %d dégâts ! (PV héros: %d/%d)\n", m.Name, h.Name, damageToHero, h.CurrentHP, h.MaxHP)

		if h.CurrentHP <= 0 {
			fmt.Printf("%s est vaincu... 💀\n", h.Name)
			break
		}
	}
}
