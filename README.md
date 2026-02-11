# goweeb-sama

> [!NOTE]
> English below but kind of useless since this utility is intended to dynamically download
> manga scans in **french** on a **french** website :)

**goweeb-sama** est un outil très rapide permettant de télécharger proprement les scans de ton manga préféré
depuis le site **anime-sama** (qui fait un excellent travail d'ailleurs, merci à vous).

## Installation

L'installation est très simple, il suffit d'aller dans l'onglet **Releases** sur la droite de l'interface GitHub et de télécharger la dernière version de l'outil pour ton système d'exploitation.

### Windows
1. Télécharge le fichier `goweeb-sama-windows.exe` (ou similaire)
2. Place-le dans un dossier de ton choix
3. Ouvre un terminal (PowerShell ou CMD) dans ce dossier et utilise l'outil

### Linux / macOS
1. Télécharge le fichier correspondant à ton système (`goweeb-sama-linux` ou `goweeb-sama-macos`)
2. Rends-le exécutable : `chmod +x goweeb-sama-*`
3. (Optionnel) Déplace-le dans `/usr/local/bin` pour l'utiliser depuis n'importe où

## Utilisation

### Format du nom de manga

**Important** : Le nom du manga doit être écrit en **minuscules** avec les **espaces remplacés par des tirets** (`-`).

Exemples :
- "One Piece" → `one-piece`
- "Jujutsu Kaisen" → `jujutsu-kaisen`
- "Chainsaw Man" → `chainsaw-man`

### Syntaxe de base

```bash
goweeb [options] <manga-slug>
```

### Options disponibles

| Option | Raccourci | Description |
|--------|-----------|-------------|
| `--all` | `-a` | Télécharge tous les chapitres disponibles |
| `--range <plage>` | `-r` | Spécifie une plage de chapitres (ex: `10-77`, `14-`) |
| `--scan-dir <dossier>` | `-d` | Dossier où sauvegarder les PDF (par défaut : `pdf`) |
| `--keep-images` | `-k` | Garde les images après la création du PDF |
| `--domain <url>` | `-u` | Remplace le domaine anime-sama (ex: `https://anime-sama.tv`) |

### Exemples d'utilisation

#### Télécharger tous les chapitres d'un manga
```bash
goweeb --all one-piece
# ou
goweeb -a one-piece
```

#### Télécharger une plage de chapitres
```bash
# Chapitres 10 à 77
goweeb --range 10-77 jujutsu-kaisen

# Raccourci
goweeb -r 10-50 one-piece

# Du chapitre 14 jusqu'au dernier disponible
goweeb -r 14- chainsaw-man
```

#### Spécifier un dossier de destination
```bash
# Les PDF seront sauvegardés dans le dossier "mes-mangas"
goweeb -d mes-mangas --all naruto

# Ou avec le chemin complet
goweeb --scan-dir ~/Documents/Mangas -r 1-100 one-piece
```

#### Spécifier un domaine personnalisé
```bash
# Si le domaine anime-sama change
goweeb -u https://anime-sama.fr --all one-piece
```

#### Combinaison d'options
```bash
# Télécharge les chapitres 1 à 50 de One piece, dans le dossier "scans" sur windows 
# en spécifiant un domaine de anime-sama personalisé
goweeb -r 1-50 -d 'C:\Users\<nom-de-ton-utilisateur>\scans' -u https://anime-sama.tv one-piece
```

### Conseils

> [!WARNING] 
> - Si jamais le manga n'est pas trouvé, peut-être qu'il faut utiliser son nom en japonais. Par exemple, `l'attaque des titans` porte le nom `shingeki no kyojin` sur anime-sama, il faut donc écrire `shingeki-no-kyojin` dans la commande
> - Garde bien en tête que le domaine d'anime-sama change pour des raisons évidentes.. donc n'hésites pas à vérifier [anime-sama.pw](https://anime-sama.pw) pour vérifier quel domaine est actif. Tu peux ensuite le spécifier avec l'argument `-u` comme présenté au dessus.

---

## 🇬🇧 English Version

**goweeb-sama** is a fast tool to download manga scans from the French website **anime-sama** (which does an excellent job by the way, thanks to them).

## Installation

Installation is very simple, just go to the **Releases** tab on the right side of the GitHub interface and download the latest version of the tool for your operating system.

### Windows
1. Download the file `goweeb-sama-windows.exe` (or similar)
2. Place it in a folder of your choice
3. Open a terminal (PowerShell or CMD) in that folder and use the tool

### Linux / macOS
1. Download the file corresponding to your system (`goweeb-sama-linux` or `goweeb-sama-macos`)
2. Make it executable: `chmod +x goweeb-sama-*`
3. (Optional) Move it to `/usr/local/bin` to use it from anywhere

## Usage

### Manga name format

**Important**: The manga name must be written in **lowercase** with **spaces replaced by hyphens** (`-`).

Examples:
- "One Piece" → `one-piece`
- "Jujutsu Kaisen" → `jujutsu-kaisen`
- "Chainsaw Man" → `chainsaw-man`

### Basic syntax

```bash
goweeb [options] <manga-slug>
```

### Available options

| Option | Shortcut | Description |
|--------|----------|-------------|
| `--all` | `-a` | Download all available chapters |
| `--range <range>` | `-r` | Specify a range of chapters (e.g., `10-77`, `14-`) |
| `--scan-dir <folder>` | `-d` | Folder to save PDF files (default: `pdf`) |
| `--keep-images` | `-k` | Keep images after PDF creation |
| `--domain <url>` | `-u` | Override anime-sama domain (e.g., `https://anime-sama.tv`) |

### Usage examples

#### Download all chapters of a manga
```bash
goweeb --all one-piece
# or
goweeb -a one-piece
```

#### Download a range of chapters
```bash
# Chapters 10 to 77
goweeb --range 10-77 jujutsu-kaisen

# Shortcut
goweeb -r 10-50 one-piece

# From chapter 14 to the last available
goweeb -r 14- chainsaw-man
```

#### Specify a destination folder
```bash
# PDFs will be saved in the "my-mangas" folder
goweeb -d my-mangas --all naruto

# Or with the full path
goweeb --scan-dir ~/Documents/Mangas -r 1-100 one-piece
```

#### Specify a custom domain
```bash
# If the anime-sama domain changes
goweeb -u https://anime-sama.fr --all one-piece
```

#### Combining options
```bash
# Download chapters 1 to 50 of One Piece, in the "scans" folder on Windows
# while specifying a custom anime-sama domain
goweeb -r 1-50 -d 'C:\Users\<your-username>\scans' -u https://anime-sama.tv one-piece
```

### Tips

> [!WARNING] 
> - If the manga is not found, you might need to use its Japanese name. For example, `Attack on Titan` is named `shingeki no kyojin` on anime-sama, so you need to write `shingeki-no-kyojin` in the command
> - Keep in mind that the anime-sama domain changes for obvious reasons... so don't hesitate to check [anime-sama.pw](https://anime-sama.pw) to verify which domain is active. You can then specify it with the `-u` argument as shown above.

---

## 🙏 Crédits

Merci beaucoup à Anime Sama pour leur travail colossal.
