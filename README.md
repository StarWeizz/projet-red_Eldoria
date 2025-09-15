# 🏰 Projet Eldoria

```
 _______   ___       ________  ________  ________  ___  ________     
|\  ___ \ |\  \     |\   ___ \|\   __  \|\   __  \|\  \|\   __  \    
\ \   __/|\ \  \    \ \  \_|\ \ \  \|\  \ \  \|\  \ \  \ \  \|\  \   
 \ \  \_|/_\ \  \    \ \  \ \\ \ \  \\\  \ \   _  _\ \  \ \   __  \  
  \ \  \_|\ \ \  \____\ \  \_\\ \ \  \\\  \ \  \\  \\ \  \ \  \ \  \ 
   \ \_______\ \_______\ \_______\ \_______\ \__\\ _\\ \__\ \__\ \__\
    \|_______|\|_______|\|_______|\|_______|\|__|\|__|\|__|\|__|\|__|
```

Un jeu d'aventure en mode console développé en Go avec tcell. Plongez dans le village mystérieux d'Ynovia, découvrez ses secrets, et partez à l'aventure dans le monde parallèle d'Eldoria !

## 🎮 Description

Eldoria est un jeu d'aventure textuel où vous incarnez un explorateur qui découvre le village d'Ynovia. Rencontrez Emeryn, le guide du village, et percez les mystères qui entourent ce lieu magique. Découvrez un portail vers un autre monde, mais attention aux monstres qui rôdent... et au redoutable boss Maximor !

**⚠️ Version Alpha** - Ce jeu est actuellement en développement. Certaines fonctionnalités peuvent être incomplètes ou instables.

## ✨ Fonctionnalités

- 🗺️ **Multiples mondes** - Explorez Ynovia et Eldoria avec des environnements uniques
- 🌳 **Système de cachette** - Cachez-vous dans les arbres pour éviter les monstres
- 🎛️ **Configuration JSON** - Personnalisez facilement vos mondes avec des fichiers de configuration
- 🎨 **Interface en console** - Affichage coloré avec des émojis et ASCII art
- ⚙️ **Système modulaire** - Architecture extensible pour ajouter facilement du contenu
- 🎯 **Interactions** - Coffres, objets à collecter, et mécaniques d'interaction

### Objets et environnements disponibles

- 🌳 **Arbres** - Cachettes pour éviter les ennemis
- 🏠 **Maisons** - Bâtiments à explorer
- 🏰 **Châteaux** - Structures imposantes
- 🟨 **Coffres** - Objets à ouvrir pour des récompenses
- 💎 **Trésors** - Objets précieux à collecter
- ⚔️ **Armes et équipements** - Épées, boucliers, potions
- 🐉 **Créatures** - Dragons et autres créatures mystiques

## 🎯 Commandes

- **Flèches directionnelles** - Déplacement du personnage
- **TAB** - Changer de monde
- **Q** - Quitter le jeu
- **X** - Commencer le jeu (écran d'accueil)

## 🚀 Installation et Lancement

### Prérequis

- Go 1.19 ou plus récent
- Terminal compatible (recommandé : terminal moderne avec support UTF-8)

### Installation

1. Clonez le projet
```bash
git clone https://github.com/StarWeizz/projet-red_Eldoria.git
cd projet-red_Eldoria
```

2. Installez les dépendances
```bash
go mod tidy
```

3. Lancez le jeu
```bash
go run main.go
```

## 🛠️ Configuration des Mondes

Le jeu utilise un système de configuration JSON pour personnaliser les mondes. Les fichiers de configuration se trouvent dans le dossier `configs/`.

### Structure d'un monde

```json
{
  "name": "Nom du monde",
  "width": 80,
  "height": 35,
  "player_start_x": 5,
  "player_start_y": 17,
  "default_tile": "🟫",
  "border_tile": "⬜",
  "game_objects": {
    "tree": {
      "symbol": "🌳",
      "name": "Arbre",
      "walkable": true,
      "interaction": "hidden"
    }
  },
  "objects": [
    {"x": 10, "y": 15, "object": "tree"}
  ]
}
```

### Création d'objets personnalisés

1. Ajoutez votre objet dans la section `game_objects`
2. Définissez ses propriétés (`symbol`, `walkable`, `interaction`)
3. Placez-le sur la carte dans la section `objects`

Consultez [CONFIG_README.md](CONFIG_README.md) pour plus de détails.

## 📁 Structure du Projet

```
projet-red_Eldoria/
├── main.go              # Point d'entrée du jeu
├── worlds/              # Package pour la gestion des mondes
│   ├── worlds.go        # Structures et fonctions de base
│   └── config.go        # Système de configuration JSON
├── configs/             # Fichiers de configuration des mondes
│   ├── ynovia.json      # Configuration du monde Ynovia
│   └── eldoria.json     # Configuration du monde Eldoria
├── CONFIG_README.md     # Guide de configuration détaillé
└── README.md           # Ce fichier
```

## 👥 Auteurs

- [@StarWeizz](https://github.com/StarWeizz) - Développeur principal
- [@mael](https://github.com/StarWeizz) - Contributeur
- [@mathis](https://github.com/StarWeizz) - Contributeur

## 📜 Licence

[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](https://choosealicense.com/licenses/mit/)

Ce projet est sous licence MIT. Voir le fichier `LICENSE` pour plus de détails.

## 🔄 Statut du Développement

![Version](https://img.shields.io/badge/Version-Alpha-orange)
![Go Version](https://img.shields.io/badge/Go-1.19+-blue)
![Platform](https://img.shields.io/badge/Platform-Cross--Platform-green)

### Fonctionnalités à venir

- 🗣️ Système de dialogues avec PNJ
- ⚔️ Combat et système de statistiques
- 📦 Inventaire et gestion d'objets
- 🎵 Effets sonores
- 💾 Sauvegarde de progression
- 🏆 Système de quêtes et objectifs

## 🤝 Contribution

Les contributions sont les bienvenues ! N'hésitez pas à :

1. Fork le projet
2. Créer une branche pour votre fonctionnalité
3. Commit vos changements
4. Ouvrir une Pull Request

## 🐛 Signaler un Bug

Si vous trouvez un bug, merci de créer une issue avec :
- Description du problème
- Étapes pour reproduire
- Environnement (OS, version de Go)
- Screenshots si applicable

---

*Prêt pour l'aventure ? Lancez le jeu et découvrez les secrets d'Eldoria !*