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
<<<<<<< HEAD
			gameState.Draw()
=======
			draw()
		}
	}
}
<<<<<<< HEAD
<<<<<<< HEAD

func showIntro(screen tcell.Screen) {
	screen.Clear()

	// Texte de bienvenue
	intro := []string{
		` _______   ___       ________  ________  ________  ___  ________     `,
		`|\  ___ \ |\  \     |\   ___ \|\   __  \|\   __  \|\  \|\   __  \    `,
		`\ \   __/|\ \  \    \ \  \_|\ \ \  \|\  \ \  \|\  \ \  \ \  \|\  \   `,
		` \ \  \_|/_\ \  \    \ \  \ \\ \ \  \\\  \ \   _  _\ \  \ \   __  \  `,
		`  \ \  \_|\ \ \  \____\ \  \_\\ \ \  \\\  \ \  \\  \\ \  \ \  \ \  \ `,
		`   \ \_______\ \_______\ \_______\ \_______\ \__\\ _\\ \__\ \__\ \__\`,
		`    \|_______|\|_______|\|_______|\|_______|\|__|\|__|\|__|\|__|\|__|`,
		`                                                                     `,
		`                                                                     `,
		`                                                                     `,
		"👋 Bienvenue dans Eldoria !",
		"",
		"Plongez dans le village d'Ynovia afin de percer ses mystères.",
		"Partez à la rencontre d'Emeryn, le guide du village, et écoutez le afin d'en apprendre davantage sur cet endroit.",
		"Vous découverirez sûrement que le village cache un portail qui mène vers un autre monde... Mais méfiez-vous des monstres qui rôdent... et du boss Maximor !",
		"",
		"⚠️ Attention : Ce jeu est en version alpha. Certaines fonctionnalités peuvent être incomplètes ou instables.",
		"",
		"Commandes :",
		"Déplacez votre personnage avec les flèches du clavier.",
		"Changez de monde avec [TAB].",
		"Quittez avec [q].",
		"",
		"Appuyez sur [x] pour commencer...",
		"",
	}

	// Centrer le texte
	w, h := screen.Size()
	for i, line := range intro {
		x := (w - len(line)) / 2
		y := h/2 - len(intro)/2 + i
		for j, r := range line {
			screen.SetContent(x+j, y, r, nil, tcell.StyleDefault)
		}
	}

	screen.Show()

	// Attendre l'appui sur "x"
	for {
		ev := screen.PollEvent()
		if key, ok := ev.(*tcell.EventKey); ok {
			if key.Rune() == 'x' || key.Rune() == 'X' {
				return
			}
>>>>>>> f8fb55b (Refactoring files)
		}
	}
}
=======
>>>>>>> origin/Mael2
=======

func handleQuit(screen tcell.Screen) {
	screen.Fini()

	// Restaurer complètement le terminal avec reset
	cmd := exec.Command("reset")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	os.Exit(0)
}
>>>>>>> antonin
