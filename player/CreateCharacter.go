package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

// Définition de la structure Character
type Character struct {
	Name      string
	Class     string
	Level     int
	MaxHP     int
	CurrentHP int
	Money     int
	Inventory []string
}

// Fonction utilitaire pour mettre la première lettre en majuscule
func capitalizeFirstLetter(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// Fonction pour créer un personnage personnalisé
func CreateCharacter() Character {
	reader := bufio.NewReader(os.Stdin)

	// Demander le nom
	fmt.Print("Entrez le nom de votre personnage : ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	name = capitalizeFirstLetter(name) // <-- Mise en majuscule de la 1ère lettre

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

	maxHP := 100
	switch chosenClass {
	case "Guerrier":
		maxHP = 100
	case "Mage":
		maxHP = 80
	case "Chasseur":
		maxHP = 90
	}

	return Character{
		Name:      name,
		Class:     chosenClass,
		Level:     1,
		MaxHP:     maxHP,
		CurrentHP: maxHP,
		Money:     100,
		Inventory: []string{},
	}
}

func main() {
	hero := CreateCharacter()
	fmt.Println("\nPersonnage créé avec succès !")
	fmt.Printf("%+v\n", hero)
}
