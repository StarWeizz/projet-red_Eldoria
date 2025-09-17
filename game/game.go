package game

import (
	"fmt"
	"strings"
	"time"

	inventory "eldoria/Inventory"
	"eldoria/interactions"
	"eldoria/money"
	createcharacter "eldoria/player"
	"eldoria/worlds"

	"github.com/gdamore/tcell/v2"
)

// GameState représente l'état du jeu
type GameState struct {
	Screen             tcell.Screen
	WorldList          []*worlds.World
	CurrentWorld       int
	PlayerCharacter    *createcharacter.Character
	PlayerMoney        *money.Money
	PlayerInventory    *inventory.Inventory
	InteractionManager *interactions.InteractionManager
	LoreMessage        string
	ShowingInventory   bool
}

// NewGameState crée un nouvel état de jeu
func NewGameState(screen tcell.Screen, playerCharacter *createcharacter.Character) *GameState {
	playerMoney := &playerCharacter.Gold
	playerInventory := playerCharacter.Inventory
	interactionManager := interactions.NewInteractionManager(playerInventory, playerMoney)

	return &GameState{
		Screen:             screen,
		WorldList:          []*worlds.World{},
		CurrentWorld:       0,
		PlayerCharacter:    playerCharacter,
		PlayerMoney:        playerMoney,
		PlayerInventory:    playerInventory,
		InteractionManager: interactionManager,
		LoreMessage:        "",
		ShowingInventory:   false,
	}
}

// GetInventoryCount compte les items dans l'inventaire
func (gs *GameState) GetInventoryCount() int {
	count := 0
	for _, qty := range gs.PlayerInventory.Items {
		count += qty
	}
	return count
}

// LoadWorlds charge les mondes depuis les fichiers de configuration
func (gs *GameState) LoadWorlds() {
	// Charger le monde Ynovia depuis JSON
	ynoviaConfig, err := worlds.LoadWorldConfig("configs/ynovia.json")
	if err != nil {
		gs.WorldList = append(gs.WorldList, worlds.NewGrid("Monde Ynovia", 80, 35, 10, 10))
	} else {
		gs.WorldList = append(gs.WorldList, worlds.NewWorldFromConfig(ynoviaConfig))
	}

	// Charger le monde Eldoria depuis JSON
	eldoriaConfig, err := worlds.LoadWorldConfig("configs/eldoria.json")
	if err != nil {
		gs.WorldList = append(gs.WorldList, worlds.NewGrid("Monde Eldoria", 40, 25, 5, 5))
	} else {
		gs.WorldList = append(gs.WorldList, worlds.NewWorldFromConfig(eldoriaConfig))
	}
}

// InitializePlayer place le joueur dans le monde initial
func (gs *GameState) InitializePlayer() {
	if len(gs.WorldList) > 0 {
		world := gs.WorldList[gs.CurrentWorld]
		world.OriginalTile = world.Grid[world.PlayerY][world.PlayerX]
		world.Grid[world.PlayerY][world.PlayerX] = '😀'
	}
}

// Draw affiche l'état du jeu à l'écran
func (gs *GameState) Draw() {
	gs.Screen.Clear()
	w := gs.WorldList[gs.CurrentWorld]
	screenWidth, screenHeight := gs.Screen.Size()

	// Topbar
	hiddenStatus := ""
	if w.IsPlayerHidden() {
		hiddenStatus = " - 🌿 CACHÉ des monstres"
	}
	inventoryCount := gs.GetInventoryCount()
	topbar := fmt.Sprintf("%s (%s) - %d/%d ♥ - 💰 %d - 🎒 %d items - %s - X:%d Y:%d%s",
		gs.PlayerCharacter.Name, gs.PlayerCharacter.Class,
		gs.PlayerCharacter.CurrentHP, gs.PlayerCharacter.MaxHP,
		gs.PlayerMoney.Get(), inventoryCount, w.Name, w.PlayerX, w.PlayerY, hiddenStatus)

	for i, r := range topbar {
		if i < screenWidth {
			gs.Screen.SetContent(i, 0, r, nil, tcell.StyleDefault)
		}
	}

	// Grille
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			r := w.Grid[y][x]
			gs.Screen.SetContent(x*2, y+1, r, nil, tcell.StyleDefault)
			gs.Screen.SetContent(x*2+1, y+1, ' ', nil, tcell.StyleDefault)
		}
	}

	// Dessiner le joueur à sa position
	gs.Screen.SetContent(w.PlayerX*2, w.PlayerY+1, gs.PlayerCharacter.Icon, nil, tcell.StyleDefault)
	gs.Screen.SetContent(w.PlayerX*2+1, w.PlayerY+1, ' ', nil, tcell.StyleDefault)

	// Bottombar - Afficher les interactions disponibles
	availableInteractions := gs.InteractionManager.CheckNearbyInteractions(w)
	bottomY := screenHeight - 1

	// Zone de lore/inventaire - Afficher sous la grille
	loreY := w.Height + 2 // Juste sous la grille avec une ligne d'espace

	if gs.ShowingInventory {
		// Afficher l'inventaire
		inventoryText := gs.PlayerInventory.GetInventoryString()
		lines := strings.Split(inventoryText, "\n")
		for i, line := range lines {
			if i < 10 && loreY+i < bottomY { // Limiter à 10 lignes et ne pas dépasser le bottombar
				for j, r := range line {
					if j < screenWidth {
						gs.Screen.SetContent(j, loreY+i, r, nil, tcell.StyleDefault.Foreground(tcell.ColorLightBlue))
					}
				}
			}
		}
	} else if gs.LoreMessage != "" {
		// Afficher le message de lore en vert clair
		for i, r := range gs.LoreMessage {
			if i < screenWidth {
				gs.Screen.SetContent(i, loreY, r, nil, tcell.StyleDefault.Foreground(tcell.ColorLightGreen))
			}
		}
	}

	if len(availableInteractions) > 0 {
		bottomText := availableInteractions[0] // Prendre la première interaction
		for i, r := range bottomText {
			if i < screenWidth {
				gs.Screen.SetContent(i, bottomY, r, nil, tcell.StyleDefault.Foreground(tcell.ColorYellow))
			}
		}
	} else {
		// Afficher les commandes de base
		defaultText := "Flèches: déplacer • [E]: interagir • [I]: inventaire • [TAB]: changer de monde • [Q]: quitter"
		for i, r := range defaultText {
			if i < screenWidth {
				gs.Screen.SetContent(i, bottomY, r, nil, tcell.StyleDefault.Foreground(tcell.ColorGray))
			}
		}
	}

	gs.Screen.Show()
}

// StartRespawnChecker démarre la vérification périodique des respawns
func (gs *GameState) StartRespawnChecker() *time.Ticker {
	respawnTicker := time.NewTicker(1 * time.Second)

	go func() {
		for range respawnTicker.C {
			w := gs.WorldList[gs.CurrentWorld]
			respawnMessages := gs.InteractionManager.CheckRespawns(w)
			if len(respawnMessages) > 0 {
				for _, msg := range respawnMessages {
					gs.LoreMessage = msg
				}
				gs.Draw() // Redessiner quand un respawn a lieu
			}
		}
	}()

	return respawnTicker
}
