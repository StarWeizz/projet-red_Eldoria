package game

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

// PrintEndGameAnimated affiche la fin du jeu avec animation
func PrintEndGameAnimated(gs *GameState) {
	screen := gs.Screen
	screenWidth, screenHeight := screen.Size()
	screen.Clear()

	victoryLines := []string{
		" _   _                                         _____                               _ ",
		"| | | |                                       |  __ \\                             | |",
		"| | | | ___  _   _ ___    __ ___   _____ ____ | |  \\/ __ _  __ _ _ __   ___ _ __  | |",
		"| | | |/ _ \\| | | / __|  / _` \\ \\ / / _ \\_  / | | __ / _` |/ _` | '_ \\ / _ \\ '__| | |",
		"\\ \\_/ / (_) | |_| \\__ \\ | (_| |\\ V /  __// /  | |_\\ \\ (_| | (_| | | | |  __/ |    |_|",
		" \\___/ \\___/ \\__,_|___/  \\__,_| \\_/ \\___/___|  \\____/\\__,_|\\__, |_| |_|\\___|_|    (_)",
		"                                                            __/ |                    ",
		"                                                           |___/                     ",
		"",
		"                                🎉 YOU WIN ! 🎉                                ",
	}

	subtitle := []string{
		"───────────────────────────────────────────────",
		"             🎉 FÉLICITATIONS ! 🎉             ",
		"      Vous avez vaincu Maximor le Boss !      ",
		"     Le royaume d'Eldoria est sauvé ! 🏰     ",
		"───────────────────────────────────────────────",
	}

	subtitle = append(subtitle, "Bravo !!! "+gs.PlayerCharacter.Name+" !")

	message := "Merci d'avoir joué à Eldoria.  développeur ! Mathis Antonin et Maël"

	startY := (screenHeight - len(victoryLines) - len(subtitle)) / 2

	// Affiche bannière
	for i, line := range victoryLines {
		startX := (screenWidth - len(line)) / 2
		for j, r := range line {
			if startX+j < screenWidth {
				screen.SetContent(startX+j, startY+i, r, nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
				screen.Show()
				time.Sleep(5 * time.Millisecond)
			}
		}
	}

	time.Sleep(400 * time.Millisecond)

	// Affiche sous-titre
	for i, line := range subtitle {
		startX := (screenWidth - len(line)) / 2
		for j, r := range line {
			if startX+j < screenWidth {
				screen.SetContent(startX+j, startY+len(victoryLines)+i, r, nil, tcell.StyleDefault.Foreground(tcell.ColorYellow))
				screen.Show()
				time.Sleep(10 * time.Millisecond)
			}
		}
	}

	time.Sleep(600 * time.Millisecond)

	// Affiche message final
	startX := (screenWidth - len(message)) / 2
	for i, r := range message {
		screen.SetContent(startX+i, startY+len(victoryLines)+len(subtitle)+2, r, nil, tcell.StyleDefault.Foreground(tcell.ColorGreen))
		screen.Show()
		time.Sleep(15 * time.Millisecond)
	}

	// Afficher "Appuyez sur [Q] pour quitter"
	quitMessage := "Appuyez sur [Q] pour quitter le jeu."
	quitX := (screenWidth - len(quitMessage)) / 2
	for i, r := range quitMessage {
		screen.SetContent(quitX+i, startY+len(victoryLines)+len(subtitle)+4, r, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
	}

	screen.Show()
}
