package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"eldoria/worlds"

	"github.com/gdamore/tcell/v2"
)

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	showIntro(screen)

	defer screen.Fini()

	// Créer plusieurs mondes
	worldList := []*worlds.World{
		worlds.NewGrid("Monde Ynovia", 80, 35, 10, 10),
		worlds.NewGrid("Monde Eldoria", 40, 25, 5, 5),
	}
	currentWorld := 0

	world := worldList[currentWorld]
	world.Grid[world.PlayerY][world.PlayerX] = '😀'

	draw := func() {
		screen.Clear()
		w := worldList[currentWorld]

		// Topbar
		topbar := fmt.Sprintf("%s - 100/100 ♥ - %s", "joueur", w.Name)
		for i, r := range topbar {
			screen.SetContent(i, 0, r, nil, tcell.StyleDefault)
		}

		// Grille
		for y := 0; y < w.Height; y++ {
			for x := 0; x < w.Width; x++ {
				r := w.Grid[y][x]
				screen.SetContent(x*2, y+1, r, nil, tcell.StyleDefault)
				screen.SetContent(x*2+1, y+1, ' ', nil, tcell.StyleDefault)
			}
		}

		// Dessiner le joueur à sa position
		screen.SetContent(w.PlayerX*2, w.PlayerY+1, '😀', nil, tcell.StyleDefault)
		screen.SetContent(w.PlayerX*2+1, w.PlayerY+1, ' ', nil, tcell.StyleDefault)

		screen.Show()
	}

	checkInteraction := func() {
		w := worldList[currentWorld]
		coords := [][2]int{
			{w.PlayerX + 1, w.PlayerY},
			{w.PlayerX - 1, w.PlayerY},
			{w.PlayerX, w.PlayerY + 1},
			{w.PlayerX, w.PlayerY - 1},
		}
		for _, c := range coords {
			x, y := c[0], c[1]
			if x >= 0 && x < w.Width && y >= 0 && y < w.Height && w.Grid[y][x] == '🟨' {
				screen.Fini()
				fmt.Println("\n⚡ Vous êtes à côté d’un bloc jaune ! Tapez 'open' pour l’ouvrir.")
				reader := bufio.NewReader(os.Stdin)
				cmd, _ := reader.ReadString('\n')
				if cmd == "open\n" {
					fmt.Println("🎁 Coffre ouvert ! Vous avez trouvé une récompense.")
				} else {
					fmt.Println("❌ Commande incorrecte, rien ne se passe.")
				}
				if err := screen.Init(); err != nil {
					log.Fatalf("%+v", err)
				}
				draw()
			}
		}
	}

	draw()

	for {
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			w := worldList[currentWorld]

			if ev.Key() == tcell.KeyTab {
				// changer de monde
				// retirer le joueur du monde courant
				worldList[currentWorld].Grid[worldList[currentWorld].PlayerY][worldList[currentWorld].PlayerX] = '🟫'

				// changer de monde
				currentWorld = (currentWorld + 1) % len(worldList)
				world := worldList[currentWorld]

				// afficher le joueur à sa position sauvegardée
				world.Grid[world.PlayerY][world.PlayerX] = '😀'
				draw()
				continue
			}

			if ev.Rune() == 'q' {
				return
			}

			w.Grid[w.PlayerY][w.PlayerX] = '🟫'

			switch ev.Key() {
			case tcell.KeyUp:
				if w.PlayerY > 0 && isWalkable(w.Grid[w.PlayerY-1][w.PlayerX]) {
					w.PlayerY--
				}
			case tcell.KeyDown:
				if w.PlayerY < w.Height-1 && isWalkable(w.Grid[w.PlayerY+1][w.PlayerX]) {
					w.PlayerY++
				}
			case tcell.KeyRight:
				if w.PlayerX < w.Width-1 && isWalkable(w.Grid[w.PlayerY][w.PlayerX+1]) {
					w.PlayerX++
				}
			case tcell.KeyLeft:
				if w.PlayerX > 0 && isWalkable(w.Grid[w.PlayerY][w.PlayerX-1]) {
					w.PlayerX--
				}
			}

			draw()
			checkInteraction()

		case *tcell.EventResize:
			screen.Sync()
			draw()
		}
	}
}

func isWalkable(tile rune) bool {
	switch tile {
	case '🟫': // mur marron
		return true
	case '⬜': // mur blanc
		return false
	}
	return true
}

func showIntro(screen tcell.Screen) {
	screen.Clear()

	// Texte de bienvenue
	intro := []string{
		"👋 Bienvenue dans le jeu !",
		"",
		"Déplacez votre personnage avec les flèches du clavier.",
		"Changez de monde avec [TAB].",
		"Quittez avec [q].",
		"",
		"Appuyez sur [x] pour commencer...",
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
		}
	}
}
