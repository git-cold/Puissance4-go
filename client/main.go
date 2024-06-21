package main

import (
	"log"
	"net"
	"fmt"
	"bufio"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"golang.org/x/image/font/opentype"
)

// Mise en place des polices d'écritures utilisées pour l'affichage.
func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	smallFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size: 30,
		DPI:  72,
	})
	if err != nil {
		log.Fatal(err)
	}

	largeFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size: 50,
		DPI:  72,
	})
	if err != nil {
		log.Fatal(err)
	}
}

// Création d'une image annexe pour l'affichage des résultats.
func init() {
	offScreenImage = ebiten.NewImage(globalWidth, globalHeight)
}

// Constantes pour l'adresse IP et le port du serveur.
const (
	IP   = "localhost"
	PORT = "8080"
)
// Création, paramétrage et lancement du jeu.
func main() {
	////////////////////////////////
	// Établir une connexion TCP avec le serveur.
	connexion, err := net.Dial("tcp", fmt.Sprintf("%s:%s", IP, PORT))
	if err != nil {
		fmt.Println("La connexion au serveur a échoué :", err)
		return
	}
	defer connexion.Close()
	fmt.Println("Vous venez de vous connecter au serveur :", connexion.RemoteAddr())
	////////////////////////////////


	g := game{}

	////////////////////////////////
	// Initialiser le writer, le channel et le reader pour la communication avec le serveur.
	g.writer = bufio.NewWriter(connexion)
	g.channel = make(chan string,64)
	reader := bufio.NewReader(connexion)

	// Goroutine pour lire les messages du serveur.
	go readClient(reader, &g)

	// Initialiser le client.
	initClient(&g)
	////////////////////////////////

	ebiten.SetWindowTitle("Programmation système : projet puissance 4")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}

}