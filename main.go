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

	// Créer plusieurs mondes à partir des configurations JSON
	worldList := []*worlds.World{}

	// Charger le monde Ynovia depuis JSON
	ynoviaConfig, err := worlds.LoadWorldConfig("configs/ynovia.json")
	if err != nil {
		fmt.Printf("Erreur lors du chargement de ynovia.json: %v\n", err)
		fmt.Println("Utilisation du monde par défaut...")
		worldList = append(worldList, worlds.NewGrid("Monde Ynovia", 80, 35, 10, 10))
	} else {
		worldList = append(worldList, worlds.NewWorldFromConfig(ynoviaConfig))
	}

	// Charger le monde Eldoria depuis JSON
	eldoriaConfig, err := worlds.LoadWorldConfig("configs/eldoria.json")
	if err != nil {
		fmt.Printf("Erreur lors du chargement de eldoria.json: %v\n", err)
		fmt.Println("Utilisation du monde par défaut...")
		worldList = append(worldList, worlds.NewGrid("Monde Eldoria", 40, 25, 5, 5))
	} else {
		worldList = append(worldList, worlds.NewWorldFromConfig(eldoriaConfig))
	}
	currentWorld := 0

	world := worldList[currentWorld]
	// Sauvegarder la tuile originale et placer le joueur
	world.OriginalTile = world.Grid[world.PlayerY][world.PlayerX]
	world.Grid[world.PlayerY][world.PlayerX] = '😀'

	draw := func() {
		screen.Clear()
		w := worldList[currentWorld]

		// Topbar
		hiddenStatus := ""
		if w.IsPlayerHidden() {
			hiddenStatus = " - 🌿 CACHÉ des monstres"
		}
		topbar := fmt.Sprintf("%s - 100/100 ♥ - %s%s", "joueur", w.Name, hiddenStatus)
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
				// restaurer la tuile originale dans le monde courant
				worldList[currentWorld].Grid[worldList[currentWorld].PlayerY][worldList[currentWorld].PlayerX] = worldList[currentWorld].OriginalTile

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

			// Restaurer la tuile originale à l'ancienne position
			w.Grid[w.PlayerY][w.PlayerX] = w.OriginalTile

			switch ev.Key() {
			case tcell.KeyUp:
				if w.PlayerY > 0 && w.IsWalkableFromConfig(w.Grid[w.PlayerY-1][w.PlayerX]) {
					w.PlayerY--
				}
			case tcell.KeyDown:
				if w.PlayerY < w.Height-1 && w.IsWalkableFromConfig(w.Grid[w.PlayerY+1][w.PlayerX]) {
					w.PlayerY++
				}
			case tcell.KeyRight:
				if w.PlayerX < w.Width-1 && w.IsWalkableFromConfig(w.Grid[w.PlayerY][w.PlayerX+1]) {
					w.PlayerX++
				}
			case tcell.KeyLeft:
				if w.PlayerX > 0 && w.IsWalkableFromConfig(w.Grid[w.PlayerY][w.PlayerX-1]) {
					w.PlayerX--
				}
			}

			// Sauvegarder la nouvelle tuile originale
			w.OriginalTile = w.Grid[w.PlayerY][w.PlayerX]

			draw()
			checkInteraction()

		case *tcell.EventResize:
			screen.Sync()
			draw()
		}
	}
}

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
		}
	}
}
