package game

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
)

// CheckInteraction vérifie les interactions du joueur dans le monde
func (gs *GameState) CheckInteraction() {
	w := gs.WorldList[gs.CurrentWorld]

	// Vérifier les respawns (sans afficher de messages)
	gs.InteractionManager.CheckRespawns(w)

	// Vérifier si le joueur est sur une porte (interaction automatique)
	currentInteraction := w.GetInteractionType(w.PlayerX, w.PlayerY)
	if currentInteraction == "door" {
		result := gs.InteractionManager.HandleInteraction(w, gs.PlayerCharacter, w.PlayerX, w.PlayerY, "door")

		if result.Success {
			// Mettre à jour le message de lore au lieu de quitter l'écran
			gs.LoreMessage = result.Message
			// Le message s'affichera automatiquement lors du prochain draw()
		}
	} else {
		// Effacer le message de lore si le joueur n'est plus sur une porte
		gs.LoreMessage = ""
	}
}

// HandleInteractionKey gère l'interaction avec la touche E
func (gs *GameState) HandleInteractionKey() {
	// Vérifier d'abord si le joueur est près du portail pour la téléportation
	if gs.CheckPortalProximity() {
		gs.TeleportToEldoria()
		gs.Draw()
		return
	}

	w := gs.WorldList[gs.CurrentWorld]
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
				result := gs.InteractionManager.HandleInteraction(w, gs.PlayerCharacter, x, y, interactionType)

				// Afficher le message dans la zone de lore au lieu de quitter l'écran
				gs.LoreMessage = result.Message

				if result.Success && result.ShouldRemove {
					// Supprimer l'objet de la grille
					w.RemoveObject(x, y)
				}

				// Vérifier si le jeu doit se terminer (victoire contre le boss)
				if result.EndGame {
					gs.EndGame()
					return
				}

				// Vérifier si le portail doit être débloqué
				if result.UnlockPortal {
					gs.PortalUnlocked = true
				}

				// Redessiner immédiatement pour afficher le message
				gs.Draw()
				return // Sortir après la première interaction trouvée
			}
		}
	}
}

// SwitchWorld change de monde (TAB) - seulement si le portail est débloqué
func (gs *GameState) SwitchWorld() {
	if !gs.PortalUnlocked {
		gs.LoreMessage = "❌ Le portail est verrouillé ! Vous devez d'abord débloquer l'accès au portail."
		return
	}

	// restaurer la tuile originale dans le monde courant
	gs.WorldList[gs.CurrentWorld].Grid[gs.WorldList[gs.CurrentWorld].PlayerY][gs.WorldList[gs.CurrentWorld].PlayerX] = gs.WorldList[gs.CurrentWorld].OriginalTile

	// changer de monde
	gs.CurrentWorld = (gs.CurrentWorld + 1) % len(gs.WorldList)
	world := gs.WorldList[gs.CurrentWorld]

	// afficher le joueur à sa position sauvegardée
	world.Grid[world.PlayerY][world.PlayerX] = '😀'

	gs.LoreMessage = fmt.Sprintf("🌍 Téléportation vers %s réussie !", world.Name)
}

// ToggleInventory bascule l'affichage de l'inventaire
func (gs *GameState) ToggleInventory() {
	gs.ShowingInventory = !gs.ShowingInventory
}

// HandleShopPurchase gère les achats dans la boutique (touches 1-5)
func (gs *GameState) HandleShopPurchase(itemIndex int) {
	w := gs.WorldList[gs.CurrentWorld]
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
				result := gs.InteractionManager.BuyItem(itemIndex)
				gs.LoreMessage = result.Message
				return
			}
		}
	}
}

// HandleSpaceKey gère l'appui sur la barre d'espace (pour Emeryn)
func (gs *GameState) HandleSpaceKey() {
	w := gs.WorldList[gs.CurrentWorld]
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
			if interactionType == "emeryn" {
				// Vérifier si on peut faire avancer l'interaction
				if gs.InteractionManager.CanAdvanceEmerynInteraction() {
					// Forcer une synchronisation complète de l'écran
					gs.Screen.Sync()

					// Clear la lore et faire avancer l'interaction
					gs.LoreMessage = ""
					gs.InteractionManager.AdvanceEmerynInteraction()

					// Récupérer le nouveau message
					result := gs.InteractionManager.HandleInteraction(w, gs.PlayerCharacter, x, y, "emeryn")
					gs.LoreMessage = result.Message

					// Un seul redraw final avec sync
					gs.Draw()
					gs.Screen.Sync()
				}
				return
			}
		}
	}
}

// MovePlayer déplace le joueur dans une direction
func (gs *GameState) MovePlayer(direction tcell.Key) bool {
	w := gs.WorldList[gs.CurrentWorld]

	// Restaurer la tuile originale à l'ancienne position
	w.Grid[w.PlayerY][w.PlayerX] = w.OriginalTile

	moved := false
	switch direction {
	case tcell.KeyUp:
		if w.PlayerY > 0 && w.IsWalkableFromConfig(w.Grid[w.PlayerY-1][w.PlayerX]) {
			w.PlayerY--
			moved = true
		}
	case tcell.KeyDown:
		if w.PlayerY < w.Height-1 && w.IsWalkableFromConfig(w.Grid[w.PlayerY+1][w.PlayerX]) {
			w.PlayerY++
			moved = true
		}
	case tcell.KeyRight:
		if w.PlayerX < w.Width-1 && w.IsWalkableFromConfig(w.Grid[w.PlayerY][w.PlayerX+1]) {
			w.PlayerX++
			moved = true
		}
	case tcell.KeyLeft:
		if w.PlayerX > 0 && w.IsWalkableFromConfig(w.Grid[w.PlayerY][w.PlayerX-1]) {
			w.PlayerX--
			moved = true
		}
	}

	if moved {
		// Sauvegarder la nouvelle tuile originale
		w.OriginalTile = w.Grid[w.PlayerY][w.PlayerX]
	} else {
		// Remettre le joueur à sa position si le mouvement a échoué
		w.Grid[w.PlayerY][w.PlayerX] = '😀'
	}

	return moved
}