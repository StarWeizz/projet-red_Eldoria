# Configuration des Mondes - Eldoria

Ce système permet de personnaliser facilement vos mondes de jeu via des fichiers JSON.

## Structure des fichiers de configuration

Les fichiers de configuration se trouvent dans le dossier `configs/` et suivent cette structure :

```json
{
  "name": "Nom du monde",
  "width": 80,
  "height": 35,
  "player_start_x": 1,
  "player_start_y": 1,
  "default_tile": "🟫",
  "border_tile": "⬜",
  "game_objects": {
    "nom_objet": {
      "symbol": "🌳",
      "name": "Nom affiché",
      "walkable": false,
      "interaction": "type_interaction"
    }
  },
  "objects": [
    {"x": 5, "y": 5, "object": "nom_objet"}
  ]
}
```

## Paramètres principaux

- `name` : Nom affiché du monde
- `width`, `height` : Dimensions de la grille
- `player_start_x`, `player_start_y` : Position de départ du joueur
- `default_tile` : Symbole utilisé pour les cases vides
- `border_tile` : Symbole utilisé pour les bordures

## Objets de jeu

### Propriétés des objets
- `symbol` : L'émoji ou caractère affiché
- `name` : Nom descriptif de l'objet
- `walkable` : `true` si le joueur peut marcher dessus, `false` sinon
- `interaction` : Type d'interaction (pour futures fonctionnalités)

### Types d'interaction disponibles
- `"none"` : Aucune interaction
- `"chest"` : Coffre (compatible avec le système existant)
- `"door"` : Porte
- `"treasure"` : Trésor
- `"weapon"` : Arme
- `"armor"` : Armure
- `"heal"` : Objet de soin
- `"damage"` : Objet dangereux
- `"boss"` : Boss à combattre

## Placement des objets

Utilisez le tableau `objects` pour placer vos objets sur la grille :

```json
"objects": [
  {"x": 10, "y": 15, "object": "tree"},
  {"x": 20, "y": 25, "object": "house"}
]
```

## Exemples d'objets prédéfinis

### Objets de décoration
- `🌳` Arbre (non traversable)
- `🌸` Fleur (traversable)
- `🗿` Rocher (non traversable)
- `🌊` Eau (non traversable)
- `🌉` Pont (traversable)

### Bâtiments
- `🏠` Maison
- `🏰` Château

### Objets interactifs
- `🟨` Coffre
- `💎` Trésor
- `⚔️` Épée
- `🛡️` Bouclier
- `🧪` Potion

### Créatures
- `🐉` Dragon
- `🍄` Champignon

## Comment ajouter un nouveau monde

1. Créez un nouveau fichier JSON dans le dossier `configs/`
2. Suivez la structure décrite ci-dessus
3. Modifiez `main.go` pour charger votre nouveau monde :

```go
// Ajouter après les autres mondes
newWorldConfig, err := worlds.LoadWorldConfig("configs/votre_monde.json")
if err != nil {
    fmt.Printf("Erreur: %v\n", err)
} else {
    worldList = append(worldList, worlds.NewWorldFromConfig(newWorldConfig))
}
```

## Conseils de création

- Gardez les dimensions raisonnables (recommandé : max 100x50)
- Laissez toujours un chemin praticable pour le joueur
- Utilisez `walkable: false` pour créer des obstacles
- Les coordonnées commencent à (0,0) en haut à gauche
- Évitez de placer des objets sur les bordures (réservées aux murs)

## Débogage

Si un monde ne se charge pas :
1. Vérifiez la syntaxe JSON avec un validateur
2. Assurez-vous que tous les objets placés existent dans `game_objects`
3. Vérifiez que les coordonnées sont dans les limites du monde
4. Le jeu utilisera un monde par défaut en cas d'erreur