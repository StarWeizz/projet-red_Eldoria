package interactions

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	inventory "eldoria/Inventory"
	"eldoria/combat"
	"eldoria/items"
	"eldoria/money"
	"eldoria/npcs"
	createcharacter "eldoria/player"
	"eldoria/worlds"
)

type ShopItem struct {
	Item  items.Item
	Price int
}

type RespawnData struct {
	RespawnTime time.Time
	ObjectType  string // "rock", "monster", etc.
}

type InteractionManager struct {
	respawnQueue   map[string]RespawnData // Clé: x_y_world, Valeur: temps de respawn
	inventory      *inventory.Inventory
	playerMoney    *money.Money
	shopItems      []ShopItem
	emeryn         *npcs.NPC
	azadorsKilled  int  // Compteur pour la quête principale
	sarhaliaRobbed bool // État pour la quête de Sarahlia
}

func NewInteractionManager(inv *inventory.Inventory, playerMoney *money.Money) *InteractionManager {
	// Créer les objets de la boutique automatiquement à partir de tous les CraftingItems
	var shopItems []ShopItem
	// Ajout Heal potion depuis PotionsList
	if healPotion, exists := items.PotionsList["Heal potion"]; exists {
		shopItems = append(shopItems, ShopItem{
			Item:  healPotion,
			Price: healPotion.GetPrice(),
		})
	}

	itemOrder := []string{"Bâton", "Pierre", "Papier", "Parchemin", "Ecaille d'Azador"}

	for _, itemName := range itemOrder {
		if item, exists := items.CraftingItems[itemName]; exists {
			shopItems = append(shopItems, ShopItem{
				Item:  item,
				Price: item.GetPrice(),
			})
		}
	}

	return &InteractionManager{
		respawnQueue:   make(map[string]RespawnData),
		inventory:      inv,
		playerMoney:    playerMoney,
		shopItems:      shopItems,
		emeryn:         npcs.CreateEmeryn(),
		azadorsKilled:  0,
		sarhaliaRobbed: false,
	}
}

type InteractionResult struct {
	Success      bool
	Message      string
	ItemGained   items.Item
	ShouldRemove bool
	RespawnTime  time.Duration
	EndGame      bool
	UnlockPortal bool
}

func (im *InteractionManager) HandleInteraction(world *worlds.World, player *createcharacter.Character, x, y int, interactionType string, playerChoice func(h *createcharacter.Character, m *combat.Monster) string) *InteractionResult {
	switch interactionType {
	case "pickup":
		return im.handlePickup(world, player, x, y)
	case "chest":
		return im.handleChest(world, x, y)
	case "door":
		return im.handleDoor(world, x, y)
	case "treasure":
		return im.handleTreasure(world, x, y)
	case "merchant":
		return im.handleMerchant(world, x, y)
	case "blacksmith":
		return im.handleBlacksmith(world, player, x, y)
	case "emeryn":
		return im.handleEmeryn(player)
	case "portal":
		return im.handlePortal(world, player, x, y)
	case "boss":
		// Combat contre Maximor, boss
		maximor := &combat.Monster{
			Name:    "Maximor",
			HP:      175,
			Attack:  10,
			Defense: 8,
		}
		win, playerDamages, monsterDamages, fled := combat.StartCombat(player, maximor, playerChoice)
		var damageLog string
		for i := 0; i < len(playerDamages) || i < len(monsterDamages); i++ {
			turn := i + 1
			dmgPlayer := 0
			dmgMonster := 0
			special := false
			if i < len(playerDamages) {
				dmgPlayer = playerDamages[i]
			}
			if i < len(monsterDamages) {
				dmgMonster = monsterDamages[i]
				if dmgMonster < 0 {
					dmgMonster = -dmgMonster
					special = true
				}
			}
			if special {
				damageLog += fmt.Sprintf("Tour %d\n→ Vous infligez %d dégâts à %s.\n← %s utilise une ATTAQUE SPÉCIALE et vous inflige %d dégâts !\n\n", turn, dmgPlayer, maximor.Name, maximor.Name, dmgMonster)
			} else {
				damageLog += fmt.Sprintf("Tour %d\n→ Vous infligez %d dégâts à %s.\n← %s vous inflige %d dégâts.\n\n", turn, dmgPlayer, maximor.Name, maximor.Name, dmgMonster)
			}
		}

		if win {
			message := fmt.Sprintf("🏆 Vous avez vaincu Maximor !\n%s", damageLog)
			player.CurrentHP = 0
			return &InteractionResult{
				Success: true,
				Message: message,
				EndGame: true, // Peut déclencher la fin du jeu
			}
		} else if fled {
			message := fmt.Sprintf("🏃‍♂️ Vous avez fui le combat contre Maximor.\n%s", damageLog)
			return &InteractionResult{
				Success: false,
				Message: message,
				EndGame: false,
			}
		} else {
			message := fmt.Sprintf("💀 Vous avez été vaincu par Maximor...\n%s", damageLog)
			return &InteractionResult{
				Success: false,
				Message: message,
				EndGame: true,
			}
		}
	case "monster":
		monster := combat.NewRandomMonster()
		win, playerDamages, monsterDamages, fled := combat.StartCombat(player, monster, playerChoice)
		var damageLog string
		for i := 0; i < len(playerDamages) || i < len(monsterDamages); i++ {
			turn := i + 1
			dmgPlayer := 0
			dmgMonster := 0
			if i < len(playerDamages) {
				dmgPlayer = playerDamages[i]
			}
			if i < len(monsterDamages) {
				dmgMonster = monsterDamages[i]
			}
			damageLog += fmt.Sprintf("Tour %d\n→ Vous infligez %d dégâts à %s.\n← %s vous inflige %d dégâts.\n\n", turn, dmgPlayer, monster.Name, monster.Name, dmgMonster)
		}

		if win {
			// Donner de l'EXP selon le type de monstre
			var expGained int
			switch monster.Name {
			case "Apprenti Azador":
				expGained = 5
			case "Azador":
				expGained = 10
			case "Azador Chevalier":
				expGained = 20
			default:
				expGained = 5 // XP par défaut
			}

			expMessage := player.AddExperience(expGained)

			// Vérifier si c'est un Azador et si le joueur est à l'étape de quête appropriée
			questMessage := ""
			specialDrop := false
			if strings.Contains(monster.Name, "Azador") {
				// Vérifier si le joueur est à l'étape de récupération de potion
				if im.shouldDropPotionForQuest() {
					// Donner une potion de soin en plus
					if healPotion, exists := items.PotionsList["Heal potion"]; exists {
						im.inventory.Add(healPotion, 1)
						specialDrop = true
					}
				}
				questMessage = im.checkAzadorKillQuest(player)
			}

			if dropItem, exists := items.CraftingItems["Ecaille d'Azador"]; exists {
				im.inventory.Add(dropItem, 1)
				respawnKey := fmt.Sprintf("%d|%d|%s", x, y, world.Name)
				im.respawnQueue[respawnKey] = RespawnData{
					RespawnTime: time.Now().Add(20 * time.Second),
					ObjectType:  "monster",
				}

				message := fmt.Sprintf("🏆 Vous avez vaincu %s et obtenu %s !\n%s", monster.Name, dropItem.GetName(), damageLog)
				if specialDrop {
					message += " et Heal potion (potion volée récupérée) !"
				} else {
					message += " !"
				}
				if expMessage != "" {
					message += "\n" + expMessage
				}
				if questMessage != "" {
					message += "\n" + questMessage
				}

				return &InteractionResult{
					Success:      true,
					Message:      message,
					ItemGained:   dropItem,
					ShouldRemove: true,
					RespawnTime:  20 * time.Second,
				}
			}

			message := fmt.Sprintf("🏆 Vous avez vaincu %s !\n%s", monster.Name, damageLog)
			if expMessage != "" {
				message += "\n" + expMessage
			}
			if questMessage != "" {
				message += "\n" + questMessage
			}

			return &InteractionResult{
				Success: true,
				Message: message,
			}
		} else if fled {
			message := fmt.Sprintf("🏃‍♂️ Vous avez fui le combat contre %s.\n%s", monster.Name, damageLog)
			return &InteractionResult{
				Success: false,
				Message: message,
			}
		} else {
			message := fmt.Sprintf("💀 Vous avez été vaincu par %s...\n%s", monster.Name, damageLog)
			return &InteractionResult{
				Success: false,
				Message: message,
			}
		}
	default:
		return &InteractionResult{
			Success: false,
			Message: "Aucune interaction disponible.",
		}
	}
}

func (im *InteractionManager) handlePickup(world *worlds.World, player *createcharacter.Character, x, y int) *InteractionResult {
	objectType := world.GetObjectTypeAt(x, y)

	switch objectType {
	case "rock":
		// Créer un item pierre depuis la map des CraftingItems
		if stoneItem, exists := items.CraftingItems["Pierre"]; exists {
			im.inventory.Add(stoneItem, 1)

			// Donner de l'EXP pour la récolte
			expMessage := player.AddExperience(1)
			message := "🪨 Pierre ramassée !"
			if expMessage != "" {
				message += "\n" + expMessage
			}

			// Vérifier les progrès de quête
			questMessage := im.checkQuestProgress(player)
			if questMessage != "" {
				message += "\n" + questMessage
			}

			// Planifier le respawn dans 10 secondes
			// Utiliser un séparateur différent pour éviter les problèmes avec les espaces dans le nom du monde
			respawnKey := fmt.Sprintf("%d|%d|%s", x, y, world.Name)
			im.respawnQueue[respawnKey] = RespawnData{
				RespawnTime: time.Now().Add(10 * time.Second),
				ObjectType:  "rock",
			}

			return &InteractionResult{
				Success:      true,
				Message:      message,
				ItemGained:   stoneItem,
				ShouldRemove: true,
				RespawnTime:  10 * time.Second,
			}
		}

		return &InteractionResult{
			Success: false,
			Message: "Erreur lors de la récupération de l'item.",
		}

	case "stick":
		// Créer un item bâton depuis la map des CraftingItems
		if stickItem, exists := items.CraftingItems["Bâton"]; exists {
			im.inventory.Add(stickItem, 1)

			// Donner de l'EXP pour la récolte
			expMessage := player.AddExperience(1)
			message := "🪵 Bâton ramassé !"
			if expMessage != "" {
				message += "\n" + expMessage
			}

			// Vérifier les progrès de quête
			questMessage := im.checkQuestProgress(player)
			if questMessage != "" {
				message += "\n" + questMessage
			}

			// Planifier le respawn du bâton dans 15 secondes
			respawnKey := fmt.Sprintf("%d|%d|%s", x, y, world.Name)
			im.respawnQueue[respawnKey] = RespawnData{
				RespawnTime: time.Now().Add(15 * time.Second),
				ObjectType:  "stick",
			}

			return &InteractionResult{
				Success:      true,
				Message:      message,
				ItemGained:   stickItem,
				ShouldRemove: true,
				RespawnTime:  15 * time.Second,
			}
		}

		return &InteractionResult{
			Success: false,
			Message: "Erreur lors de la récupération de l'item.",
		}

	default:
		return &InteractionResult{
			Success: false,
			Message: "Rien à ramasser ici.",
		}
	}
}

func (im *InteractionManager) handleChest(world *worlds.World, x, y int) *InteractionResult {
	return &InteractionResult{
		Success: true,
		Message: "🎁 Coffre ouvert ! Vous avez trouvé une récompense.",
	}
}

func (im *InteractionManager) handleDoor(world *worlds.World, x, y int) *InteractionResult {
	// Messages de lore pour différentes portes selon leur position
	loreMessages := []string{
		"🚪 Vous entrez dans une demeure chaleureuse. Un parfum de soupe flotte dans l'air.",
		"🚪 Cette maison semble abandonnée depuis longtemps. Des toiles d'araignée ornent les coins.",
		"🚪 Vous pénétrez dans une forge. Le bruit du marteau résonne encore faiblement.",
		"🚪 Une bibliothèque poussiéreuse s'étend devant vous. Des grimoires anciens tapissent les étagères.",
		"🚪 L'intérieur de cette maison révèle un laboratoire d'alchimie mystérieux.",
	}

	// Utiliser la position pour déterminer quel message afficher
	messageIndex := (x + y) % len(loreMessages)

	return &InteractionResult{
		Success: true,
		Message: loreMessages[messageIndex],
	}
}

func (im *InteractionManager) handleTreasure(world *worlds.World, x, y int) *InteractionResult {
	return &InteractionResult{
		Success: true,
		Message: "💎 Trésor trouvé !",
	}
}

func (im *InteractionManager) handleMerchant(world *worlds.World, x, y int) *InteractionResult {
	// Cette fonction n'a pas accès au player, donc on ne peut pas valider les quêtes ici
	// Les quêtes de Sarahlia seront gérées via une interaction spéciale

	// Afficher la liste des objets disponibles
	shopMessage := "💎 Sarhalia : \"Bienvenue dans ma boutique !\"\n\nArticles disponibles :\n"
	for i, shopItem := range im.shopItems {
		shopMessage += fmt.Sprintf("%d. %s - %d 💰\n", i+1, shopItem.Item.GetName(), shopItem.Price)
	}

	return &InteractionResult{
		Success: true,
		Message: shopMessage,
	}
}

// Nouvelle méthode pour gérer l'achat d'un objet spécifique
func (im *InteractionManager) BuyItem(itemIndex int) *InteractionResult {
	if itemIndex < 0 || itemIndex >= len(im.shopItems) {
		return &InteractionResult{
			Success: false,
			Message: "❌ Objet invalide.",
		}
	}

	shopItem := im.shopItems[itemIndex]

	// Vérifier si le joueur a assez d'argent
	if im.playerMoney.Get() < shopItem.Price {
		return &InteractionResult{
			Success: false,
			Message: fmt.Sprintf("❌ Pas assez d'argent ! Il vous faut %d 💰 mais vous n'avez que %d 💰.", shopItem.Price, im.playerMoney.Get()),
		}
	}

	// Effectuer la transaction
	if im.playerMoney.Remove(shopItem.Price) {
		im.inventory.Add(shopItem.Item, 1)
		return &InteractionResult{
			Success:    true,
			Message:    fmt.Sprintf("✅ Vous avez acheté %s pour %d 💰 ! Il vous reste %d 💰.", shopItem.Item.GetName(), shopItem.Price, im.playerMoney.Get()),
			ItemGained: shopItem.Item,
		}
	}

	return &InteractionResult{
		Success: false,
		Message: "❌ Erreur lors de l'achat.",
	}
}

func (im *InteractionManager) handleBlacksmith(world *worlds.World, player *createcharacter.Character, x, y int) *InteractionResult {
	// Vérifier d'abord si c'est la première visite pour la quête
	questMessage := im.checkBlacksmithQuestProgress(player)

	// Générer la liste des armes upgrade possibles
	upgradeOptions := im.getWeaponUpgradeOptions()

	message := "⚒️ Valenric : \"Salut aventurier ! Je peux upgrader tes armes !\"\n\n"

	if questMessage != "" {
		message = questMessage + "\n\n" + message
	}

	if len(upgradeOptions) == 0 {
		message += "Tu n'as pas d'armes upgradables pour le moment.\n"
		message += "Apporte-moi 2 armes du même niveau et je les combinerai en une arme plus puissante !"
	} else {
		message += "Armes upgradables :\n"
		for i, option := range upgradeOptions {
			// Afficher la quantité disponible
			availableQuantity := im.inventory.Items[option.Current.GetName()]
			message += fmt.Sprintf("%d. %s (2 sur %d disponibles) → %s\n", i+1, option.Current.GetName(), availableQuantity, option.Next.GetName())
		}
		message += "\nAppuyez sur [1-%d] pour upgrader l'arme correspondante."
		message = fmt.Sprintf(message, len(upgradeOptions))
	}

	return &InteractionResult{
		Success: true,
		Message: message,
	}
}

func (im *InteractionManager) handleEmeryn(player *createcharacter.Character) *InteractionResult {
	// Utiliser le nouveau système de messages d'Emeryn
	message := im.emeryn.GetEmerynMessage(player)

	return &InteractionResult{
		Success: true,
		Message: message,
	}
}

func (im *InteractionManager) handlePortal(world *worlds.World, player *createcharacter.Character, x, y int) *InteractionResult {
	// Vérifier si c'est une interaction pour débloquer le portail
	if im.emeryn != nil {
		for _, quest := range im.emeryn.Quests {
			if quest.ID == "main_quest" {
				if !quest.Completed && quest.CurrentStep == 6 {
					// Étape 6 (index 6) : Trouver le portail et le débloquer
					if im.emeryn.ValidateQuestStep(player, "main_quest") {
						return &InteractionResult{
							Success:      true,
							Message:      "🌟 Bravo ! Vous avez trouvé le portail et l'avez débloqué ! Quête principale terminée !\n\nVous pouvez maintenant utiliser [TAB] pour changer de monde !",
							UnlockPortal: true,
						}
					}
				} else if quest.Completed {
					// Quête terminée, portail déjà débloqué
					return &InteractionResult{
						Success: true,
						Message: "🌀 Portail vers Eldoria. Utilisez [TAB] pour changer de monde !",
					}
				}
			}
		}
	}

	return &InteractionResult{
		Success: true,
		Message: "🌀 Portail mystérieux vers Eldoria. Vous devez d'abord terminer votre quête pour pouvoir l'activer.",
	}
}

// AdvanceEmerynInteraction fait avancer l'interaction avec Emeryn (appelé lors de l'appui sur espace)
func (im *InteractionManager) AdvanceEmerynInteraction() {
	im.emeryn.AdvanceEmerynPhase()
}

// CanAdvanceEmerynInteraction vérifie si on peut faire avancer l'interaction avec Emeryn
func (im *InteractionManager) CanAdvanceEmerynInteraction() bool {
	if im.emeryn == nil {
		return false
	}
	return im.emeryn.CanAdvanceEmeryn()
}

// GetEmerynQuests retourne les quêtes d'Emeryn
func (im *InteractionManager) GetEmerynQuests() []npcs.Quest {
	if im.emeryn == nil {
		return []npcs.Quest{}
	}
	return im.emeryn.Quests
}

// --- Respawn ---
func (im *InteractionManager) CheckRespawns(world *worlds.World) []string {
	var messages []string
	now := time.Now()

	for respawnKey, data := range im.respawnQueue {
		if now.After(data.RespawnTime) {
			parts := strings.Split(respawnKey, "|")
			if len(parts) != 3 {
				delete(im.respawnQueue, respawnKey)
				continue
			}

			x, _ := strconv.Atoi(parts[0])
			y, _ := strconv.Atoi(parts[1])
			worldName := parts[2]

			if worldName == world.Name {
				_ = world.RespawnObject(x, y, data.ObjectType)
				// Messages de respawn supprimés pour une expérience plus fluide
			}
			delete(im.respawnQueue, respawnKey)
		}
	}
	return messages
}

// Méthode pour vérifier si le joueur peut interagir avec un objet adjacent
func (im *InteractionManager) CheckNearbyInteractions(world *worlds.World) []string {
	var availableInteractions []string

	// Vérifier les 4 directions autour du joueur
	coords := [][2]int{
		{world.PlayerX + 1, world.PlayerY},
		{world.PlayerX - 1, world.PlayerY},
		{world.PlayerX, world.PlayerY + 1},
		{world.PlayerX, world.PlayerY - 1},
	}

	for _, coord := range coords {
		x, y := coord[0], coord[1]
		if x >= 0 && x < world.Width && y >= 0 && y < world.Height {
			interactionType := world.GetInteractionType(x, y)
			if interactionType != "none" && interactionType != "" && interactionType != "door" {
				objectName := world.GetObjectNameAt(x, y)
				availableInteractions = append(availableInteractions,
					fmt.Sprintf("Appuyez sur [E] près de %s pour %s", objectName, interactionType))
			}
		}
	}

	return availableInteractions
}

// checkAzadorKillQuest vérifie si tuer un Azador fait progresser les quêtes
func (im *InteractionManager) checkAzadorKillQuest(player *createcharacter.Character) string {
	if im.emeryn == nil {
		return ""
	}

	// Vérifier les quêtes d'Emeryn
	for _, quest := range im.emeryn.Quests {
		// Quête d'introduction
		if quest.ID == "intro_quest" && !quest.Completed {
			if quest.CurrentStep == 1 {
				if im.emeryn.ValidateQuestStep(player, "intro_quest") {
					return "✨ Quête mise à jour : Récolte maintenant 2 pierres !"
				}
			}
		}
		// Quête principale
		if quest.ID == "main_quest" && !quest.Completed {
			if quest.CurrentStep == 1 {
				// Compter les Azadors tués pour la quête principale
				im.azadorsKilled++
				if im.azadorsKilled >= 3 {
					if im.emeryn.ValidateQuestStep(player, "main_quest") {
						return "✨ Quête mise à jour : Va voir Sarahlia la marchande !"
					}
				} else {
					return fmt.Sprintf("✨ Azadors éliminés : %d/3", im.azadorsKilled)
				}
			} else if quest.CurrentStep == 3 {
				// Récupérer la potion volée
				if im.emeryn.ValidateQuestStep(player, "main_quest") {
					return "✨ Quête mise à jour : Retourne voir Sarahlia avec la potion !"
				}
			}
		}
	}

	return ""
}

// shouldDropPotionForQuest vérifie si un Azador doit drop une potion pour la quête
func (im *InteractionManager) shouldDropPotionForQuest() bool {
	if im.emeryn == nil {
		return false
	}

	// Vérifier si le joueur est à l'étape 3 de la quête principale (récupérer la potion)
	for _, quest := range im.emeryn.Quests {
		if quest.ID == "main_quest" && !quest.Completed && quest.CurrentStep == 3 {
			return true
		}
	}

	return false
}

// checkQuestProgress vérifie et fait progresser automatiquement les quêtes selon l'inventaire
func (im *InteractionManager) checkQuestProgress(player *createcharacter.Character) string {
	if im.emeryn == nil {
		return ""
	}

	// Vérifier la quête d'introduction d'Emeryn
	for _, quest := range im.emeryn.Quests {
		if quest.ID == "intro_quest" && !quest.Completed {
			// Étape 2 (index 2) : Récolter 2 pierres
			if quest.CurrentStep == 2 {
				if player.Inventory.HasItem("Pierre", 2) {
					if im.emeryn.ValidateQuestStep(player, "intro_quest") {
						return "✨ Quête mise à jour : Récolte maintenant 1 bâton !"
					}
				}
			}
			// Étape 3 (index 3) : Récolter 1 bâton
			if quest.CurrentStep == 3 {
				if player.Inventory.HasItem("Bâton", 1) {
					if im.emeryn.ValidateQuestStep(player, "intro_quest") {
						return "✨ Quête mise à jour : Appuie sur [C] pour crafter une lame rouillée !"
					}
				}
			}
			// Étape 4 (index 4) : Crafter une lame rouillée
			if quest.CurrentStep == 4 {
				if player.Inventory.HasItem("Lame rouillé", 2) {
					if im.emeryn.ValidateQuestStep(player, "intro_quest") {
						return "✨ Quête mise à jour : Va voir Valenric le forgeron !"
					}
				}
			}
			break
		}
	}

	// Vérifier la quête principale d'Emeryn
	for _, quest := range im.emeryn.Quests {
		if quest.ID == "main_quest" && !quest.Completed {
			// Étape 5 (index 5) : Atteindre niveau 3 puis retourner voir Emeryn
			if quest.CurrentStep == 5 {
				if player.Level >= 3 || player.Name == "God" {
					// Ne pas faire avancer automatiquement, le joueur doit retourner voir Emeryn
					return "✨ Objectif accompli ! Retourne voir Emeryn pour la suite de ta mission !"
				}
			}
			break
		}
	}

	return ""
}

// checkBlacksmithQuestProgress vérifie si le joueur visite Valenric pour la première fois dans la quête
func (im *InteractionManager) checkBlacksmithQuestProgress(player *createcharacter.Character) string {
	if im.emeryn == nil {
		return ""
	}

	// Vérifier la quête d'introduction d'Emeryn
	for _, quest := range im.emeryn.Quests {
		if quest.ID == "intro_quest" && !quest.Completed {
			// Étape 5 (index 5) : Aller voir Valenric
			if quest.CurrentStep == 5 {
				if im.emeryn.ValidateQuestStep(player, "intro_quest") {
					return "✨ Quête mise à jour : Upgrade maintenant ton arme !"
				}
			}
			break
		}
	}

	return ""
}

// checkUpgradeQuestProgress vérifie si l'upgrade d'arme termine la quête d'introduction
func (im *InteractionManager) checkUpgradeQuestProgress(player *createcharacter.Character) string {
	if im.emeryn == nil {
		return ""
	}

	// Vérifier la quête d'introduction d'Emeryn
	for _, quest := range im.emeryn.Quests {
		if quest.ID == "intro_quest" && !quest.Completed {
			// Étape 6 (index 6) : Upgrader l'arme
			if quest.CurrentStep == 6 {
				// Vérifier si le joueur a une épée de chevalier (arme upgradée)
				if im.inventory.Items["épée de chevalier"] >= 1 {
					if im.emeryn.ValidateQuestStep(player, "intro_quest") {
						return "🎉✨ QUÊTE TERMINÉE ! Félicitations, vous avez complété votre introduction à l'aventure !"
					}
				}
			}
			break
		}
	}

	return ""
}

// checkSarhaliaQuestProgress gère les interactions spéciales avec Sarahlia selon la quête
func (im *InteractionManager) checkSarhaliaQuestProgress() string {
	if im.emeryn == nil {
		return ""
	}

	// Vérifier la quête principale d'Emeryn
	for _, quest := range im.emeryn.Quests {
		if quest.ID == "main_quest" && !quest.Completed {
			// Étape 2 (index 2) : Première visite chez Sarahlia
			if quest.CurrentStep == 2 {
				im.sarhaliaRobbed = true
				if im.emeryn.ValidateQuestStep(nil, "main_quest") {
					return "💎 Sarahlia : \"Oh non ! Des Azadors ont volé mes précieuses potions de soin !\n\nPeux-tu m'aider à en récupérer au moins une ? Je te récompenserai généreusement !\""
				}
			}
			// Étape 4 (index 4) : Rapporter la potion récupérée
			if quest.CurrentStep == 4 {
				if im.inventory.Items["Heal potion"] >= 1 {
					// Retirer la potion de l'inventaire (Sarahlia la récupère)
					if healPotion, exists := items.PotionsList["Heal potion"]; exists {
						im.inventory.Remove(healPotion, 1)
					}
					if im.emeryn.ValidateQuestStep(nil, "main_quest") {
						// Donner une potion bonus comme récompense
						if healPotion, exists := items.PotionsList["Heal potion"]; exists {
							im.inventory.Add(healPotion, 1)
						}
						return "💎 Sarahlia : \"Merci infiniment ! Tu as récupéré ma potion !\n\nVoici une potion supplémentaire en remerciement. Tu es un vrai héros !\""
					}
				} else {
					return "💎 Sarahlia : \"As-tu récupéré ma potion volée ? Je ne la vois pas dans ton inventaire...\""
				}
			}
		}
	}

	return ""
}

// CheckQuestProgressPublic est une méthode publique pour vérifier les progrès de quête depuis l'extérieur
func (im *InteractionManager) CheckQuestProgressPublic(player *createcharacter.Character) string {
	return im.checkQuestProgress(player)
}

// CheckSarhaliaQuestPublic gère les interactions spéciales avec Sarahlia
func (im *InteractionManager) CheckSarhaliaQuestPublic(player *createcharacter.Character) string {
	if im.emeryn == nil {
		return ""
	}

	// Vérifier la quête principale d'Emeryn
	for _, quest := range im.emeryn.Quests {
		if quest.ID == "main_quest" && !quest.Completed {
			// Étape 2 (index 2) : Première visite chez Sarahlia
			if quest.CurrentStep == 2 {
				im.sarhaliaRobbed = true
				if im.emeryn.ValidateQuestStep(player, "main_quest") {
					return "💎 Sarahlia : \"Oh non ! Des Azadors ont volé mes précieuses potions de soin !\n\nPeux-tu m'aider à en récupérer au moins une ? Je te récompenserai généreusement !\""
				}
			}
			// Étape 4 (index 4) : Rapporter la potion récupérée
			if quest.CurrentStep == 4 {
				if im.inventory.Items["Heal potion"] >= 1 {
					// Retirer la potion de l'inventaire (Sarahlia la récupère)
					if healPotion, exists := items.PotionsList["Heal potion"]; exists {
						im.inventory.Remove(healPotion, 1)
					}
					if im.emeryn.ValidateQuestStep(player, "main_quest") {
						// Donner une potion bonus comme récompense
						if healPotion, exists := items.PotionsList["Heal potion"]; exists {
							im.inventory.Add(healPotion, 1)
						}
						return "💎 Sarahlia : \"Merci infiniment ! Tu as récupéré ma potion !\n\nVoici une potion supplémentaire en remerciement. Tu es un vrai héros !\""
					}
				} else {
					return "💎 Sarahlia : \"As-tu récupéré ma potion volée ? Je ne la vois pas dans ton inventaire...\""
				}
			}
		}
	}

	return ""
}

// UpgradeOption représente une option d'upgrade d'arme
type UpgradeOption struct {
	Current items.Item
	Next    items.Item
}

// getWeaponUpgradeOptions retourne les options d'upgrade disponibles
func (im *InteractionManager) getWeaponUpgradeOptions() []UpgradeOption {
	var options []UpgradeOption

	// Mapper les upgrades d'armes
	weaponUpgrades := map[string]string{
		"Lame rouillé":       "épée de chevalier",
		"épée de chevalier":  "Epée Démoniaque",
		"Grimoire":           "Livre de Magie",
		"Livre de Magie":     "Livre des Ombre",
		"Couteaux de Chasse": "épée court runique",
		"épée court runique": "Dague de l'ombre",
	}

	// Vérifier chaque arme upgradable dans l'inventaire
	for currentName, nextName := range weaponUpgrades {
		// Vérifier si le joueur a au moins 2 de l'arme actuelle
		if quantity, exists := im.inventory.Items[currentName]; exists && quantity >= 2 {
			// Vérifier si l'arme suivante existe
			if currentWeapon, currentExists := items.WeaponList[currentName]; currentExists {
				if nextWeapon, nextExists := items.WeaponList[nextName]; nextExists {
					options = append(options, UpgradeOption{
						Current: currentWeapon,
						Next:    nextWeapon,
					})
				}
			}
		}
	}

	return options
}

// PerformWeaponUpgrade effectue l'upgrade d'une arme
func (im *InteractionManager) PerformWeaponUpgrade(player *createcharacter.Character, optionIndex int) *InteractionResult {
	options := im.getWeaponUpgradeOptions()

	if optionIndex < 0 || optionIndex >= len(options) {
		return &InteractionResult{
			Success: false,
			Message: "❌ Option d'upgrade invalide.",
		}
	}

	option := options[optionIndex]

	// Retirer 2 armes actuelles
	if !im.inventory.Remove(option.Current, 2) {
		return &InteractionResult{
			Success: false,
			Message: "❌ Pas assez d'armes à upgrader.",
		}
	}

	// Ajouter l'arme upgradée
	im.inventory.Add(option.Next, 1)

	// Calculer combien d'armes de base il reste
	remainingCount := im.inventory.Items[option.Current.GetName()]
	message := fmt.Sprintf("⚒️ Upgrade réussie ! Vous avez obtenu %s !", option.Next.GetName())
	if remainingCount > 0 {
		message += fmt.Sprintf(" Il vous reste %d x %s.", remainingCount, option.Current.GetName())
	}

	// Vérifier si cela termine la quête d'introduction
	questMessage := im.checkUpgradeQuestProgress(player)
	if questMessage != "" {
		message += "\n" + questMessage
	}

	return &InteractionResult{
		Success: true,
		Message: message,
	}
}

// GetEmeryn retourne la référence vers Emeryn pour les vérifications de quête externes
func (im *InteractionManager) GetEmeryn() *npcs.NPC {
	return im.emeryn
}
