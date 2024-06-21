package main

import (
	// "math/rand"
	"fmt"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Mise à jour de l'état du jeu en fonction des entrées au clavier et des messages du serveur.
func (g *game) Update() error {

	g.stateFrame++

	switch g.gameState {
	case titleState:
		if g.nombreDeJoueurConnecte == 1  {
			select {
				case recu := <-g.channel:
					if recu == "connexion" {
						//si on reçoit connexion du serveur on considère que les deux joueurs sont connectés
						g.nombreDeJoueurConnecte = 2

					}else {
						g.channel <- recu
					}
				default:
			}
		}
		if g.titleUpdate() && (g.nombreDeJoueurConnecte == 2) {
			clear(g.channel)
			g.gameState++
		}

	case colorSelectState:

		select { // Position renvoyée par le serveur
			case recu := <-g.channel:
				//Si on reçois une nouvelle position de la couleur du joueur on l'actualise
				div := diviseur(recu)
				if div[0] == "couleur" {
					conv, err := strconv.Atoi(div[2])
					if err == nil {
						g.p2Color = conv
					}
					//Si c'est un choix on considère que l'autre joueur est pret
					if div[1] == "choix" {
						g.couleurPret[1] = true
					}
					
				}
			default:
		}
		
		// Si la touche entrée est appuyée et les deux joueurs sont prêts, on peut commence à jouer.
		if g.colorSelectUpdate() && (g.couleurPret[0] && g.couleurPret[1]) {
			g.gameState++
		}
	case playState:
		g.tokenPosUpdate()
		var lastXPositionPlayed int
		var lastYPositionPlayed int
		if g.turn == p1Turn {
			lastXPositionPlayed, lastYPositionPlayed = g.p1Update()
		} else {
			lastXPositionPlayed, lastYPositionPlayed = g.p2Update()
		}
		if lastXPositionPlayed >= 0 {
			finished, result := g.checkGameEnd(lastXPositionPlayed, lastYPositionPlayed)
			if finished {
				g.result = result
				g.gameState++
			}
		}
	case resultState:
		select {
			case recu := <-g.channel:
				if recu == "pret" {
					g.joueurPret[1] = true
					g.nombreJoueursPrets++
				}else {
					g.channel <- recu
				}
			default:
		}

		if g.resultUpdate() {
			g.reset()
			g.gameState = playState
			// On réinitialise les joeurs pret pour une prochaine fin de partie
			g.joueurPret = [2]bool{false, false}
			g.nombreJoueursPrets = 0
		}
	}

	return nil
}

// Mise à jour de l'état du jeu à l'écran titre.
func (g *game) titleUpdate() bool {
	g.stateFrame = g.stateFrame % globalBlinkDuration
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter)
}

// Mise à jour de l'état du jeu lors de la sélection des couleurs.
func (g *game) colorSelectUpdate() bool {
	// Calcul des coordonnées de la couleur dans la grille de sélection
	col := g.p1Color % globalNumColorCol
	line := g.p1Color / globalNumColorLine

	// Déplacement à droite
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		col = (col + 1) % globalNumColorCol
	}

	// Déplacement à gauche
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		col = (col - 1 + globalNumColorCol) % globalNumColorCol
	}	

	// Déplacement vers le bas
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		line = (line + 1) % globalNumColorLine
	}

	// Déplacement vers le haut
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		line = (line - 1 + globalNumColorLine) % globalNumColorLine
	}

	// Calcul de la nouvelle couleur sélectionnée, si c'est la même on ne fait rien
	newColor := line*globalNumColorLine + col
	if newColor != g.p1Color {
		g.p1Color = newColor
		writeClient(g.writer, "couleur" + "," + "selection" + "," + strconv.Itoa(g.p1Color))
	}

	// Validation de la couleur sélectionnée ou on regarde si les deux joueurs sont prêts pour lance sans appuyer sur la touche entrée.
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || (g.couleurPret[0] && g.couleurPret[1]) {
		if g.p2Color == g.p1Color {
			fmt.Println("Cette couleur est déjà prise par l'autre joueur.")
			return false
		}
		writeClient(g.writer, "couleur" + "," + "choix" + "," + strconv.Itoa(g.p1Color))
		g.couleurPret[0] = true
		return true
	}

	return false
}

// Gestion de la position du prochain pion à jouer par le joueur 1.
func (g *game) tokenPosUpdate() {
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.tokenPosition = (g.tokenPosition - 1 + globalNumTilesX) % globalNumTilesX
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.tokenPosition = (g.tokenPosition + 1) % globalNumTilesX
	}
}

// Gestion du moment où le prochain pion est joué par le joueur 1.
//Si on est joueur 1 on joue et on renvoie notre position sinon on attend un message de l'autre joueur.
func (g *game) p1Update() (int, int) {
	if (g.numeroJoueur == 1) {
		lastXPositionPlayed := -1
		lastYPositionPlayed := -1
		if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			if updated, yPos := g.updateGrid(p1Token, g.tokenPosition); updated {
				g.turn = p2Turn
				lastXPositionPlayed = g.tokenPosition
				lastYPositionPlayed = yPos
				conv := strconv.Itoa(lastXPositionPlayed)
				//envoi de la position jouer à l'autre joueur
				writeClient(g.writer, "p" + "," + conv )
			}
		}
		return lastXPositionPlayed, lastYPositionPlayed
	}else if (g.numeroJoueur == 2) {
		position_colonne := -1
		select { // Position renvoyée par le serveur
			case recu := <-g.channel:
				div := diviseur(recu)
				if div[0] == "p" {
					conv, err := strconv.Atoi(div[1])
					if err == nil {
						position_colonne = conv
					}
				}
			default:
		}

		if position_colonne == -1 {
			return -1, -1
		}
		position_ligne := -1
		if updated, nouvelle_position_ligne := g.updateGrid(p2Token, position_colonne); updated {
			g.turn = p2Turn
			position_ligne = nouvelle_position_ligne
		}
		return position_colonne, position_ligne
	}

	return -1, -1
}

// Gestion de la position du prochain pion joué par le joueur 2 et du moment où ce pion est joué.
//Si on est joueur 2 on joue et on renvoie notre position sinon on attend un message de l'autre joueur.
func (g *game) p2Update() (int, int) {
	if (g.numeroJoueur == 1) {
		position_colonne := -1 

		select { // Position renvoyée par le serveur
			case recu := <-g.channel:
				div := diviseur(recu)
				if div[0] == "p" {
					conv, err := strconv.Atoi(div[1])
					if err == nil {
						position_colonne = conv
					}
				}
			default:
		}

		if position_colonne == -1 {
			return -1, -1
		}
		position_ligne := -1
		if updated, nouvelle_position_ligne := g.updateGrid(p2Token, position_colonne); updated {
			g.turn = p1Turn
			position_ligne = nouvelle_position_ligne
		}
		return position_colonne, position_ligne
	}else if (g.numeroJoueur == 2) {
		lastXPositionPlayed := -1
		lastYPositionPlayed := -1
		if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			if updated, yPos := g.updateGrid(p1Token, g.tokenPosition); updated {
				g.turn = p1Turn
				lastXPositionPlayed = g.tokenPosition
				lastYPositionPlayed = yPos
				conv := strconv.Itoa(lastXPositionPlayed)
				//envoi de la position jouer à l'autre joueur
				writeClient(g.writer, "p" + "," + conv )
			}
		}
		return lastXPositionPlayed, lastYPositionPlayed
	}

	return -1, -1
}

// Mise à jour de l'état du jeu à l'écran des résultats.
func (g *game) resultUpdate() bool {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if g.joueurPret[0] == false {
			writeClient(g.writer, "pret")
			g.joueurPret[0] = true
			g.nombreJoueursPrets++
		}
	}
	return g.nombreJoueursPrets == 2
}

// Mise à jour de la grille de jeu lorsqu'un pion est inséré dans la
// colonne de coordonnée (x) position.
func (g *game) updateGrid(token, position int) (updated bool, yPos int) {
	for y := globalNumTilesY - 1; y >= 0; y-- {
		if g.grid[position][y] == noToken {
			updated = true
			yPos = y
			g.grid[position][y] = token
			return
		}
	}
	return
}

// Vérification de la fin du jeu : est-ce que le dernier joueur qui
// a placé un pion gagne ? est-ce que la grille est remplie sans gagnant
// (égalité) ? ou est-ce que le jeu doit continuer ?
func (g game) checkGameEnd(xPos, yPos int) (finished bool, result int) {

	tokenType := g.grid[xPos][yPos]

	// horizontal
	count := 0
	for x := xPos; x < globalNumTilesX && g.grid[x][yPos] == tokenType; x++ {
		count++
	}
	for x := xPos - 1; x >= 0 && g.grid[x][yPos] == tokenType; x-- {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// vertical
	count = 0
	for y := yPos; y < globalNumTilesY && g.grid[xPos][y] == tokenType; y++ {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// diag haut gauche/bas droit
	count = 0
	for x, y := xPos, yPos; x < globalNumTilesX && y < globalNumTilesY && g.grid[x][y] == tokenType; x, y = x+1, y+1 {
		count++
	}

	for x, y := xPos-1, yPos-1; x >= 0 && y >= 0 && g.grid[x][y] == tokenType; x, y = x-1, y-1 {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// diag haut droit/bas gauche
	count = 0
	for x, y := xPos, yPos; x >= 0 && y < globalNumTilesY && g.grid[x][y] == tokenType; x, y = x-1, y+1 {
		count++
	}

	for x, y := xPos+1, yPos-1; x < globalNumTilesX && y >= 0 && g.grid[x][y] == tokenType; x, y = x+1, y-1 {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// egalité ?
	if yPos == 0 {
		for x := 0; x < globalNumTilesX; x++ {
			if g.grid[x][0] == noToken {
				return
			}
		}
		return true, equality
	}

	return
}
