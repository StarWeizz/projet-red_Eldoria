package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func main() {
	oldState, _ := term.MakeRaw(int(os.Stdin.Fd()))
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	buf := make([]byte, 3)
	x, y := 5, 5 // position du joueur

	clearScreen()
	printAt(x, y, "X")

	for {
		n, _ := os.Stdin.Read(buf)
		if n == 1 && buf[0] == 'q' {
			break
		}

		if n == 3 && buf[0] == 27 && buf[1] == 91 {
			// efface ancienne position
			printAt(x, y, " ")

			switch buf[2] {
			case 'A': // haut
				y--
			case 'B': // bas
				y++
			case 'C': // droite
				x++
			case 'D': // gauche
				x--
			}
			// dessine nouvelle position
			printAt(x, y, "X")
		}
	}
}

// Efface l’écran
func clearScreen() {
	fmt.Print("\033[2J")
}

// Déplace le curseur et écrit
func printAt(x, y int, s string) {
	fmt.Printf("\033[%d;%dH%s", y, x, s)
}
