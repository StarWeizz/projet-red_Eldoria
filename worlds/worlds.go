package worlds

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

type World struct {
	Name          string
	Grid          [][]rune
	Width, Height int
	PlayerX       int
	PlayerY       int
	Config        *WorldConfig // Configuration du monde pour les interactions
	OriginalTile  rune         // Sauvegarde de la tuile originale sous le joueur
	Sticks        []Stick      // Liste des bâtons dans le monde
}

// Fonction pour créer une grille simple avec bordure
func NewGrid(name string, width, height int, specialX, specialY int) *World {
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
	if specialX >= 0 && specialY >= 0 && specialX < width && specialY < height {
		grid[specialY][specialX] = '🤭'
	}
	return &World{
		Name:         name,
		Grid:         grid,
		Width:        width,
		Height:       height,
		PlayerX:      1, // position initiale du joueur
		PlayerY:      1,
		Config:       nil, // Pas de configuration pour les mondes créés à l'ancienne
		OriginalTile: '🟫', // Tuile par défaut
		Sticks:       []Stick{},
	}
}

// GetObjectTypeAt retourne le type d'objet à une position donnée
func (w *World) GetObjectTypeAt(x, y int) string {
	if w.Config == nil {
		return ""
	}

	// Chercher dans les objets placés
	for _, obj := range w.Config.Objects {
		if obj.X == x && obj.Y == y {
			return obj.Object
		}
	}
	return ""
}

// GetObjectNameAt retourne le nom d'un objet à une position donnée
func (w *World) GetObjectNameAt(x, y int) string {
	if w.Config == nil {
		return ""
	}

	objectType := w.GetObjectTypeAt(x, y)
	if objectType != "" {
		if gameObj, exists := w.Config.GameObjects[objectType]; exists {
			return gameObj.Name
		}
	}
	return ""
}

// GetInteractionType retourne le type d'interaction pour un objet à une position
func (w *World) GetInteractionType(x, y int) string {
	if w.Config == nil {
		return "none"
	}

	objectType := w.GetObjectTypeAt(x, y)
	if objectType != "" {
		if gameObj, exists := w.Config.GameObjects[objectType]; exists {
			return gameObj.Interaction
		}
	}
	return "none"
}

// RespawnObject fait réapparaître un objet à une position donnée
func (w *World) RespawnObject(x, y int, objectType string) error {
	if w.Config == nil {
		return fmt.Errorf("pas de configuration disponible")
	}

	// Vérifier que l'objet existe dans la configuration
	if gameObj, exists := w.Config.GameObjects[objectType]; exists {
		// Vérifier les limites
		if x >= 0 && x < w.Width && y >= 0 && y < w.Height {
			w.Grid[y][x] = []rune(gameObj.Symbol)[0]

			// Ajouter l'objet à la liste des objets si pas déjà présent
			found := false
			for i, obj := range w.Config.Objects {
				if obj.X == x && obj.Y == y {
					w.Config.Objects[i].Object = objectType
					found = true
					break
				}
			}
			if !found {
				w.Config.Objects = append(w.Config.Objects, ObjectPlacement{
					X:      x,
					Y:      y,
					Object: objectType,
				})
			}
			return nil
		}
		return fmt.Errorf("position hors limites")
	}
	return fmt.Errorf("type d'objet inconnu: %s", objectType)
}

// RemoveObject supprime un objet à une position donnée
func (w *World) RemoveObject(x, y int) {
	if x >= 0 && x < w.Width && y >= 0 && y < w.Height {
		// Remplacer par la tuile par défaut
		if w.Config != nil && w.Config.DefaultTile != "" {
			w.Grid[y][x] = []rune(w.Config.DefaultTile)[0]
		} else {
			w.Grid[y][x] = '🟫'
		}

		// Supprimer de la liste des objets
		if w.Config != nil {
			for i, obj := range w.Config.Objects {
				if obj.X == x && obj.Y == y {
					w.Config.Objects = append(w.Config.Objects[:i], w.Config.Objects[i+1:]...)
					break
				}
			}
		}
	}
}

// Stick représente un bâton dans le monde
// Il contient des informations sur sa disponibilité et sa position
type Stick struct {
	X, Y        int  // Position du bâton
	IsAvailable bool // Indique si le bâton est disponible pour interaction
}

// InitializeSticks initialise les bâtons dans le monde
func (w *World) InitializeSticks() {
	w.Sticks = []Stick{
		{X: 5, Y: 10, IsAvailable: true},
		{X: 15, Y: 20, IsAvailable: true},
	}
}

// DrawSticks dessine les bâtons disponibles sur la carte
func (w *World) DrawSticks(screen tcell.Screen) {
	for _, stick := range w.Sticks {
		if stick.IsAvailable {
			screen.SetContent(stick.X, stick.Y, '⚪', nil, tcell.StyleDefault.Foreground(tcell.ColorGreen))
		}
	}
}
