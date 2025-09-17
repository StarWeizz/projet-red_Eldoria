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
	respawnQueue map[string]RespawnData // Clé: x|y|worldName
	inventory    *inventory.Inventory
	playerMoney  *money.Money
	shopItems    []ShopItem
}

func NewInteractionManager(inv *inventory.Inventory, playerMoney *money.Money) *InteractionManager {
	// Créer les objets de la boutique automatiquement
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
		respawnQueue: make(map[string]RespawnData),
		inventory:    inv,
		playerMoney:  playerMoney,
		shopItems:    shopItems,
	}
}

type InteractionResult struct {
	Success      bool
	Message      string
	ItemGained   items.Item
	ShouldRemove bool
	RespawnTime  time.Duration
	EndGame      bool // <- nouveau champ
}

// --- Gestion des interactions ---
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
	case "boss":
		boss := combat.NewMaximor()
		win := combat.StartCombat(player, &boss.Monster)
		if win {
			return &InteractionResult{
				Success: true,
				Message: fmt.Sprintf("🏆 Vous avez vaincu %s ! Le royaume est sauvé !", boss.Monster.Name),
				EndGame: true, // <- signal fin du jeu
			}
		} else {
			return &InteractionResult{
				Success: false,
				Message: fmt.Sprintf("💀 %s vous a vaincu... Le royaume est perdu...", boss.Monster.Name),
			}
		}
	case "monster":
		monster := combat.NewRandomMonster()
		win := combat.StartCombat(player, monster)
		if win {
			if dropItem, exists := items.CraftingItems["Ecaille d'Azador"]; exists {
				im.inventory.Add(dropItem, 1)
				respawnKey := fmt.Sprintf("%d|%d|%s", x, y, world.Name)
				im.respawnQueue[respawnKey] = RespawnData{
					RespawnTime: time.Now().Add(20 * time.Second),
					ObjectType:  "monster",
				}
				return &InteractionResult{
					Success:      true,
					Message:      fmt.Sprintf("🏆 Vous avez vaincu %s et obtenu %s !", monster.Name, dropItem.GetName()),
					ItemGained:   dropItem,
					ShouldRemove: true,
					RespawnTime:  20 * time.Second,
				}
			}
			return &InteractionResult{
				Success: true,
				Message: fmt.Sprintf("🏆 Vous avez vaincu %s !", monster.Name),
			}
		} else {
			return &InteractionResult{
				Success: false,
				Message: fmt.Sprintf("💀 Vous avez été vaincu par %s...", monster.Name),
			}
		}
	default:
		return &InteractionResult{
			Success: false,
			Message: "Aucune interaction disponible.",
		}
	}
}

// --- Pickup ---
func (im *InteractionManager) handlePickup(world *worlds.World, player *createcharacter.Character, x, y int) *InteractionResult {
	objectType := world.GetObjectTypeAt(x, y)
	switch objectType {
	case "rock":
		if stoneItem, exists := items.CraftingItems["Pierre"]; exists {
			im.inventory.Add(stoneItem, 1)
			respawnKey := fmt.Sprintf("%d|%d|%s", x, y, world.Name)
			im.respawnQueue[respawnKey] = RespawnData{
				RespawnTime: time.Now().Add(10 * time.Second),
				ObjectType:  "rock",
			}
			return &InteractionResult{
				Success:      true,
				Message:      fmt.Sprintf("🪨 Pierre ramassée !"),
				ItemGained:   stoneItem,
				ShouldRemove: true,
				RespawnTime:  10 * time.Second,
			}
		}
		return &InteractionResult{Success: false, Message: "Erreur lors de la récupération de l'item."}
	default:
		return &InteractionResult{Success: false, Message: "Rien à ramasser ici."}
	}
}

// --- Autres handlers ---
func (im *InteractionManager) handleChest(world *worlds.World, x, y int) *InteractionResult {
	return &InteractionResult{Success: true, Message: "🎁 Coffre ouvert ! Vous avez trouvé une récompense."}
}

func (im *InteractionManager) handleDoor(world *worlds.World, x, y int) *InteractionResult {
	loreMessages := []string{
		"🚪 Vous entrez dans une demeure chaleureuse.",
		"🚪 Cette maison semble abandonnée depuis longtemps.",
		"🚪 Vous pénétrez dans une forge.",
		"🚪 Une bibliothèque poussiéreuse s'étend devant vous.",
		"🚪 L'intérieur de cette maison révèle un laboratoire d'alchimie mystérieux.",
	}
	messageIndex := (x + y) % len(loreMessages)
	return &InteractionResult{Success: true, Message: loreMessages[messageIndex]}
}

func (im *InteractionManager) handleTreasure(world *worlds.World, x, y int) *InteractionResult {
	return &InteractionResult{Success: true, Message: "💎 Trésor trouvé !"}
}

func (im *InteractionManager) handleMerchant(world *worlds.World, x, y int) *InteractionResult {
	shopMessage := "🏪 Marchand : \"Bienvenue !\"\nObjets disponibles :\n"
	for i, shopItem := range im.shopItems {
		shopMessage += fmt.Sprintf("%d. %s - %d 💰\n", i+1, shopItem.Item.GetName(), shopItem.Price)
	}
	shopMessage += fmt.Sprintf("\nVotre argent : %d 💰\nAppuyez sur [1-5] pour acheter ou [Échap] pour quitter.", im.playerMoney.Get())
	return &InteractionResult{Success: true, Message: shopMessage}
}

func (im *InteractionManager) BuyItem(itemIndex int) *InteractionResult {
	if itemIndex < 0 || itemIndex >= len(im.shopItems) {
		return &InteractionResult{Success: false, Message: "❌ Objet invalide."}
	}
	shopItem := im.shopItems[itemIndex]
	if im.playerMoney.Get() < shopItem.Price {
		return &InteractionResult{Success: false, Message: "❌ Pas assez d'argent !"}
	}
	if im.playerMoney.Remove(shopItem.Price) {
		im.inventory.Add(shopItem.Item, 1)
		return &InteractionResult{Success: true, Message: fmt.Sprintf("✅ Vous avez acheté %s !", shopItem.Item.GetName()), ItemGained: shopItem.Item}
	}
	return &InteractionResult{Success: false, Message: "❌ Erreur lors de l'achat."}
}

func (im *InteractionManager) handleBlacksmith(world *worlds.World, x, y int) *InteractionResult {
	return &InteractionResult{Success: true, Message: "⚒️ Forgeron : \"Salut aventurier !\""}
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
				if data.ObjectType == "rock" {
					messages = append(messages, fmt.Sprintf("🪨 Un rocher a réapparu en (%d, %d)", x, y))
				} else if data.ObjectType == "monster" {
					messages = append(messages, fmt.Sprintf("👹 Un monstre a réapparu en (%d, %d)", x, y))
				}
			}
			delete(im.respawnQueue, respawnKey)
		}
	}
	return messages
}

// --- Interactions proches ---
func (im *InteractionManager) CheckNearbyInteractions(world *worlds.World) []string {
	var available []string
	coords := [][2]int{
		{world.PlayerX + 1, world.PlayerY},
		{world.PlayerX - 1, world.PlayerY},
		{world.PlayerX, world.PlayerY + 1},
		{world.PlayerX, world.PlayerY - 1},
	}

	for _, coord := range coords {
		x, y := coord[0], coord[1]
		if x >= 0 && x < world.Width && y >= 0 && y < world.Height {
			it := world.GetInteractionType(x, y)
			if it != "none" && it != "" {
				name := world.GetObjectNameAt(x, y)
				available = append(available, fmt.Sprintf("Appuyez sur [E] près de %s pour %s", name, it))
			}
		}
	}
	return available
}
