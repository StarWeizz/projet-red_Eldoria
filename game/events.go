package game

// HandleInteractionKey gère la touche E
func (gs *GameState) HandleInteractionKey() {
	w := gs.WorldList[gs.CurrentWorld]
	// Vérifie la case du joueur et les cases adjacentes
	coords := [][2]int{
		{w.PlayerX, w.PlayerY},
		{w.PlayerX + 1, w.PlayerY},
		{w.PlayerX - 1, w.PlayerY},
		{w.PlayerX, w.PlayerY + 1},
		{w.PlayerX, w.PlayerY - 1},
	}
	for _, coord := range coords {
		x, y := coord[0], coord[1]
		if x >= 0 && x < w.Width && y >= 0 && y < w.Height {
			interaction := w.GetInteractionType(x, y)
			if interaction != "none" && interaction != "" {
				result := gs.InteractionManager.HandleInteraction(w, gs.PlayerCharacter, x, y, interaction)
				if result != nil {
					gs.LoreMessage = result.Message
					if result.ShouldRemove {
						w.RemoveObject(x, y)
					}
					if result.EndGame {
						gs.EndGame()
					}
				}
				break // Une seule interaction par touche E
			}
		}
	}
}

// ToggleInventory bascule l'affichage de l'inventaire
func (gs *GameState) ToggleInventory() {
	gs.ShowingInventory = !gs.ShowingInventory
}

// HandleShopPurchase gère l'achat depuis le marchand
func (gs *GameState) HandleShopPurchase(itemIndex int) {
	result := gs.InteractionManager.BuyItem(itemIndex)
	if result != nil {
		gs.LoreMessage = result.Message
	}
}
