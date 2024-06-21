package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

// initialise les propriétés du client en fonction du message reçu du serveur.
func initClient(g *game) {
	// Lire son numero de joueur
	message := <-g.channel
	// Convertir le numero de joueur en entier.
	nbJoueur, err := strconv.Atoi(message)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Initialiser le nombre de joueur connecté en fonction du numero de joueur
	g.numeroJoueur = nbJoueur
	if g.numeroJoueur == 1 {
		g.nombreDeJoueurConnecte = 1
	} else if g.numeroJoueur == 2 {
		g.nombreDeJoueurConnecte = 2
	}

	// Afficher le numéro du joueur et le nombre de joueurs connectés.
	fmt.Println("numero joueur : ", g.numeroJoueur)
	fmt.Println("nombre joueur connecté : ", g.nombreDeJoueurConnecte)
}

// writeClient envoie un message au serveur via le writer.
func writeClient(writer *bufio.Writer, message string) {
    _, err := writer.WriteString(message + "\n")
	if err != nil {
		fmt.Println("WriteString error:", err)
		return
	}

	// Vider le buffer du writer.
	err = writer.Flush()
    if err != nil {
		fmt.Println("Flush error:", err)
		return
	}
}

// diviser une chaîne de caractères en utilisant la virgule comme séparateur et renvoie un tableau de sous-chaînes.
func diviseur(read string) (tableau []string) {
	tableau = strings.Split(read, ",")
	return tableau
}

// readClient lit en continu les messages du serveur et les envoie au canal du jeu.
func readClient(reader *bufio.Reader, g *game) {
	for {
		message, err := reader.ReadString(byte('\n'))
        if err!= nil {
            return
        }

		// Supprimer le \n de et envoyer le message au canal du jeu.
        messageConvert := strings.Replace(message,"\n","",-1)
		fmt.Println("message recu: ",messageConvert)
        g.channel <- messageConvert
	}
}

//Fonction pour effacer le canal du jeu et ne garder que le dernier message
func clear(channel chan string) {
	var lastValue string
	var verif = false

	for {
		select {
			case lastValue = <-channel:
				verif = true
			default:
				if verif {
					channel <- lastValue
				}
				return
			}
	}
}


// Peu servir pour d'autre fonctionnalité : viderChannel vide le canal en supprimant toutes les valeurs présentes. 
// func viderChannel(c chan string) {
//     for {
//         select {
//         case <-c:
//             // Ignorer la valeur
//         default:
//             // Sortir de la boucle, si le canal est vide
//             return
//         }
//     }
// } 
