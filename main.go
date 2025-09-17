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
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
<<<<<<< HEAD
			// Gestion des touches spéciales - Tab avec vérification du portail
			if ev.Key() == tcell.KeyTab {
=======
			// Gestion des touches spéciales
			switch ev.Key() {
			case tcell.KeyTab:
>>>>>>> origin/Mael2
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
				// Restaurer le terminal proprement
				screen.Fini()

				// Restaurer complètement le terminal avec reset
				cmd := exec.Command("reset")
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.Run()

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
<<<<<<< HEAD
				continue

			case ' ':
				gameState.HandleSpaceKey()
				continue

			case 'p', 'P':
				gameState.UnlockPortal()
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
=======
>>>>>>> origin/Mael2
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
