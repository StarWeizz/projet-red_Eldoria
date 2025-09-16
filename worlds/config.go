package worlds

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Objet représentant un élément sur la grille
type GameObject struct {
	Symbol      string `json:"symbol"`      // L'émoji ou caractère à afficher
	Name        string `json:"name"`        // Nom de l'objet
	Walkable    bool   `json:"walkable"`    // Le joueur peut-il marcher dessus ?
	Interaction string `json:"interaction"` // Type d'interaction (ex: "chest", "door", etc.)
}

// Configuration d'un monde
type WorldConfig struct {
	Name        string                `json:"name"`
	Width       int                   `json:"width"`
	Height      int                   `json:"height"`
	PlayerStartX int                  `json:"player_start_x"`
	PlayerStartY int                  `json:"player_start_y"`
	DefaultTile  string               `json:"default_tile"`
	BorderTile   string               `json:"border_tile"`
	Objects     []ObjectPlacement     `json:"objects"`
	GameObjects map[string]GameObject `json:"game_objects"`
}

// Placement d'un objet sur la grille
type ObjectPlacement struct {
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Object string `json:"object"` // Référence à un objet dans GameObjects
}

// Charger une configuration depuis un fichier JSON
func LoadWorldConfig(filename string) (*WorldConfig, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("impossible de lire le fichier %s: %v", filename, err)
	}

	var config WorldConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("impossible de parser le JSON: %v", err)
	}

	return &config, nil
}

// Créer un monde à partir d'une configuration
func NewWorldFromConfig(config *WorldConfig) *World {
	// Créer la grille de base
	grid := make([][]rune, config.Height)
	for y := range grid {
		grid[y] = make([]rune, config.Width)
		for x := range grid[y] {
			if y == 0 || y == config.Height-1 || x == 0 || x == config.Width-1 {
				// Bordures
				if config.BorderTile != "" {
					grid[y][x] = []rune(config.BorderTile)[0]
				} else {
					grid[y][x] = '⬜'
				}
			} else {
				// Tuile par défaut
				if config.DefaultTile != "" {
					grid[y][x] = []rune(config.DefaultTile)[0]
				} else {
					grid[y][x] = '🟫'
				}
			}
		}
	}

	// Placer les objets
	for _, obj := range config.Objects {
		if obj.X >= 0 && obj.X < config.Width && obj.Y >= 0 && obj.Y < config.Height {
			if gameObj, exists := config.GameObjects[obj.Object]; exists {
				grid[obj.Y][obj.X] = []rune(gameObj.Symbol)[0]
			}
		}
	}

	// Sauvegarder la tuile originale à la position du joueur
	originalTile := grid[config.PlayerStartY][config.PlayerStartX]
	
	return &World{
		Name:         config.Name,
		Grid:         grid,
		Width:        config.Width,
		Height:       config.Height,
		PlayerX:      config.PlayerStartX,
		PlayerY:      config.PlayerStartY,
		Config:       config, // Stocker la config pour les interactions
		OriginalTile: originalTile, // Sauvegarder la tuile originale
	}
}

// Vérifier si une tuile est praticable selon la configuration
func (w *World) IsWalkableFromConfig(tile rune) bool {
	if w.Config == nil {
		// Fallback vers l'ancien système
		return isWalkableOld(tile)
	}

	tileStr := string(tile)
	for _, gameObj := range w.Config.GameObjects {
		if gameObj.Symbol == tileStr {
			return gameObj.Walkable
		}
	}

	// Si pas trouvé dans la config, utiliser l'ancien système pour les tuiles communes
	return isWalkableOld(tile)
}

// Vérifier si le joueur est caché (sur une tuile avec interaction "hidden")
func (w *World) IsPlayerHidden() bool {
	if w.Config == nil {
		return false
	}

	currentTile := w.Grid[w.PlayerY][w.PlayerX]
	tileStr := string(currentTile)
	
	for _, gameObj := range w.Config.GameObjects {
		if gameObj.Symbol == tileStr && gameObj.Interaction == "hidden" {
			return true
		}
	}
	
	return false
}

func isWalkableOld(tile rune) bool {
	switch tile {
	case '🟫':
		return true
	case '⬜':
		return false
	}
	return true
}