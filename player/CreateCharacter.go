package player
package createcharacter

import (
	"bufio"
	inventory "eldoria/Inventory"
	money "eldoria/money"
	"fmt"
	"os"
	"strings"
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

// Fonction pour créer un personnage personnalisé
func CreateCharacter() *Character {
	reader := bufio.NewReader(os.Stdin)

	// Demander le nom
	fmt.Print("Entrez le nom de votre personnage : ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

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