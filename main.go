package main

import (
	"log"

	"eldoria/game"
	"eldoria/items"
	"fmt"
	"strings"

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
						return
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
				// Changer de monde
				gameState.CurrentWorld = (gameState.CurrentWorld + 1) % len(gameState.WorldList)
				gameState.Draw()
				continue
			case tcell.KeyUp, tcell.KeyDown, tcell.KeyLeft, tcell.KeyRight:
				// Déplacement du joueur
				w := gameState.WorldList[gameState.CurrentWorld]
				oldX, oldY := w.PlayerX, w.PlayerY
				// Calculer la nouvelle position
				newX, newY := w.PlayerX, w.PlayerY
				switch ev.Key() {
				case tcell.KeyUp:
					if w.PlayerY > 0 {
						newY--
					}
				case tcell.KeyDown:
					if w.PlayerY < w.Height-1 {
						newY++
					}
				case tcell.KeyLeft:
					if w.PlayerX > 0 {
						newX--
					}
				case tcell.KeyRight:
					if w.PlayerX < w.Width-1 {
						newX++
					}
				}
				// Vérifier si la case est praticable et sans ennemi
				targetTile := w.Grid[newY][newX]
				walkable := w.IsWalkableFromConfig(targetTile)
				enemyOnTile := false
				for i := range w.Config.Enemies {
					enemy := &w.Config.Enemies[i]
					if enemy.X == newX && enemy.Y == newY && enemy.HP > 0 {
						enemyOnTile = true
						break
					}
				}
				if walkable && !enemyOnTile {
					// Déplacement autorisé
					w.Grid[oldY][oldX] = w.OriginalTile
					w.PlayerX, w.PlayerY = newX, newY
					w.OriginalTile = w.Grid[w.PlayerY][w.PlayerX]
					w.Grid[w.PlayerY][w.PlayerX] = '😀'
				}
				// Affiche un message si interaction possible autour
				nearby := gameState.InteractionManager.CheckNearbyInteractions(w)
				if len(nearby) > 0 {
					gameState.LoreMessage = nearby[0]
				} else {
					gameState.LoreMessage = ""
				}
				gameState.Draw()
				continue
			}

			// Gestion des touches par caractère
			switch ev.Rune() {
			case 'c', 'C':
				// Affiche le menu de crafting sous la map
				craftable := []items.Recipe{}
				for _, recipe := range items.Recipes {
					canCraft := true
					needed := make(map[string]int)
					for _, item := range recipe.Needs {
						needed[item]++
					}
					for item, qty := range needed {
						if gameState.PlayerInventory.Items[item] < qty {
							canCraft = false
							break
						}
					}
					if canCraft {
						craftable = append(craftable, recipe)
					}
				}
				// Compose le menu
				menu := "🛠️ Menu de Crafting\n"
				if len(craftable) == 0 {
					menu += "Aucune recette craftable avec votre inventaire."
				} else {
					menu += "Sélectionnez le numéro de l'objet à crafter :\n"
					for i, recipe := range craftable {
						menu += fmt.Sprintf("%d. %s (nécessite: %v)\n", i+1, recipe.Result, recipe.Needs)
					}
				}
				menu += "\nAppuyez sur [1-%d] pour crafter, [i] pour fermer."
				menu = fmt.Sprintf(menu, len(craftable))
				// Affiche la map puis le menu en dessous
				gameState.Draw()
				screenWidth, _ := screen.Size()
				lines := strings.Split(menu, "\n")
				loreY := gameState.WorldList[gameState.CurrentWorld].Height + 3
				for i, line := range lines {
					for j, r := range line {
						if j < screenWidth {
							screen.SetContent(j, loreY+i, r, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
						}
					}
				}
				screen.Show()
				// Attend le choix de l'utilisateur sans bloquer le jeu
				// On ne bloque pas, on attend juste un chiffre ou 'i' à la prochaine touche
				// Le joueur peut continuer à jouer, le menu reste affiché jusqu'à la prochaine touche
				// On utilise une variable pour signaler que le menu est affiché
				gameState.LoreMessage = "[CRAFTING]"
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
			case 'q', 'Q':
				screen.Clear()
				screen.Show()
				return
			case 'e', 'E':
				gameState.HandleInteractionKey()
				if gameState.Ended {
					// Quitte la boucle après la victoire
					break
				}
				gameState.Draw()
			case 'i', 'I':
				gameState.ToggleInventory()
				gameState.Draw()
			case '1', '2', '3', '4', '5':
				if gameState.LoreMessage == "[CRAFTING]" {
					// Récupère la liste des craftables
					craftable := []items.Recipe{}
					for _, recipe := range items.Recipes {
						canCraft := true
						needed := make(map[string]int)
						for _, item := range recipe.Needs {
							needed[item]++
						}
						for item, qty := range needed {
							if gameState.PlayerInventory.Items[item] < qty {
								canCraft = false
								break
							}
						}
						if canCraft {
							craftable = append(craftable, recipe)
						}
					}
					idx := int(ev.Rune() - '1')
					if idx >= 0 && idx < len(craftable) {
						recipe := craftable[idx]
						needed := make(map[string]int)
						for _, item := range recipe.Needs {
							needed[item]++
						}
						for item, qty := range needed {
							if ref, ok := gameState.PlayerInventory.Refs[item]; ok {
								gameState.PlayerInventory.Remove(ref, qty)
							}
						}
						if craftItem, ok := items.CraftingItems[recipe.Result]; ok {
							gameState.PlayerInventory.Add(craftItem, 1)
							gameState.LoreMessage = "🛠️ Vous avez crafté : " + recipe.Result
						} else {
							gameState.LoreMessage = "Recette craftée : " + recipe.Result + " (objet non trouvé dans CraftingItems)"
						}
						gameState.Draw()
					}
				} else if gameState.LoreMessage != "" && (strings.Contains(gameState.LoreMessage, "Marchand") || strings.Contains(gameState.LoreMessage, "Bienvenue")) {
					itemIndex := int(ev.Rune() - '1')
					gameState.HandleShopPurchase(itemIndex)
					gameState.Draw()
				}
			}

		case *tcell.EventResize:
			screen.Sync()
			gameState.Draw()
		}
	}
}
