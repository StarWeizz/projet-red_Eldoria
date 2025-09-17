package main

import (
	"log"
	"os"
	"os/exec"

	"eldoria/game"

	"github.com/gdamore/tcell/v2"
)

func main() {
	// Initialiser l'écran
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	defer screen.Fini()

	// Créer le personnage via l'intro
	playerCharacter := game.ShowIntroAndCreateCharacter(screen)

	// Créer l'état du jeu
	gameState := game.NewGameState(screen, playerCharacter)

	// Charger les mondes
	gameState.LoadWorlds()

	// Initialiser le joueur
	gameState.InitializePlayer()

	// Démarrer le système de respawn
	respawnTicker := gameState.StartRespawnChecker()
	defer respawnTicker.Stop()

	// Dessiner l'état initial
	gameState.Draw()

	// Boucle principale du jeu
	for {
		if gameState.Ended {
			break
		}

		// Vérifie si le héros est mort
		if gameState.PlayerCharacter.CurrentHP <= 0 {
			screen.Clear()
			loseMsg := []string{
				"__     ______  _    _   _      ____   _____ ______ ",
				"\\ \\   / / __ \\| |  | | | |    / __ \\ / ____|  ____|",
				" \\ \\_/ / |  | | |  | | | |   | |  | | (___ | |__   ",
				"  \\   /| |  | | |  | | | |   | |  | |\\___ \\|  __|  ",
				"   | | | |__| | |__| | | |___| |__| |____) | |____ ",
				"   |_|  \\____/ \\____/  |______\\____/|_____/|______|",
				"",
				"                                         ",
				"             VOUS AVEZ PERDU !           ",
				"                                         ",
				"Appuyez sur [Q] pour quitter le jeu."}
			screenWidth, screenHeight := screen.Size()
			startY := (screenHeight - len(loseMsg)) / 2
			for i, line := range loseMsg {
				startX := (screenWidth - len(line)) / 2
				for j, r := range line {
					if startX+j < screenWidth {
						screen.SetContent(startX+j, startY+i, r, nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
					}
				}
			}
			screen.Show()
			// Attend que l'utilisateur appuie sur Q pour quitter
			for {
				ev := screen.PollEvent()
				if keyEv, ok := ev.(*tcell.EventKey); ok {
					if keyEv.Rune() == 'q' || keyEv.Rune() == 'Q' {
						handleQuit(screen)
					}
				}
			}
		}

		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			// Gestion des touches spéciales - Tab avec vérification du portail
			if ev.Key() == tcell.KeyTab {
				gameState.SwitchWorld()
				gameState.Draw()
				continue
			}

			// Gestion des touches par caractère
			switch ev.Rune() {
			case 'q', 'Q':
				handleQuit(screen)

			case 'e', 'E':
				gameState.HandleInteractionKey()
				if gameState.Ended {
					// Quitte la boucle après la victoire
					break
				}
				continue

			case 'i', 'I':
				gameState.ToggleInventory()
				gameState.Draw()
				continue

			case '1', '2', '3', '4', '5':
				itemIndex := int(ev.Rune() - '1')
				gameState.HandleShopPurchase(itemIndex)
				gameState.Draw()
				continue

			case ' ':
				gameState.HandleSpaceKey()
				continue

			case 'p', 'P':
				gameState.UnlockPortal()
				gameState.Draw()
				continue

			case 'a', 'A':
				// Utiliser une potion de soin si disponible
				qty, ok := gameState.PlayerInventory.Items["Heal potion"]
				if ok && qty > 0 {
					potion, exists := gameState.PlayerInventory.Refs["Heal potion"]
					if exists {
						healValue := 20 // Valeur par défaut
						if p, ok := potion.(interface{ GetHeal() int }); ok {
							healValue = p.GetHeal()
						}
						gameState.PlayerCharacter.CurrentHP += healValue
						if gameState.PlayerCharacter.CurrentHP > gameState.PlayerCharacter.MaxHP {
							gameState.PlayerCharacter.CurrentHP = gameState.PlayerCharacter.MaxHP
						}
						gameState.PlayerInventory.Remove(potion, 1)
						gameState.LoreMessage = "💊 Vous avez utilisé une potion de soin (+20 PV) !"
					} else {
						gameState.LoreMessage = "Potion de soin introuvable dans la liste des références."
					}
				} else {
					gameState.LoreMessage = "Vous n'avez pas de potion de soin dans votre inventaire."
				}
				gameState.Draw()
				continue
			}

			// Gestion du mouvement
			if ev.Key() == tcell.KeyUp || ev.Key() == tcell.KeyDown ||
				ev.Key() == tcell.KeyLeft || ev.Key() == tcell.KeyRight {
				if gameState.MovePlayer(ev.Key()) {
					gameState.Draw()
					gameState.CheckInteraction()
				}
			}

		case *tcell.EventResize:
			screen.Sync()
			gameState.Draw()
		}
	}
}

func handleQuit(screen tcell.Screen) {
	screen.Fini()

	// Restaurer complètement le terminal avec reset
	cmd := exec.Command("reset")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	os.Exit(0)
}
