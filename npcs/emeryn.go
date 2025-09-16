package npcs

import (
	items "eldoria/items"
	createcharacter "eldoria/player"
	"fmt"
)

func CreateEmeryn() *NPC {
	emeryn := NewNPC(
		"Emeryn",
		"Centre du village",
		"Un guide sage qui accueille les nouveaux aventuriers et leur enseigne les bases de la survie.",
	)

	// Créer la quête d'introduction
	introQuest := Quest{
		ID:          "intro_quest",
		Title:       "Quête d'introduction",
		Description: "Apprenez les bases de l'aventure et préparez-vous pour votre première mission.",
		CurrentStep: 0,
		Completed:   false,
		Steps: []QuestStep{
			{
				ID:          1,
				Title:       "Trouver Emeryn au centre du village",
				Description: "Rendez-vous au centre du village pour rencontrer Emeryn.",
				Condition: func(player *createcharacter.Character) bool {
					// Cette condition sera toujours vraie une fois que le joueur interagit avec Emeryn
					return true
				},
				Reward: QuestReward{
					Money:  50,
					Weapon: items.WeaponList["Lame rouillé"], // Première arme
					Items:  make(map[string]int),
				},
				Completed: false,
			},
			{
				ID:          2,
				Title:       "Tuer votre premier Azador",
				Description: "Combattez un Azador et récoltez de la pierre et des bâtons.",
				Condition: func(player *createcharacter.Character) bool {
					// Vérifier si le joueur a de la pierre et des bâtons dans son inventaire
					hasStone := player.Inventory.HasItem("Pierre", 1)
					hasStick := player.Inventory.HasItem("Bâton", 1)
					return hasStone && hasStick
				},
				Reward: QuestReward{
					Money: 100,
					Items: make(map[string]int),
				},
				Completed: false,
			},
			{
				ID:          3,
				Title:       "Trouver Valenric le forgeron et crafter votre première arme",
				Description: "Apportez pierre et bâtons au forgeron pour créer une arme améliorée.",
				Condition: func(player *createcharacter.Character) bool {
					// Cette condition sera gérée par le forgeron lui-même
					// On peut vérifier si le joueur a une arme craftée spécifique
					return player.Inventory.HasItem("Épée en fer", 1) ||
						   player.Inventory.HasItem("Hache en fer", 1) ||
						   player.Inventory.HasItem("Arc en bois", 1)
				},
				Reward: QuestReward{
					Money: 150,
					Items: make(map[string]int),
				},
				Completed: false,
			},
		},
	}

	emeryn.AddQuest(introQuest)

	// Ajouter des dialogues contextuels
	emeryn.Dialogue["greeting"] = "Bienvenue, jeune aventurier ! Je suis Emeryn, votre guide dans ce monde dangereux."
	emeryn.Dialogue["quest_available"] = "J'ai une mission parfaite pour un débutant comme vous."
	emeryn.Dialogue["quest_in_progress"] = "Comment avancez-vous dans votre quête ?"
	emeryn.Dialogue["quest_completed"] = "Félicitations ! Vous êtes maintenant prêt pour de plus grands défis."

	return emeryn
}

// EmerynPhase indique la phase d'interaction actuelle avec Emeryn
var EmerynPhase int = 0
var EmerynInteractionStarted bool = false

func (npc *NPC) EmerynSpecialInteraction(player *createcharacter.Character) {
	// Ne rien faire ici, l'interaction sera gérée par GetEmerynMessage
}

func (npc *NPC) GetEmerynMessage(player *createcharacter.Character) string {
	// Vérifier la quête d'introduction
	for _, quest := range npc.Quests {
		if quest.ID == "intro_quest" && !quest.Completed {
			if quest.CurrentStep == 0 {
				if EmerynPhase == 0 {
					// Phase 1 : Message d'accueil personnalisé selon la classe
					EmerynInteractionStarted = true
					return npc.GetWelcomeMessage(player) + "\n\nAppuyez sur [Espace] pour continuer."
				} else if EmerynPhase == 1 {
					// Phase 2 : Instructions après avoir reçu l'arme
					npc.ValidateQuestStep(player, "intro_quest")

					// Donner 20 XP pour avoir visité Emeryn
					expMessage := player.AddExperience(20)

					EmerynInteractionStarted = false // Réinitialiser après la quête
					EmerynPhase = 2 // Marquer que la lame a été donnée

					message := "Je viens de te donner 1 lame rouillée.\n\nVa à la sortie du village, tu trouveras un Azador, bats-le et ensuite va récolter des pierres et des bâtons puis trouve Valenric le forgeron."
					if expMessage != "" {
						message += "\n\n" + expMessage
					}
					return message
				}
			} else if quest.CurrentStep == 1 {
				// Étape 2 : Tuer un Azador
				return "Trouve un Azador à la sortie du village et récolte des matériaux !"
			} else if quest.CurrentStep == 2 {
				// Étape 3 : Voir le forgeron
				return "Va voir Valenric le forgeron avec tes matériaux !"
			}
		}
	}

	// Si la quête est terminée ou si on a déjà donné la lame
	if EmerynPhase >= 2 {
		return "Bon voyage, aventurier ! N'oublie pas de voir Valenric le forgeron."
	}

	return "Bon voyage, aventurier !"
}

func (npc *NPC) GetWelcomeMessage(player *createcharacter.Character) string {
	baseMessage := "Bienvenue dans le village d'Eldoria, écoute mes conseils.\n\n"

	switch player.Class {
	case "Guerrier":
		return baseMessage + "Jeune guerrier, tes avantages :\n• Combat corps à corps\n• Haute résistance\n• Armes lourdes"
	case "Mage":
		return baseMessage + "Jeune mage, tes avantages :\n• Maîtrise magique\n• Sorts puissants\n• Intelligence élevée"
	case "Chasseur":
		return baseMessage + "Jeune chasseur, tes avantages :\n• Combat à distance\n• Agilité élevée\n• Précision redoutable"
	default:
		return baseMessage + "Jeune aventurier, prépare-toi pour ton voyage."
	}
}

func (npc *NPC) AdvanceEmerynPhase() {
	// Ne faire avancer que si l'interaction a été démarrée avec E et qu'on n'a pas déjà donné la lame
	if EmerynInteractionStarted && EmerynPhase < 2 {
		EmerynPhase++
	}
}

func (npc *NPC) CanAdvanceEmeryn() bool {
	// Peut avancer seulement si l'interaction est démarrée et qu'on n'a pas encore donné la lame
	return EmerynInteractionStarted && EmerynPhase < 2
}

func (npc *NPC) ExplainClasses(player *createcharacter.Character) {
	fmt.Printf("\n=== Explication des classes ===\n")
	fmt.Printf("📚 Laissez-moi vous expliquer les différentes classes d'aventuriers :\n\n")

	fmt.Printf("⚔️ **Guerrier**\n")
	fmt.Printf("   • +3 Force, +2 Endurance, +1 Agilité\n")
	fmt.Printf("   • Spécialisé dans le combat au corps à corps\n")
	fmt.Printf("   • Utilise épées, haches et armures lourdes\n\n")

	fmt.Printf("🏹 **Archer**\n")
	fmt.Printf("   • +3 Agilité, +2 Force, +1 Endurance\n")
	fmt.Printf("   • Expert en combat à distance\n")
	fmt.Printf("   • Utilise arcs, flèches et armures légères\n\n")

	fmt.Printf("🔮 **Mage**\n")
	fmt.Printf("   • +3 Intelligence, +2 Agilité, +1 Force\n")
	fmt.Printf("   • Maître des sorts et de la magie\n")
	fmt.Printf("   • Utilise bâtons magiques et sorts puissants\n\n")

	fmt.Printf("Votre classe actuelle : %s\n", player.Class)
	fmt.Printf("Vos caractéristiques actuelles :\n")
	fmt.Printf("• Niveau : %d\n", player.Level)
	fmt.Printf("• HP : %d/%d\n", player.CurrentHP, player.MaxHP)
	fmt.Printf("• Or : %d\n", player.Gold.Get())
}

func (npc *NPC) GetClassExplanation(player *createcharacter.Character) string {
	result := "\n=== Explication des classes ===\n"
	result += "📚 Laissez-moi vous expliquer les différentes classes d'aventuriers :\n\n"

	result += "⚔️ **Guerrier**\n"
	result += "   • +3 Force, +2 Endurance, +1 Agilité\n"
	result += "   • Spécialisé dans le combat au corps à corps\n"
	result += "   • Utilise épées, haches et armures lourdes\n\n"

	result += "🏹 **Archer**\n"
	result += "   • +3 Agilité, +2 Force, +1 Endurance\n"
	result += "   • Expert en combat à distance\n"
	result += "   • Utilise arcs, flèches et armures légères\n\n"

	result += "🔮 **Mage**\n"
	result += "   • +3 Intelligence, +2 Agilité, +1 Force\n"
	result += "   • Maître des sorts et de la magie\n"
	result += "   • Utilise bâtons magiques et sorts puissants\n\n"

	result += fmt.Sprintf("Votre classe actuelle : %s\n", player.Class)
	result += "Vos caractéristiques actuelles :\n"
	result += fmt.Sprintf("• Niveau : %d\n", player.Level)
	result += fmt.Sprintf("• HP : %d/%d\n", player.CurrentHP, player.MaxHP)
	result += fmt.Sprintf("• Or : %d\n", player.Gold.Get())

	return result
}