package main

import (
	"log"

	"eldoria/game"

	"github.com/gdamore/tcell/v2"
)

func main() {
	// Initialiser l'écran
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("Erreur écran: %+v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("Erreur Init écran: %+v", err)
	}
	defer screen.Fini()

	// Créer le personnage via l'intro
	playerCharacter := game.ShowIntroAndCreateCharacter(screen)

	// Créer l'état du jeu avec le joueur
	gameState := game.NewGameState(screen, playerCharacter)

	// Charger les mondes
	gameState.LoadWorlds()

	// Initialiser le joueur dans le monde courant
	gameState.InitializePlayer()

	// Démarrer le système de respawn
	respawnTicker := gameState.StartRespawnChecker()
	defer respawnTicker.Stop()

	// Dessiner l'état initial
	gameState.Draw()

	// Boucle principale du jeu
	for {
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			// Gestion des touches spéciales
			switch ev.Key() {
			case tcell.KeyTab:
				gameState.SwitchWorld()
				gameState.Draw()
				continue
			case tcell.KeyUp, tcell.KeyDown, tcell.KeyLeft, tcell.KeyRight:
				if gameState.MovePlayer(ev.Key()) {
					gameState.Draw()
					gameState.CheckInteraction()
				}
				continue
			}

			// Gestion des touches par caractère
			switch ev.Rune() {
			case 'q', 'Q':
				screen.Clear()
				screen.Show()
				return
			case 'e', 'E':
				gameState.HandleInteractionKey()
				gameState.Draw()
			case 'i', 'I':
				gameState.ToggleInventory()
				gameState.Draw()
			case '1', '2', '3', '4', '5':
				itemIndex := int(ev.Rune() - '1')
				gameState.HandleShopPurchase(itemIndex)
				gameState.Draw()
			}

		case *tcell.EventResize:
			screen.Sync()
			gameState.Draw()
		}
	}
}
