# le-password

Générateur de mots de passe en ligne de commande pour **Windows**, écrit en Go.

Les mots de passe sont générés à partir du générateur aléatoire cryptographique de Windows (`BCryptGenRandom`, via `bcrypt.dll`) et affichés dans la console.

```
 3  passwords generated
╭──────────────────────────────╮
│  › vW}H2tgsr7xuHLKCK]TI      │
│  › Vfp!Ddd)CM#U)i@fC<D4      │
│  › um<C%^%pn~#zeySeu2Ls      │
╰──────────────────────────────╯
```

## Fonctionnalités

- Aléatoire cryptographique fourni directement par l'API Windows CNG (`BCryptGenRandom`).
- Chaque mot de passe contient au moins une minuscule, une majuscule, un chiffre et un symbole.

## Prérequis

- Windows 10/11
- Go 1.27 (pour le compiler)

## Installation

```sh
git clone https://github.com/l3nnyz/le-password.git
cd le-password
go build -o le-password.exe .
```

## Utilisation

```sh
le-password.exe [-l longueur] [-n nombre]
```

| Option | Description | Défaut |
|--------|-------------|--------|
| `-l`   | Longueur de chaque mot de passe (minimum 4) | `16` |
| `-n`   | Nombre de mots de passe à générer (minimum 1) | `1` |

Exemples :

```sh
le-password.exe                # un mot de passe de 16 caractères
le-password.exe -l 32          # un mot de passe de 32 caractères
le-password.exe -l 20 -n 5     # cinq mots de passe de 20 caractères
```

## Structure du projet

| Fichier | Rôle |
|---------|------|
| `main.go` | Lecture des options (`flag`) et affichage (Bubble Tea, Lip Gloss) |
| `gen.go` | Génération des mots de passe |
| `imports.go` | Chargement de `bcrypt.dll` et déclaration de `BCryptGenRandom` |

## Dépendances

- [golang.org/x/sys/windows](https://pkg.go.dev/golang.org/x/sys/windows)
- [github.com/charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea)
- [github.com/charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss)
