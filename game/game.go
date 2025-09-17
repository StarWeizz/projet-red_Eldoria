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
	"github.com/mattn/go-runewidth"
)

// GameState représente l'état du jeu
type GameState struct {
<<<<<<< HEAD
	Screen              tcell.Screen
	WorldList           []*worlds.World
	CurrentWorld        int
	PlayerCharacter     *createcharacter.Character
	PlayerMoney         *money.Money
	PlayerInventory     *inventory.Inventory
	InteractionManager  *interactions.InteractionManager
	LoreMessage         string
	ShowingInventory    bool
	PortalUnlocked      bool
=======
	Screen             tcell.Screen
	WorldList          []*worlds.World
	CurrentWorld       int
	PlayerCharacter    *createcharacter.Character
	PlayerMoney        *money.Money
	PlayerInventory    *inventory.Inventory
	InteractionManager *interactions.InteractionManager
	LoreMessage        string
	ShowingInventory   bool
>>>>>>> origin/Mael2
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
		PortalUnlocked:     false,
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

// GetCurrentQuest récupère la quête actuelle du joueur
func (gs *GameState) GetCurrentQuest() string {
	if gs.InteractionManager == nil {
		return ""
	}

	// Accéder à la quête d'Emeryn
	quests := gs.InteractionManager.GetEmerynQuests()
	for _, quest := range quests {
		if !quest.Completed && quest.CurrentStep < len(quest.Steps) {
			currentStep := quest.Steps[quest.CurrentStep]

			// Personnaliser les titres selon l'étape
			switch quest.CurrentStep {
			case 0:
				return fmt.Sprintf("%s (%d/%d)", currentStep.Title, quest.CurrentStep+1, len(quest.Steps))
			case 1:
				return fmt.Sprintf("Tuer un Azador à la sortie du village (%d/%d)", quest.CurrentStep+1, len(quest.Steps))
			case 2:
				return fmt.Sprintf("Voir Valenric le forgeron (%d/%d)", quest.CurrentStep+1, len(quest.Steps))
			default:
				return fmt.Sprintf("%s (%d/%d)", currentStep.Title, quest.CurrentStep+1, len(quest.Steps))
			}
		}
	}

	return ""
}

// WrapText découpe un texte en lignes qui respectent la largeur maximale
// et respecte les sauts de ligne explicites (\n)
func (gs *GameState) WrapText(text string, maxWidth int) []string {
	if maxWidth <= 0 {
		return []string{text}
	}

	var allLines []string

	// D'abord séparer par les sauts de ligne explicites
	paragraphs := strings.Split(text, "\n")

	for _, paragraph := range paragraphs {
		if paragraph == "" {
			// Ligne vide
			allLines = append(allLines, "")
			continue
		}

		words := strings.Fields(paragraph)
		if len(words) == 0 {
			allLines = append(allLines, "")
			continue
		}

		currentLine := ""
		for _, word := range words {
			// Si ajouter ce mot dépasse la largeur max
			if len(currentLine)+len(word)+1 > maxWidth {
				if currentLine != "" {
					allLines = append(allLines, currentLine)
					currentLine = word
				} else {
					// Le mot lui-même est trop long, le couper
					for len(word) > maxWidth {
						allLines = append(allLines, word[:maxWidth])
						word = word[maxWidth:]
					}
					currentLine = word
				}
			} else {
				if currentLine != "" {
					currentLine += " " + word
				} else {
					currentLine = word
				}
			}
		}

		if currentLine != "" {
			allLines = append(allLines, currentLine)
		}
	}

	return allLines
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
	currentQuest := gs.GetCurrentQuest()
	questStatus := ""
	if currentQuest != "" {
		questStatus = fmt.Sprintf(" - ⚔️ %s", currentQuest)
	}

	// Formater avec des coordonnées de taille fixe pour éviter les variations
	// Afficher l'EXP selon le niveau
	expInfo := ""
	if gs.PlayerCharacter.Level >= 5 {
		expInfo = fmt.Sprintf("Lv%d(MAX)", gs.PlayerCharacter.Level)
	} else {
		nextLevelExp := gs.PlayerCharacter.GetExpForLevel(gs.PlayerCharacter.Level + 1)
		expInfo = fmt.Sprintf("Lv%d(%d/%d)", gs.PlayerCharacter.Level, gs.PlayerCharacter.Experience, nextLevelExp)
	}

	topbar := fmt.Sprintf("%s (%s) - %d/%d ♥ - %s - 💰 %d - 🎒 %d items - %s - X:%02d Y:%02d%s%s",
		gs.PlayerCharacter.Name, gs.PlayerCharacter.Class,
		gs.PlayerCharacter.CurrentHP, gs.PlayerCharacter.MaxHP,
		expInfo, gs.PlayerMoney.Get(), inventoryCount, w.Name, w.PlayerX, w.PlayerY, hiddenStatus, questStatus)

	// Effacer complètement la première ligne avant d'afficher la topbar
	for i := 0; i < screenWidth; i++ {
		gs.Screen.SetContent(i, 0, ' ', nil, tcell.StyleDefault)
	}

	// Afficher la topbar en gérant correctement la largeur des caractères Unicode
	displayPos := 0
	for _, r := range topbar {
		charWidth := runewidth.RuneWidth(r)

		// Vérifier s'il y a assez d'espace pour ce caractère
		if displayPos + charWidth > screenWidth {
			break
		}

		// Afficher le caractère
		gs.Screen.SetContent(displayPos, 0, r, nil, tcell.StyleDefault)

		// Si c'est un caractère large (emoji), marquer la position suivante comme occupée
		if charWidth == 2 {
			displayPos++
			if displayPos < screenWidth {
				gs.Screen.SetContent(displayPos, 0, ' ', nil, tcell.StyleDefault)
			}
		}

		displayPos++
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
	gs.Screen.SetContent(w.PlayerX*2, w.PlayerY+1, '😀', nil, tcell.StyleDefault)
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
		// Afficher le message de lore avec retour à la ligne automatique
		wrappedLines := gs.WrapText(gs.LoreMessage, screenWidth)
		for lineIndex, line := range wrappedLines {
			if lineIndex < 10 && loreY+lineIndex < bottomY { // Limiter à 10 lignes
				for charIndex, r := range line {
					if charIndex < screenWidth {
						gs.Screen.SetContent(charIndex, loreY+lineIndex, r, nil, tcell.StyleDefault.Foreground(tcell.ColorLightGreen))
					}
				}
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
<<<<<<< HEAD

// UnlockPortal débloque l'accès au portail
func (gs *GameState) UnlockPortal() {
	gs.PortalUnlocked = true
	gs.LoreMessage = "🌟 PORTAIL DÉBLOQUÉ ! Vous pouvez maintenant utiliser [TAB] pour changer de monde ou [E] près du portail pour vous téléporter !"
}


// CheckPortalProximity vérifie si le joueur est près du portail
func (gs *GameState) CheckPortalProximity() bool {
	if gs.CurrentWorld != 0 { // Le portail est seulement dans Ynovia (monde 0)
		return false
	}

	world := gs.WorldList[gs.CurrentWorld]
	portalX, portalY := 10, 10 // Position du portail dans ynovia.json

	// Vérifier si le joueur est adjacent au portail (distance de 1)
	distance := abs(world.PlayerX - portalX) + abs(world.PlayerY - portalY)
	return distance <= 1
}

// TeleportToEldoria téléporte le joueur vers Eldoria via le portail
func (gs *GameState) TeleportToEldoria() {
	if !gs.PortalUnlocked {
		gs.LoreMessage = "❌ Le portail est verrouillé ! Vous n'avez pas encore débloqué l'accès."
		return
	}

	if !gs.CheckPortalProximity() {
		gs.LoreMessage = "❌ Vous devez être près du portail pour l'utiliser !"
		return
	}

	if len(gs.WorldList) > 1 {
		// Retirer le joueur du monde actuel
		currentWorld := gs.WorldList[gs.CurrentWorld]
		currentWorld.Grid[currentWorld.PlayerY][currentWorld.PlayerX] = currentWorld.OriginalTile

		// Aller vers Eldoria (monde 1)
		gs.CurrentWorld = 1
		newWorld := gs.WorldList[gs.CurrentWorld]

		// Placer le joueur dans Eldoria
		newWorld.OriginalTile = newWorld.Grid[newWorld.PlayerY][newWorld.PlayerX]
		newWorld.Grid[newWorld.PlayerY][newWorld.PlayerX] = '😀'

		gs.LoreMessage = "🌟 Téléportation vers Eldoria via le portail réussie !"
	}
}

// abs retourne la valeur absolue d'un entier
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
=======
>>>>>>> origin/Mael2
