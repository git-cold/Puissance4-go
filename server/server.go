package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"strings"
)

// Constantes pour l'adresse IP et le port du serveur.
const (
    IP   = "localhost" 
    PORT = "8080"       
)

var (
	joueursConnectes = 0

	globalWriter_un *bufio.Writer
	globalWriter_deux *bufio.Writer

	globalReader_un *bufio.Reader
	globalReader_deux *bufio.Reader
)

// handleClient gère la connexion d'un client en fonction du nombre de joueurs connectés.
func handleClient(connexion net.Conn) {
	joueursConnectes++

	if joueursConnectes == 1 {
		globalWriter_un = bufio.NewWriter(connexion)
		globalReader_un = bufio.NewReader(connexion)

		// Informe le premier client qu'il est connecté avec son numero de joueur.
		writeServer(globalWriter_un,"1")

		fmt.Println("Premier client connecté depuis", connexion.RemoteAddr())
	} else if joueursConnectes == 2 {
		globalWriter_deux = bufio.NewWriter(connexion)
		globalReader_deux = bufio.NewReader(connexion)
		
		// Informe le deuxième client qu'il est connecté avec son numero de joueur.
		writeServer(globalWriter_deux,"2")
		// informer le premier client de la connexion
		writeServer(globalWriter_un,"connexion")

		fmt.Println("Deuxième client depuis", connexion.RemoteAddr())
		fmt.Println("La partie peut commencer")
	}
}

// fonction qui gère la création du serveur, l'attente des connexions des clients, et la communication entre eux.
func main() {
	// Création du listener
    fmt.Println("Lancement du serveur ...")
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", IP, PORT))
	if err != nil {
		fmt.Println("Erreur lors de la création du serveur :", err)
		return
	}
	defer listener.Close()

	// Attente de la connexion du premier client
	connexion_un, err := listener.Accept()
	if err != nil {
		fmt.Println("Accept error:", err)
		return
	}
	defer connexion_un.Close()
	handleClient(connexion_un)

	// Attente de la connexion du deuxième client
	connexion_deux, err := listener.Accept()
	if err != nil {
		fmt.Println("Accept error:", err)
		return
	}
	defer connexion_deux.Close()
	handleClient(connexion_deux)

	// Initialisation de la WaitGroup pour que la goroutine principale n'arrête pas les deux goroutines lancées.
	var wg sync.WaitGroup
	wg.Add(2)

	// Goroutine pour la communication du premier client vers le deuxième
    go func() {
        defer wg.Done()
        for {
			message, err := globalReader_un.ReadString(byte('\n'))
			if err != nil {
				continue
			}
			
			// Supprimer le caractère de saut de ligne et transmettre le message au deuxième client.
			messageConvert := strings.Replace(message,"\n","",-1)
			writeServer(globalWriter_deux, messageConvert)
        }
    }()

	// Goroutine pour la communication du deuxième client vers le premier
	go func() {
        defer wg.Done()
        for {
			message, err := globalReader_deux.ReadString(byte('\n'))
			if err != nil {
				continue
			}
			
			// Supprimer le caractère de saut de ligne et transmettre le message au deuxième client.
			messageConvert := strings.Replace(message,"\n","",-1)
			writeServer(globalWriter_un, messageConvert)
        }
    }()

	// Attente pour que les deux goroutines ne se terminent pas
	wg.Wait()
}

// writeServer envoie un message à un client via le writer.
func writeServer(writer *bufio.Writer, message string) {
    _, err := writer.WriteString(message + "\n")
	if err != nil {
		fmt.Println("WriteString error:", err)
		return
	}

	// Supprimer le \n de et envoyer le message au canal du jeu.
	err = writer.Flush()
    if err != nil {
		fmt.Println("Flush error:", err)
		return
	}
}

// diviseur divise une chaîne de caractères en utilisant la virgule comme séparateur et renvoie un tableau de sous-chaînes.
func diviseur(read string) (tableau []string) {
	tableau = strings.Split(read, ",")
	return tableau
}