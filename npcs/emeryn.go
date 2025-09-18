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
					return true
				},
				Reward: QuestReward{
					Money:  50,
					Weapon: items.WeaponList["Lame rouillé"],
					Items:  make(map[string]int),
				},
				Completed: false,
			},
			{
				ID:          2,
				Title:       "Tuer votre premier Azador",
				Description: "Combattez un Azador à la sortie du village.",
				Condition: func(player *createcharacter.Character) bool {
					return true // Validé manuellement dans checkAzadorKillQuest
				},
				Reward: QuestReward{
					Money: 75,
					Items: make(map[string]int),
				},
				Completed: false,
			},
			{
				ID:          3,
				Title:       "Récolter 2 pierres",
				Description: "Trouvez et récoltez 2 pierres sur la carte.",
				Condition: func(player *createcharacter.Character) bool {
					return player.Inventory.HasItem("Pierre", 2)
				},
				Reward: QuestReward{
					Money: 25,
					Items: make(map[string]int),
				},
				Completed: false,
			},
			{
				ID:          4,
				Title:       "Récolter 1 bâton",
				Description: "Trouvez et récoltez 1 bâton sur la carte.",
				Condition: func(player *createcharacter.Character) bool {
					return player.Inventory.HasItem("Bâton", 1)
				},
				Reward: QuestReward{
					Money: 25,
					Items: make(map[string]int),
				},
				Completed: false,
			},
			{
				ID:          5,
				Title:       "Crafter une lame rouillée",
				Description: "Appuyez sur [C] pour ouvrir le menu de crafting et craftez une lame rouillée.",
				Condition: func(player *createcharacter.Character) bool {
					return player.Inventory.HasItem("Lame rouillé", 2)
				},
				Reward: QuestReward{
					Money: 50,
					Items: make(map[string]int),
				},
				Completed: false,
			},
			{
				ID:          6,
				Title:       "Aller voir Valenric le forgeron",
				Description: "Trouvez Valenric le forgeron dans le village.",
				Condition: func(player *createcharacter.Character) bool {
					return true // Validé manuellement lors de l'interaction avec Valenric
				},
				Reward: QuestReward{
					Money: 75,
					Items: make(map[string]int),
				},
				Completed: false,
			},
			{
				ID:          7,
				Title:       "Upgrader votre arme chez Valenric",
				Description: "Utilisez les services de Valenric pour améliorer votre arme.",
				Condition: func(player *createcharacter.Character) bool {
					return player.Inventory.HasItem("épée de chevalier", 1)
				},
				Reward: QuestReward{
					Money: 200,
					Items: make(map[string]int),
				},
				Completed: false,
			},
		},
	}

	emeryn.AddQuest(introQuest)

	// Créer la quête principale
	mainQuest := Quest{
		ID:          "main_quest",
		Title:       "La menace des Azadors",
		Description: "Les Azadors deviennent de plus en plus agressifs. Aidez le village à faire face à cette menace.",
		CurrentStep: 0,
		Completed:   false,
		Steps: []QuestStep{
			{
				ID:          1,
				Title:       "Retourner voir Emeryn",
				Description: "Retournez voir Emeryn au centre du village pour votre prochaine mission.",
				Condition: func(player *createcharacter.Character) bool {
					return true // Validé manuellement lors de l'interaction
				},
				Reward: QuestReward{
					Money: 100,
					Items: make(map[string]int),
				},
				Completed: false,
			},
			{
				ID:          2,
				Title:       "Éliminer 3 Azadors",
				Description: "Tuez 3 Azadors pour protéger le village.",
				Condition: func(player *createcharacter.Character) bool {
					return true // Validé manuellement dans le système de combat
				},
				Reward: QuestReward{
					Money: 150,
					Items: make(map[string]int),
				},
				Completed: false,
			},
			{
				ID:          3,
				Title:       "Parler à la marchande Sarahlia",
				Description: "Allez voir Sarahlia qui a un problème avec des voleurs.",
				Condition: func(player *createcharacter.Character) bool {
					return true // Validé manuellement lors de l'interaction
				},
				Reward: QuestReward{
					Money: 75,
					Items: make(map[string]int),
				},
				Completed: false,
			},
			{
				ID:          4,
				Title:       "Récupérer la potion volée",
				Description: "Tuez 1 Azador pour récupérer la potion volée de Sarahlia.",
				Condition: func(player *createcharacter.Character) bool {
					return player.Inventory.HasItem("Heal potion", 1)
				},
				Reward: QuestReward{
					Money: 100,
					Items: make(map[string]int),
				},
				Completed: false,
			},
			{
				ID:          5,
				Title:       "Rapporter la potion à Sarahlia",
				Description: "Retournez voir Sarahlia avec la potion récupérée.",
				Condition: func(player *createcharacter.Character) bool {
					return true // Validé manuellement lors de l'interaction
				},
				Reward: QuestReward{
					Money: 100,
					Items: map[string]int{"Heal potion": 1},
				},
				Completed: false,
			},
			{
				ID:          6,
				Title:       "Retourner voir Emeryn",
				Description: "Retournez voir Emeryn une fois que vous avez atteint le niveau 3.",
				Condition: func(player *createcharacter.Character) bool {
					return player.Level >= 3 || player.Name == "God"
				},
				Reward: QuestReward{
					Money: 200,
					Items: make(map[string]int),
				},
				Completed: false,
			},
			{
				ID:          7,
				Title:       "Trouver et débloquer le portail",
				Description: "Trouvez le portail dans la forêt et appuyez sur [E] pour le débloquer.",
				Condition: func(player *createcharacter.Character) bool {
					return true // Validé manuellement lors de l'interaction avec le portail
				},
				Reward: QuestReward{
					Money: 300,
					Items: make(map[string]int),
				},
				Completed: false,
			},
		},
	}

	emeryn.AddQuest(mainQuest)

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
	// Vérifier d'abord la quête d'introduction
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

					message := "Je viens de te donner 1 lame rouillée.\n\nVa à la sortie du village, tu trouveras un Azador, bats-le !"
					if expMessage != "" {
						message += "\n\n" + expMessage
					}
					return message
				}
			} else if quest.CurrentStep == 1 {
				// Étape 2 : Tuer un Azador
				return "Trouve un Azador à la sortie du village et bats-le !"
			} else if quest.CurrentStep == 2 {
				// Étape 3 : Récolter 2 pierres
				return "Maintenant, récolte 2 pierres sur la carte !"
			} else if quest.CurrentStep == 3 {
				// Étape 4 : Récolter 1 bâton
				return "Maintenant, récolte 1 bâton sur la carte !"
			} else if quest.CurrentStep == 4 {
				// Étape 5 : Crafter une lame rouillée
				return "Maintenant, appuie sur [C] pour crafter une lame rouillée avec tes matériaux !"
			} else if quest.CurrentStep == 5 {
				// Étape 6 : Aller voir Valenric
				return "Va voir Valenric le forgeron dans le village !"
			} else if quest.CurrentStep == 6 {
				// Étape 7 : Upgrader l'arme
				return "Utilise les services de Valenric pour upgrader ton arme !"
			}
		}
	}

	// Vérifier la quête principale
	for _, quest := range npc.Quests {
		if quest.ID == "main_quest" && !quest.Completed {
			if quest.CurrentStep == 0 {
				// Démarrer automatiquement la quête principale
				npc.ValidateQuestStep(player, "main_quest")
				return "Excellent travail pour ton introduction ! Maintenant, les Azadors deviennent de plus en plus menaçants.\n\nJe vais avoir besoin de ton aide pour protéger le village. Es-tu prêt pour de vrais défis ?"
			} else if quest.CurrentStep == 1 {
				// Étape 2 : Éliminer 3 Azadors
				return "Va éliminer 3 Azadors pour protéger le village. Ils deviennent de plus en plus agressifs !"
			} else if quest.CurrentStep == 2 {
				// Étape 3 : Parler à Sarahlia
				return "Va voir Sarahlia la marchande, elle a des problèmes avec des voleurs !"
			} else if quest.CurrentStep == 3 {
				// Étape 4 : Récupérer la potion
				return "Tue un Azador pour récupérer la potion volée de Sarahlia !"
			} else if quest.CurrentStep == 4 {
				// Étape 5 : Rapporter à Sarahlia
				return "Retourne voir Sarahlia avec la potion récupérée !"
			} else if quest.CurrentStep == 5 {
				// Étape 6 : Retourner voir Emeryn (après niveau 3)
				if player.Name == "God" {
					npc.ValidateQuestStep(player, "main_quest")
					return "Bienvenue, divinité ! Tu n'as pas besoin d'entraînement. Tu es prêt pour la mission finale ! Il est temps de trouver le portail dans la forêt."
				} else if player.Level >= 3 {
					npc.ValidateQuestStep(player, "main_quest")
					return "Excellent ! Tu as atteint le niveau 3. Tu es maintenant prêt pour le combat final ! Il est temps de trouver le portail dans la forêt."
				} else {
					return fmt.Sprintf("Tu te débrouilles bien ! Continue à t'entraîner en tuant des Azadors jusqu'au niveau 3, puis reviens me voir. (Niveau actuel: %d)", player.Level)
				}
			} else if quest.CurrentStep == 6 {
				// Étape 7 : Trouver le portail
				return "Il est temps de trouver le portail dans la forêt. Approche-toi du portail et appuie sur [E] pour le débloquer !"
			}
		}
	}

	// Si toutes les quêtes sont terminées
	return "Tu es devenu un vrai héros ! Le village est en sécurité grâce à toi."
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