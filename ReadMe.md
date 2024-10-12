# Maze Generator & Solver

Ce projet est une application en Go qui génère un labyrinthe procédural et propose plusieurs algorithmes pour le résoudre, notamment BFS (Breadth-First Search), DFS (Depth-First Search) et A*.

## Fonctionnalités

- Génération de labyrinthes aléatoires à l'aide de l'algorithme de Prim.
- Résolution du labyrinthe avec 3 algorithmes différents :
  - **BFS (Breadth-First Search)** - Affichage en vert.
  - **DFS (Depth-First Search)** - Affichage en bleu.
  - **A* (A-Star Search)** - Affichage en rouge.
- Réinitialisation et génération de nouveaux labyrinthes à chaque itération.
- Affichage visuel des chemins résolus.

## Prérequis

- [Go](https://golang.org/doc/install) (version 1.23.1+ recommandée)
- Make (pour la commande `make all`).
- Un compilateur gcc 64 bits.

## Installation

1. Clonez le dépôt sur votre machine locale :

    ```bash
    git clone https://github.com/votre-utilisateur/votre-projet.git
    cd Maze_Solver
    ```

1. Compilez et exécutez le projet en utilisant la commande `make` :

    ```bash
    make all
    ```

## Utilisation

- **Génération du labyrinthe :** Lorsque vous démarrez l'application, un labyrinthe aléatoire est généré.
- **Contrôles :**
  - Appuyez sur `1` pour résoudre le labyrinthe avec BFS (parcours en largeur).
  - Appuyez sur `2` pour résoudre le labyrinthe avec DFS (parcours en profondeur).
  - Appuyez sur `3` pour résoudre le labyrinthe avec A*.
  - Appuyez sur `ESPACE` pour générer un nouveau labyrinthe.

## Détails des algorithmes

### BFS (Breadth-First Search)
Le BFS explore les nœuds couche par couche, garantissant ainsi la découverte du chemin le plus court. Il est représenté en vert dans l'application.

### DFS (Depth-First Search)
Le DFS explore chaque chemin aussi profondément que possible avant de revenir en arrière. Il est représenté en bleu dans l'application.

### A* (A-Star Search)
L'algorithme A* utilise une heuristique (la distance de Manhattan) pour optimiser la recherche du chemin, combinant à la fois le coût et l'heuristique. Il est représenté en rouge dans l'application.
