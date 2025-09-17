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
	Ended              bool // <- Indique si le jeu est terminé
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
		Ended:              false,
	}
}

// LoadWorlds charge les mondes depuis les fichiers de configuration
func (gs *GameState) LoadWorlds() {
	// Exemple de chargement, adapte selon ton projet
	ynoviaConfig, err := worlds.LoadWorldConfig("configs/ynovia.json")
	if err != nil {
		gs.WorldList = append(gs.WorldList, worlds.NewGrid("Monde Ynovia", 80, 35, 10, 10))
	} else {
		gs.WorldList = append(gs.WorldList, worlds.NewWorldFromConfig(ynoviaConfig))
	}

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
	inventoryCount := 0
	for _, qty := range gs.PlayerInventory.Items {
		inventoryCount += qty
	}
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

	// Dessiner le joueur
	gs.Screen.SetContent(w.PlayerX*2, w.PlayerY+1, '😀', nil, tcell.StyleDefault)
	gs.Screen.SetContent(w.PlayerX*2+1, w.PlayerY+1, ' ', nil, tcell.StyleDefault)

	// Interactions et lore
	loreY := w.Height + 2
	if gs.ShowingInventory {
		inventoryText := gs.PlayerInventory.GetInventoryString()
		lines := strings.Split(inventoryText, "\n")
		for i, line := range lines {
			if i < 10 && loreY+i < screenHeight-1 {
				for j, r := range line {
					if j < screenWidth {
						gs.Screen.SetContent(j, loreY+i, r, nil, tcell.StyleDefault.Foreground(tcell.ColorLightBlue))
					}
				}
			}
		}
	} else if gs.LoreMessage != "" {
		for i, r := range gs.LoreMessage {
			if i < screenWidth {
				gs.Screen.SetContent(i, loreY, r, nil, tcell.StyleDefault.Foreground(tcell.ColorLightGreen))
			}
		}
	}

	gs.Screen.Show()
}

// StartRespawnChecker démarre le ticker de respawn
func (gs *GameState) StartRespawnChecker() *time.Ticker {
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			w := gs.WorldList[gs.CurrentWorld]
			messages := gs.InteractionManager.CheckRespawns(w)
			if len(messages) > 0 {
				for _, msg := range messages {
					gs.LoreMessage = msg
				}
				gs.Draw()
			}
		}
	}()
	return ticker
}

// --- Fin du jeu ---
func (gs *GameState) EndGame() {
	gs.Ended = true
	gs.Screen.Clear()
	PrintEndGameAnimated(gs)
}
