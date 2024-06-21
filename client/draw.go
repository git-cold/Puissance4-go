package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"fmt"
)

var blueColor = color.RGBA{100,149,237, 255}
var redColor = color.RGBA{139,0,0, 255}


// Affichage des graphismes à l'écran selon l'état actuel du jeu.
func (g *game) Draw(screen *ebiten.Image) {
	
	screen.Fill(globalBackgroundColor)

	switch g.gameState {
	case titleState:
		g.titleDraw(screen)
	case colorSelectState:
		g.colorSelectDraw(screen)
	case playState:
		g.playDraw(screen)
	case resultState:
		g.resultDraw(screen)
	}

}

// Affichage des graphismes de l'écran titre.
func (g game) titleDraw(screen *ebiten.Image) {
	text.Draw(screen, "Puissance 4 en réseau", largeFont, 90, 150, globalTextColor)
	text.Draw(screen, "Projet de programmation système", smallFont, 105, 190, globalTextColor)
	text.Draw(screen, "Année 2023-2024", smallFont, 210, 230, globalTextColor)

	//Ajout : Afficher le nombre de joueurs connectés
	joueurConnecteTexte := fmt.Sprintf("Joueurs connecté(s) : %d/2", g.nombreDeJoueurConnecte)
	text.Draw(screen, joueurConnecteTexte, smallFont, 210, 320, redColor)
	//

	if g.stateFrame >= globalBlinkDuration/3 {
		text.Draw(screen, "Appuyez sur entrée", smallFont, 210, 500, globalTextColor)
	}
}

// Affichage des graphismes de l'écran de sélection des couleurs des joueurs.
func (g game) colorSelectDraw(screen *ebiten.Image) {
	text.Draw(screen, "Quelle couleur pour vos pions ?", smallFont, 110, 80, globalTextColor)

	line := 0
	col := 0

	//Ajout : Affiche qui doit jouer en fonctions du numero de joueurs
	if g.numeroJoueur == 1 {
		text.Draw(screen, "Vous jouez en premier", smallFont, 110, 120, blueColor)
	} else if g.numeroJoueur == 2 {
		text.Draw(screen, "Vous jouez en deuxième", smallFont, 110, 120, blueColor)
	}
	//

	for numColor := 0; numColor < globalNumColor; numColor++ {

		xPos := (globalNumTilesX-globalNumColorCol)/2 + col
		yPos := (globalNumTilesY-globalNumColorLine)/2 + line
		
		if numColor == g.p1Color {
			vector.DrawFilledCircle(screen, float32(globalTileSize/2+xPos*globalTileSize), float32(globalTileSize+globalTileSize/2+yPos*globalTileSize), globalTileSize/2, globalSelectColor, true)
		}

		//Ajout : Ajout du cercle de l'autre joueur
		if numColor == g.p2Color {
			vector.DrawFilledCircle(screen, float32(globalTileSize/2+xPos*globalTileSize), float32(globalTileSize+globalTileSize/2+yPos*globalTileSize), globalTileSize/2, blueColor, true)
		}
		//

		vector.DrawFilledCircle(screen, float32(globalTileSize/2+xPos*globalTileSize), float32(globalTileSize+globalTileSize/2+yPos*globalTileSize), globalTileSize/2-globalCircleMargin, globalTokenColors[numColor], true)

		col++
		if col >= globalNumColorCol {
			col = 0
			line++
		}
	}
}

// Affichage des graphismes durant le jeu.
func (g *game) playDraw(screen *ebiten.Image) {
	g.drawGrid(screen)

	vector.DrawFilledCircle(screen, float32(globalTileSize/2+g.tokenPosition*globalTileSize), float32(globalTileSize/2), globalTileSize/2-globalCircleMargin, globalTokenColors[g.p1Color], true)

	//Ajout : Affichage du message "À vous de jouer" ou "En attente de l'adversaire"
	var messageJouer string
	if (g.turn == p1Turn && g.numeroJoueur == 1) || (g.turn == p2Turn && g.numeroJoueur == 2) {
		messageJouer = "À vous de jouer"
	}else {
		messageJouer = "En attente de l'adversaire"
	}
	text.Draw(screen, messageJouer, smallFont, 10, 50, redColor)
	//
}

// Affichage des graphismes à l'écran des résultats.
func (g game) resultDraw(screen *ebiten.Image) {
	g.drawGrid(offScreenImage)

	options := &ebiten.DrawImageOptions{}
	options.ColorScale.ScaleAlpha(0.2)
	screen.DrawImage(offScreenImage, options)

	message := "Aucun vainqueur cette fois."
	if g.result == p1wins {
		message = "Victoire !"
	} else if g.result == p2wins {
		message = "Oh non ! Vous avez perdu..."
	}
	text.Draw(screen, message, smallFont, 300, 400, globalTextColor)

	//Ajout
	// Nombre de joueur pret
	joueurPretsText := fmt.Sprintf("Nombre de joueur prêts : %d", g.nombreJoueursPrets)
    text.Draw(screen, joueurPretsText, smallFont, 300, 450, redColor)

	// Affiche si l'on est pret à relancer une partie ou pas
	var statut string
	if g.joueurPret[0] {
		statut = "Prêt"
	} else {
		statut = "Non prêt"
	}
	joueurUnText := fmt.Sprintf("Vous : %s", statut)
	text.Draw(screen, joueurUnText, smallFont, 300, 200, blueColor)

	// Affiche si l'autre joueur est pret à relancer une partie ou pas
	if g.joueurPret[1] {
		statut = "Prêt"
	} else {
		statut = "Non prêt"
	}
	joueurDeuxText := fmt.Sprintf("Autre joueur : %s", statut)
	text.Draw(screen, joueurDeuxText, smallFont, 300, 230, blueColor)

	// Affiche que l'on peut appuyer sur la touche entrée pour rejouer
	if g.stateFrame >= globalBlinkDuration/3 {
        blinkInterval := 30  // intervalle
        if (g.stateFrame/blinkInterval)%2 == 0 {
            text.Draw(screen, "Appuyez sur entrée pour rejouer", smallFont, 210, 500, globalTextColor)
        }
    }
	//
}

// Affichage de la grille de puissance 4, incluant les pions déjà joués.
func (g game) drawGrid(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 0, globalTileSize, globalTileSize*globalNumTilesX, globalTileSize*globalNumTilesY, globalGridColor, true)

	for x := 0; x < globalNumTilesX; x++ {
		for y := 0; y < globalNumTilesY; y++ {

			var tileColor color.Color
			switch g.grid[x][y] {
			case p1Token:
				tileColor = globalTokenColors[g.p1Color]
			case p2Token:
				tileColor = globalTokenColors[g.p2Color]
			default:
				tileColor = globalBackgroundColor
			}

			vector.DrawFilledCircle(screen, float32(globalTileSize/2+x*globalTileSize), float32(globalTileSize+globalTileSize/2+y*globalTileSize), globalTileSize/2-globalCircleMargin, tileColor, true)
		}
	}
}
