package game

import (
	"bufio"
	"fmt"
	"log"
	"os"

	createcharacter "eldoria/player"

	"github.com/gdamore/tcell/v2"
)

// ShowIntroAndCreateCharacter affiche l'introduction et crée un personnage
func ShowIntroAndCreateCharacter(screen tcell.Screen) *createcharacter.Character {
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

				// Utiliser la fonction CreateCharacter existante
				character := createcharacter.CreateCharacter()

				fmt.Printf("\n🎉 Personnage créé avec succès !\n")
				fmt.Printf("Nom: %s\n", character.Name)
				fmt.Printf("Classe: %s\n", character.Class)
				fmt.Printf("HP: %d/%d\n", character.MaxHP, character.MaxHP)
				fmt.Printf("\nAppuyez sur Entrée pour commencer l'aventure...")

				bufio.NewReader(os.Stdin).ReadString('\n')

				// Réinitialiser l'écran pour le jeu
				if err := screen.Init(); err != nil {
					log.Fatalf("Erreur lors de la réinitialisation de l'écran: %+v", err)
				}

				return character
			}
		}
	}
}

