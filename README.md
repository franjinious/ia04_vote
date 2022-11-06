# IA04-Vote

#### Brief Introduction

Un système simple de vote, comprenant des programmes du serveur et des API de support côté client, basé sur un modèle multi-agent, qui peut initier, participer un vote et obtenir des résultats. Il prend en charge plusieurs algorithmes de vote courants.

---

#### Quick Start

##### A. Go Build

1. Clonez le code source du projet à partir du site officiel utc gitlab.

```bash
git clone https://gitlab.utc.fr/wanhongz/ia04-vote.git
```

2. Modifiez le fichier de configuration dans le répertoire racine du projet, modifiez l'adresse IP et le port de votre serveur ( L'adresse par défaut est **"127.0.0.1:8082"** ).

3. Compilez avec la commande go build.

```go
go build
```

4. Ensuite, vous pouvez trouver l'exécutable **ia04-vote** dans le répertoire racine du projet, l'exécutez

```bash
./ia04-vote
```



##### B. Go Install

Vous pouvez également utiliser la commande **go install** pour installer.

```go
go install -v gitlab.utc.fr/wanhongz/ia04-vote@latest
```

Si tout se passe bien, vous pouvez trouver le fichier exécutable **ia04-vote** dans dossier **bin** sous le chemin **$GOPATH**. 

````shell
$GOPATH/bin/ia04-vote
````



Si tout est normal, vous pouvez voir l'invite de démarrage du serveur

![start](./image/start.png)

---

#### UML

![uml](./image/uml.jpg)

---

#### Client API 

