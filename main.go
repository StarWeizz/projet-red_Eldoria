package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"

	inventory "eldoria/Inventory"
	"eldoria/interactions"
	"eldoria/money"
	"eldoria/worlds"
	createcharacter "eldoria/player"

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
	// Créer le personnage via l'intro
	playerCharacter := showIntroAndCreateCharacter(screen)

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

	// Utiliser les données du personnage créé
	playerMoney := &playerCharacter.Gold
	playerInventory := playerCharacter.Inventory
	interactionManager := interactions.NewInteractionManager(playerInventory, playerMoney)

	// Variable pour stocker le message de lore à afficher
	loreMessage := ""

	// Fonction pour compter les items dans l'inventaire
	getInventoryCount := func() int {
		count := 0
		for _, qty := range playerInventory.Items {
			count += qty
		}
		return count
	}

	world := worldList[currentWorld]
	// Sauvegarder la tuile originale et placer le joueur
	world.OriginalTile = world.Grid[world.PlayerY][world.PlayerX]
	world.Grid[world.PlayerY][world.PlayerX] = '😀'

	draw := func() {
		screen.Clear()
		w := worldList[currentWorld]
		screenWidth, screenHeight := screen.Size()

		// Topbar
		hiddenStatus := ""
		if w.IsPlayerHidden() {
			hiddenStatus = " - 🌿 CACHÉ des monstres"
		}
		inventoryCount := getInventoryCount()
		topbar := fmt.Sprintf("%s (%s) - %d/%d ♥ - 💰 %d - 🎒 %d items - %s - X:%d Y:%d%s",
			playerCharacter.Name, playerCharacter.Class,
			playerCharacter.CurrentHP, playerCharacter.MaxHP,
			playerMoney.Get(), inventoryCount, w.Name, w.PlayerX, w.PlayerY, hiddenStatus)
		for i, r := range topbar {
			if i < screenWidth {
				screen.SetContent(i, 0, r, nil, tcell.StyleDefault)
			}
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

		// Zone de lore - Afficher sous la grille si il y a un message
		loreY := w.Height + 2 // Juste sous la grille avec une ligne d'espace
		if loreMessage != "" {
			// Afficher le message de lore en vert clair
			for i, r := range loreMessage {
				if i < screenWidth {
					screen.SetContent(i, loreY, r, nil, tcell.StyleDefault.Foreground(tcell.ColorLightGreen))
				}
			}
		}

		// Bottombar - Afficher les interactions disponibles
		availableInteractions := interactionManager.CheckNearbyInteractions(w)
		bottomY := screenHeight - 1

		if len(availableInteractions) > 0 {
			bottomText := availableInteractions[0] // Prendre la première interaction
			for i, r := range bottomText {
				if i < screenWidth {
					screen.SetContent(i, bottomY, r, nil, tcell.StyleDefault.Foreground(tcell.ColorYellow))
				}
			}
		} else {
			// Afficher les commandes de base
			defaultText := "Déplacez-vous avec les flèches • [E] pour interagir • [TAB] changer de monde • [Q] quitter"
			for i, r := range defaultText {
				if i < screenWidth {
					screen.SetContent(i, bottomY, r, nil, tcell.StyleDefault.Foreground(tcell.ColorGray))
				}
			}
		}

		screen.Show()
	}

	checkInteraction := func() {
		w := worldList[currentWorld]

		// Vérifier les respawns
		respawnMessages := interactionManager.CheckRespawns(w)
		for _, msg := range respawnMessages {
			// Afficher brièvement les messages de respawn (optionnel)
			_ = msg
		}

		// Vérifier si le joueur est sur une porte (interaction automatique)
		currentInteraction := w.GetInteractionType(w.PlayerX, w.PlayerY)
		if currentInteraction == "door" {
			result := interactionManager.HandleInteraction(w, w.PlayerX, w.PlayerY, "door")

			if result.Success {
				// Mettre à jour le message de lore au lieu de quitter l'écran
				loreMessage = result.Message
				// Le message s'affichera automatiquement lors du prochain draw()
			}
		} else {
			// Effacer le message de lore si le joueur n'est plus sur une porte
			loreMessage = ""
		}
	}

	handleInteractionKey := func() {
		w := worldList[currentWorld]
		coords := [][2]int{
			{w.PlayerX + 1, w.PlayerY},
			{w.PlayerX - 1, w.PlayerY},
			{w.PlayerX, w.PlayerY + 1},
			{w.PlayerX, w.PlayerY - 1},
		}

		for _, coord := range coords {
			x, y := coord[0], coord[1]
			if x >= 0 && x < w.Width && y >= 0 && y < w.Height {
				interactionType := w.GetInteractionType(x, y)
				if interactionType != "none" && interactionType != "" && interactionType != "door" {
					result := interactionManager.HandleInteraction(w, x, y, interactionType)

					// Afficher le message dans la zone de lore au lieu de quitter l'écran
					loreMessage = result.Message

					if result.Success && result.ShouldRemove {
						// Supprimer l'objet de la grille
						w.RemoveObject(x, y)
					}

					// Redessiner immédiatement pour afficher le message
					draw()
					return // Sortir après la première interaction trouvée
				}
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

			if ev.Rune() == 'e' || ev.Rune() == 'E' {
				handleInteractionKey()
				continue
			}

			// Gestion des achats dans la boutique (touches 1-5)
			if ev.Rune() >= '1' && ev.Rune() <= '5' {
				// Vérifier si le joueur est près d'un marchand
				coords := [][2]int{
					{w.PlayerX + 1, w.PlayerY},
					{w.PlayerX - 1, w.PlayerY},
					{w.PlayerX, w.PlayerY + 1},
					{w.PlayerX, w.PlayerY - 1},
				}

				for _, coord := range coords {
					x, y := coord[0], coord[1]
					if x >= 0 && x < w.Width && y >= 0 && y < w.Height {
						interactionType := w.GetInteractionType(x, y)
						if interactionType == "merchant" {
							itemIndex := int(ev.Rune() - '1')
							result := interactionManager.BuyItem(itemIndex)
							loreMessage = result.Message
							draw() // Redessiner pour afficher le message
							continue
						}
					}
				}
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

func showIntroAndCreateCharacter(screen tcell.Screen) *createcharacter.Character {
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
		"Vous découvrirez sûrement que le village cache un portail qui mène vers un autre monde... Mais méfiez-vous des monstres qui rôdent... et du boss Maximor !",
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
				screen.Fini()
				character := createCharacterInTerminal()

				// Réinitialiser l'écran pour le jeu
				if err := screen.Init(); err != nil {
					log.Fatalf("Erreur lors de la réinitialisation de l'écran: %+v", err)
				}

				return character
			}
		}
	}
}

func createCharacterInTerminal() *createcharacter.Character {
	reader := bufio.NewReader(os.Stdin)

	// Demander le nom
	fmt.Print("Entrez le nom de votre personnage : ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	name = capitalizeFirstLetter(name)

	// Choix de la classe
	classes := []string{"Guerrier", "Mage", "Chasseur"}
	fmt.Println("Choisissez la classe de votre personnage :")
	for i, class := range classes {
		fmt.Printf("%d. %s\n", i+1, class)
	}

	var classChoice int
	for {
		fmt.Print("Entrez le numéro de la classe : ")
		fmt.Scan(&classChoice)
		if classChoice >= 1 && classChoice <= len(classes) {
			break
		}
		fmt.Println("Choix invalide, réessayez.")
	}

	chosenClass := classes[classChoice-1]

	// HP de base selon la classe
	maxHP := 100
	switch chosenClass {
	case "Guerrier":
		maxHP = 100
	case "Mage":
		maxHP = 80
	case "Chasseur":
		maxHP = 90
	}

	fmt.Printf("\n🎉 Personnage créé avec succès !\n")
	fmt.Printf("Nom: %s\n", name)
	fmt.Printf("Classe: %s\n", chosenClass)
	fmt.Printf("HP: %d/%d\n", maxHP, maxHP)
	fmt.Printf("\nAppuyez sur Entrée pour commencer l'aventure...")
	reader.ReadString('\n')

	return &createcharacter.Character{
		Name:      name,
		Class:     chosenClass,
		Level:     1,
		MaxHP:     maxHP,
		CurrentHP: maxHP,
		Gold:      *money.NewMoney(100),
		Inventory: inventory.NewInventory(),
	}
}

func capitalizeFirstLetter(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(strings.ToLower(s))
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
