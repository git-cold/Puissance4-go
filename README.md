# Projet Puissance 4 Multijoueur en GO

<div style="display: flex; align-items: center;">
    <img src="readmefile/go.png" alt="GO logo" style="width: 200px;">
    <img src="readmefile/ebitengine.png" alt="Ebitengine logo" style="width: 200px;">
</div>

[![Vidéo de démonstration du jeu](readmefile/Puissance4-go.png)](readmefile/Puissance4.mp4)

## Lancement du Projet

Pour démarrer le projet, lancez le serveur et les deux clients.

### Serveur

Dans le dossier `/serveur`, exécutez la commande :

```sh
go run server.go
```

### Client

Dans le dossier `/client`, exécutez la commande suivante, pour obtenir un exécutable permettant de lancer l'application (assurer vous d'avoir la version 1.21.3 sinon changer la version dans le fichier `/client/go.mod`) :

```sh
go build
```

## Fonctionnement d'Ebitengine

![Schéma du fonctionnement d'ebitengine](readmefile/schema1.png)

## Fonctionnement du Client-Serveur

![Schéma fonctionnement du Client-Serveur](readmefile/schema2.png)
