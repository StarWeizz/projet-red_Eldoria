package main

import (
	"bufio"
	"eldoria/forgeron"
	"eldoria/marchant"
	createcharacter "eldoria/player"
	"fmt"
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
)

const (
	width  = 40
	height = 22 // on garde 20 cases + 2 lignes pour HUD
)

func main() {
	// Créer un joueur
	hero := createcharacter.CreateCharacter()

	// Créer PNJ
	marchand := marchant.NewMerchant("Jean le Marchand")
	forgeron := forgeron.NewBlacksmith("Durin le Forgeron")

	// Initialiser écran
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	defer screen.Fini()

	// Construire la grille
	grid := make([][]rune, height-2) // -2 car on garde 2 lignes pour HUD
	for y := range grid {
		grid[y] = make([]rune, width)
		for x := range grid[y] {
			if y == 0 || y == len(grid)-1 || x == 0 || x == width-1 {
				grid[y][x] = '⬜'
			} else {
				grid[y][x] = '🟫'
			}
		}
	}

	// Ajouter PNJ
	grid[5][5] = '🔵'   // Marchand
	grid[10][10] = '🔴' // Forgeron

	// Position initiale du joueur
	px, py := 1, 1
	grid[py][px] = '😀'

	// Fonction d’affichage
	draw := func() {
		screen.Clear()
		for y := 0; y < len(grid); y++ {
			for x := 0; x < width; x++ {
				r := grid[y][x]
				screen.SetContent(x*2, y, r, nil, tcell.StyleDefault)
				screen.SetContent(x*2+1, y, ' ', nil, tcell.StyleDefault)
			}
		}

		// Affichage de l’or
		goldText := fmt.Sprintf("💰 Or: %d", hero.Gold.Get())
		for i, r := range goldText {
			screen.SetContent(i, height-2, r, nil, tcell.StyleDefault)
		}

		// Affichage de l’inventaire
		invText := "🎒 Inventaire: "
		for name, qty := range hero.Inventory.Items {
			invText += fmt.Sprintf("%s x%d  ", name, qty)
		}

		for i, r := range invText {
			if i < width*2 {
				screen.SetContent(i, height-1, r, nil, tcell.StyleDefault)
			}
		}

		screen.Show()
	}

	draw()

	// Vérifier si joueur est à côté d’un PNJ
	checkInteraction := func() {
		coords := [][2]int{
			{px + 1, py}, {px - 1, py},
			{px, py + 1}, {px, py - 1},
		}
		for _, c := range coords {
			x, y := c[0], c[1]
			switch grid[y][x] {
			case '🔵': // Marchand
				screen.Suspend()
				fmt.Println("\n💰 Vous parlez au marchand !")
				marchand.ShowStock()
				fmt.Println("Tapez le nom d’un objet pour l’acheter, ou 'exit' pour quitter.")
				reader := bufio.NewReader(os.Stdin)
				cmd, _ := reader.ReadString('\n')
				if cmd != "exit\n" {
					marchand.Buy(hero, cmd[:len(cmd)-1])
				}
				fmt.Println("Appuyez sur Entrée pour continuer...")
				reader.ReadString('\n')
				screen.Resume()
				draw()

			case '🔴': // Forgeron
				screen.Suspend()
				fmt.Println("\n⚒️ Vous parlez au forgeron !")
				forgeron.ShowStock()
				fmt.Println("Tapez le nom d’une arme pour l’acheter, ou 'exit' pour quitter.")
				reader := bufio.NewReader(os.Stdin)
				cmd, _ := reader.ReadString('\n')
				if cmd != "exit\n" {
					forgeron.Buy(hero, cmd[:len(cmd)-1])
				}
				fmt.Println("Appuyez sur Entrée pour continuer...")
				reader.ReadString('\n')
				screen.Resume()
				draw()
			}
		}
	}

	// Boucle d’événements
	for {
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Rune() == 'q' {
				return
			}

			grid[py][px] = '🟫'

			switch ev.Key() {
			case tcell.KeyUp:
				if py > 1 {
					py--
				}
			case tcell.KeyDown:
				if py < len(grid)-2 {
					py++
				}
			case tcell.KeyRight:
				if px < width-2 {
					px++
				}
			case tcell.KeyLeft:
				if px > 1 {
					px--
				}
			}

			grid[py][px] = '😀'
			draw()
			checkInteraction()

		case *tcell.EventResize:
			screen.Sync()
			draw()
		}
	}
}
