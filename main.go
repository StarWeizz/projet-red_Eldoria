package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
)

const (
	width  = 40
	height = 20
)

func main() {

	//inv := inventory.NewInventory()
	//inv.Add(items.WeaponList["Epée simple"], 1)
	//inv.Add(items.PotionsList["Heal potion"], 3)
	//inv.Add(items.CraftingItems["Bâton"], 2)
	//inv.List()

	// Comande inventaire en gros c'est des /give

	// Initialiser l’écran
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	defer screen.Fini()

	// Construire la grille
	grid := make([][]rune, height)
	for y := range grid {
		grid[y] = make([]rune, width)
		for x := range grid[y] {
			if y == 0 || y == height-1 || x == 0 || x == width-1 {
				grid[y][x] = '⬜'
			} else {
				grid[y][x] = '🟫'
			}
		}
	}

	// Ajouter un bloc spécial
	grid[10][10] = '🟨'

	// Position initiale du joueur
	px, py := 1, 1
	grid[py][px] = '😀'

	// Fonction d’affichage
	draw := func() {
		screen.Clear()
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				r := grid[y][x]
				screen.SetContent(x*2, y, r, nil, tcell.StyleDefault)
				screen.SetContent(x*2+1, y, ' ', nil, tcell.StyleDefault)
			}
		}
		screen.Show()
	}

	draw()

	// Fonction pour vérifier si joueur est à côté du bloc spécial
	checkInteraction := func() {
		coords := [][2]int{
			{px + 1, py}, {px - 1, py},
			{px, py + 1}, {px, py - 1},
		}
		for _, c := range coords {
			x, y := c[0], c[1]
			if grid[y][x] == '🟨' {
				screen.Fini() // désactiver l’écran pour afficher dans terminal
				fmt.Println("\n⚡ Vous êtes à côté d’un bloc jaune ! Tapez 'open' pour l’ouvrir.")
				reader := bufio.NewReader(os.Stdin)
				cmd, _ := reader.ReadString('\n')
				if cmd == "open\n" {
					fmt.Println("🎁 Coffre ouvert ! Vous avez trouvé une récompense.")
				} else {
					fmt.Println("❌ Commande incorrecte, rien ne se passe.")
				}

				// Réactiver l’écran
				if err := screen.Init(); err != nil {
					log.Fatalf("%+v", err)
				}
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
				if py < height-2 {
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
