package npcs

import (
	items "eldoria/items"
	createcharacter "eldoria/player"
	"fmt"
)

type QuestStep struct {
	ID          int
	Title       string
	Description string
	Condition   func(*createcharacter.Character) bool
	Reward      QuestReward
	Completed   bool
}

type QuestReward struct {
	Money int
	Items map[string]int // item name -> quantity
	Weapon items.Item
}

type Quest struct {
	ID          string
	Title       string
	Description string
	Steps       []QuestStep
	CurrentStep int
	Completed   bool
}

type NPC struct {
	Name        string
	Location    string
	Description string
	Quests      []Quest
	Dialogue    map[string]string // context -> dialogue
}

func NewNPC(name, location, description string) *NPC {
	return &NPC{
		Name:        name,
		Location:    location,
		Description: description,
		Quests:      []Quest{},
		Dialogue:    make(map[string]string),
	}
}

func (npc *NPC) AddQuest(quest Quest) {
	npc.Quests = append(npc.Quests, quest)
}

func (npc *NPC) ShowAvailableQuests(player *createcharacter.Character) {
	fmt.Printf("\n=== %s ===\n", npc.Name)
	fmt.Printf("📍 Lieu : %s\n", npc.Location)
	fmt.Printf("💬 %s\n\n", npc.Description)

	hasActiveQuest := false
	for i, quest := range npc.Quests {
		if !quest.Completed {
			hasActiveQuest = true
			fmt.Printf("## ⚔️ %s\n", quest.Title)
			fmt.Printf("%s\n\n", quest.Description)

			// Afficher les étapes
			for j, step := range quest.Steps {
				status := "❌"
				if step.Completed {
					status = "✅"
				} else if j == quest.CurrentStep {
					status = "🔄"
				}

				fmt.Printf("%s %d. %s\n", status, j+1, step.Title)
				fmt.Printf("   📋 %s\n", step.Description)

				// Afficher la récompense
				if step.Reward.Money > 0 || len(step.Reward.Items) > 0 || step.Reward.Weapon != nil {
					fmt.Printf("   🎁 Récompense : ")
					if step.Reward.Money > 0 {
						fmt.Printf("%d or ", step.Reward.Money)
					}
					if step.Reward.Weapon != nil {
						fmt.Printf("Arme: %s ", step.Reward.Weapon.GetName())
					}
					for itemName, qty := range step.Reward.Items {
						fmt.Printf("%s x%d ", itemName, qty)
					}
					fmt.Println()
				}
				fmt.Println()
			}

			// Montrer l'action possible
			if quest.CurrentStep < len(quest.Steps) {
				currentStep := quest.Steps[quest.CurrentStep]
				if currentStep.Condition(player) {
					fmt.Printf("✨ Vous pouvez valider l'étape actuelle !\n")
				} else {
					fmt.Printf("⏳ Conditions non remplies pour l'étape actuelle.\n")
				}
			}

			npc.Quests[i] = quest
			break
		}
	}

	if !hasActiveQuest {
		fmt.Printf("💤 %s n'a pas de quêtes disponibles pour le moment.\n", npc.Name)
	}
}

func (npc *NPC) GetQuestInfo(player *createcharacter.Character) string {
	result := fmt.Sprintf("=== %s ===\n", npc.Name)
	result += fmt.Sprintf("📍 Lieu : %s\n", npc.Location)
	result += fmt.Sprintf("💬 %s\n\n", npc.Description)

	hasActiveQuest := false
	for _, quest := range npc.Quests {
		if !quest.Completed {
			hasActiveQuest = true
			result += fmt.Sprintf("## ⚔️ %s\n", quest.Title)
			result += fmt.Sprintf("%s\n\n", quest.Description)

			// Afficher les étapes
			for j, step := range quest.Steps {
				status := "❌"
				if step.Completed {
					status = "✅"
				} else if j == quest.CurrentStep {
					status = "🔄"
				}

				result += fmt.Sprintf("%s %d. %s\n", status, j+1, step.Title)
				result += fmt.Sprintf("   📋 %s\n", step.Description)

				// Afficher la récompense
				if step.Reward.Money > 0 || len(step.Reward.Items) > 0 || step.Reward.Weapon != nil {
					result += "   🎁 Récompense : "
					if step.Reward.Money > 0 {
						result += fmt.Sprintf("%d or ", step.Reward.Money)
					}
					if step.Reward.Weapon != nil {
						result += fmt.Sprintf("Arme: %s ", step.Reward.Weapon.GetName())
					}
					for itemName, qty := range step.Reward.Items {
						result += fmt.Sprintf("%s x%d ", itemName, qty)
					}
					result += "\n"
				}
				result += "\n"
			}

			// Montrer l'action possible
			if quest.CurrentStep < len(quest.Steps) {
				currentStep := quest.Steps[quest.CurrentStep]
				if currentStep.Condition(player) {
					result += "✨ Vous pouvez valider l'étape actuelle !\n"
				} else {
					result += "⏳ Conditions non remplies pour l'étape actuelle.\n"
				}
			}

			break
		}
	}

	if !hasActiveQuest {
		result += fmt.Sprintf("💤 %s n'a pas de quêtes disponibles pour le moment.\n", npc.Name)
	}

	return result
}

func (npc *NPC) ValidateQuestStep(player *createcharacter.Character, questID string) bool {
	for i, quest := range npc.Quests {
		if quest.ID == questID && !quest.Completed && quest.CurrentStep < len(quest.Steps) {
			currentStep := &npc.Quests[i].Steps[quest.CurrentStep]

			if currentStep.Condition(player) {
				// Donner la récompense silencieusement
				if currentStep.Reward.Money > 0 {
					player.Gold.Add(currentStep.Reward.Money)
				}

				if currentStep.Reward.Weapon != nil {
					player.Inventory.Add(currentStep.Reward.Weapon, 1)
				}

				for itemName, qty := range currentStep.Reward.Items {
					if item, exists := items.CraftingItems[itemName]; exists {
						player.Inventory.Add(item, qty)
					}
				}

				// Marquer l'étape comme complétée
				currentStep.Completed = true

				// Passer à l'étape suivante
				npc.Quests[i].CurrentStep++

				// Vérifier si la quête est terminée
				if npc.Quests[i].CurrentStep >= len(quest.Steps) {
					npc.Quests[i].Completed = true
				}

				return true
			} else {
				return false
			}
		}
	}
	return false
}

func (npc *NPC) Interact(player *createcharacter.Character) {
	npc.ShowAvailableQuests(player)
}