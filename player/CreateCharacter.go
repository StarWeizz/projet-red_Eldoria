package createcharacter

import (
	"bufio"
	inventory "eldoria/Inventory"
	money "eldoria/money"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

// Structure Character
type Character struct {
	Name      string
	Class     string
	Level     int
	MaxHP     int
	CurrentHP int
	Gold      money.Money
	Inventory *inventory.Inventory
}

// Fonction utilitaire pour mettre la première lettre en majuscule
func capitalizeFirstLetter(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(strings.ToLower(s)) // mettre tout en minuscule d'abord
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// Retourne (true, "") si valide, sinon (false, raison)
func validateName(name string) (bool, string) {
	runes := []rune(name)

	// Vérifie longueur (en runes, pour gérer les accents correctement)
	if len(runes) == 0 {
		return false, "Le nom ne peut pas être vide."
	}
	if len(runes) > 15 {
		return false, fmt.Sprintf("Le nom est trop long (%d caractères). Maximum : 15.", len(runes))
	}

	// Vérifie que chaque rune est une lettre (lettres latines + accents autorisés)
	for i, r := range runes {
		if !unicode.IsLetter(r) {
			return false, fmt.Sprintf(
				"Caractère non autorisé '%c' à la position %d : seules les lettres sont autorisées.",
				r, i+1,
			)
		}
	}

	return true, ""
}

// Fonction pour créer un personnage personnalisé
func CreateCharacter() *Character {
	reader := bufio.NewReader(os.Stdin)

	var name string
	for {
		fmt.Print("Entrez le nom de votre personnage : ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if ok, reason := validateName(input); ok {
			name = capitalizeFirstLetter(input)
			break
		} else {
			fmt.Println("Nom invalide :", reason)
		}
	}

	// Choix de la classe
	classes := []string{"Guerrier", "Mage", "Chasseur"}
	fmt.Println("Choisissez la classe de votre personnage :")
	for i, class := range classes {
		fmt.Printf("%d. %s\n", i+1, class)
	}

	var classChoice int
	for {
		fmt.Print("Entrez le numéro de la classe : ")
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		n, err := strconv.Atoi(line)
		if err == nil && n >= 1 && n <= len(classes) {
			classChoice = n
			break
		}
		fmt.Println("Choix invalide, réessayez.")
	}

	chosenClass := classes[classChoice-1]

	// HP de base selon la classe
	maxHP := 100
	switch chosenClass {
	case "Guerrier":
		maxHP = 100
	case "Mage":
		maxHP = 80
	case "Chasseur":
		maxHP = 90
	}

	return &Character{
		Name:      name,
		Class:     chosenClass,
		Level:     1,
		MaxHP:     maxHP,
		CurrentHP: maxHP,
		Gold:      *money.NewMoney(100),     // chaque perso démarre avec 100 or
		Inventory: inventory.NewInventory(), // inventaire vide au départ
	}
}
