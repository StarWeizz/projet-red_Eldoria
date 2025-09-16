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

type InteractionManager struct {
	respawnQueue map[string]time.Time // Clé: x_y_world, Valeur: temps de respawn
	inventory    *inventory.Inventory
	playerMoney  *money.Money
	shopItems    []ShopItem
	emeryn       *npcs.NPC
}

func NewInteractionManager(inv *inventory.Inventory, playerMoney *money.Money) *InteractionManager {
	// Créer les objets de la boutique automatiquement à partir de tous les CraftingItems
	var shopItems []ShopItem
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
		respawnQueue: make(map[string]time.Time),
		inventory:    inv,
		playerMoney:  playerMoney,
		shopItems:    shopItems,
		emeryn:       npcs.CreateEmeryn(),
	}
}

type InteractionResult struct {
	Success      bool
	Message      string
	ItemGained   items.Item
	ShouldRemove bool
	RespawnTime  time.Duration
}

func (im *InteractionManager) HandleInteraction(world *worlds.World, player *createcharacter.Character, x, y int, interactionType string) *InteractionResult {
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
		return im.handleBlacksmith(world, x, y)
	case "emeryn":
		return im.handleEmeryn(player)
	case "monster":
		monster := combat.GetRandomMonster()
		combat.StartCombat(player, monster) // <-- passer l'instance du joueur
		if player.CurrentHP > 0 {
			return &InteractionResult{
				Success: true,
				Message: fmt.Sprintf("Vous avez vaincu le %s ! 🏆", monster.Name),
			}
		} else {
			return &InteractionResult{
				Success: false,
				Message: fmt.Sprintf("Vous avez été vaincu par le %s... 💀", monster.Name),
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

			// Planifier le respawn dans 10 secondes
			// Utiliser un séparateur différent pour éviter les problèmes avec les espaces dans le nom du monde
			respawnKey := fmt.Sprintf("%d|%d|%s", x, y, world.Name)
			im.respawnQueue[respawnKey] = time.Now().Add(10 * time.Second)

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

func (im *InteractionManager) handleBlacksmith(world *worlds.World, x, y int) *InteractionResult {
	return &InteractionResult{
		Success: true,
		Message: "⚒️ Forgeron : \"Salut aventurier ! Je peux forger des armes et armures pour toi. Tu as des matériaux à transformer ?\"",
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

// Méthode pour vérifier et respawner les objets
func (im *InteractionManager) CheckRespawns(world *worlds.World) []string {
	var messages []string
	now := time.Now()

	for respawnKey, respawnTime := range im.respawnQueue {
		if now.After(respawnTime) {
			// Parser la clé pour récupérer x, y et le nom du monde
			parts := strings.Split(respawnKey, "|")
			if len(parts) != 3 {
				messages = append(messages, fmt.Sprintf("❌ Format respawn key invalide: %s", respawnKey))
				delete(im.respawnQueue, respawnKey)
				continue
			}

			x, err1 := strconv.Atoi(parts[0])
			y, err2 := strconv.Atoi(parts[1])
			worldName := parts[2]

			if err1 != nil || err2 != nil {
				messages = append(messages, fmt.Sprintf("❌ Erreur conversion coordonnées: %s", respawnKey))
				delete(im.respawnQueue, respawnKey)
				continue
			}

			if worldName == world.Name {
				// Faire réapparaître l'objet
				err := world.RespawnObject(x, y, "rock")
				if err != nil {
					messages = append(messages, fmt.Sprintf("❌ Erreur respawn du rocher: %v", err))
				}
			}

			// Supprimer de la queue
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
