package worlds

type World struct {
	Name          string
	Grid          [][]rune
	Width, Height int
	PlayerX       int
	PlayerY       int
	Config        *WorldConfig // Configuration du monde pour les interactions
	OriginalTile  rune         // Sauvegarde de la tuile originale sous le joueur
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
	}
}
