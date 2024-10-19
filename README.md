# Cine Circle API

## Fonctionnement

Ce projet fonctionne avec une base de données Postgres pour sauvegarder l'ensemble des données utilisateurs ainsi que les données relatives aux différents médias (films et séries).

L'API [TheMovieDatabase](https://developer.themoviedb.org/reference/intro/getting-started) est utilisée pour récupérer et rechercher des informations sur les médias.

## Démarrer le projet en local avec Docker

### Pré-requis 

**Docker**

Pour fonctionner, ce projet nécessite l'installation de Docker et Docker Compose. [Documentation](https://docs.docker.com/engine/install/).

**Variables d'environnement**

Pour que l'API puisse fonctionner avec Postgres mais également d'autres services externes, il est important de configurer des variables d'environnement.

Pour ce faire, créer un fichier ``.env`` et copier le contenu suivant dedans :
```dotenv
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
DATABASE_NAME=cine-circle-api
THE_MOVIE_DATABASE_API_KEY="XXX"
SMTP_SERVER="smtp.gmail.com"
SMTP_USERNAME="teasycinecircle@gmail.com"
SMTP_PASSWORD="XXX"
```

Remplacez les `XXX` par les valeurs correspondantes :
- `THE_MOVIE_DATABASE_API_KEY` : la clé API associée à votre compte [TheMovieDb](https://developer.themoviedb.org/reference/intro/getting-started), qui peut s'obtenir gratuitement. 
Il suffit de se créer un compte sur le site.
- `SMTP_PASSWORD` : mot de passe du compte Gmail Teasy pour l'utiliser comme server SMTP.

## Lancement

Démarrez l'API avec la base de données PostgreSQL :
```bash
docker compose up -d --build
```

Cela peut prendre quelques minutes lors de la première compilation de l'image Docker.

Une fois démarrée, les logs de l'API seront disponibles via la commande :
```bash
docker compose logs api -f
```

Rendez-vous à la section [Authentification & Swagger](#authentification--swagger) pour découvrir les webservices disponibles.

## Démarrer le projet en local avec Java

### Pré-requis

**Génération de clés**

```bash
# create rsa key pair
openssl genrsa -out keypair.pem 2048

# extract public key
openssl rsa -in keypair.pem -pubout -out public.pem

# create private key in PKCS#8 format
openssl pkcs8 -topk8 -inform PEM -outform PEM -nocrypt -in keypair.pem -out private.pem
```

Copier les clés `public.pem` et `private.pem` dans le dossier `src/main/resources/certs`. Si le dossier n'existe pas, il faut le créer.
La clé `keypair.pem` peut ensuite être supprimée.

**PostgreSQL**

Pour fonctionner, l'API a besoin d'une base de données PostgresSQL.
Plusieurs possibilités :
1. Utiliser Docker Compose : créer un fichier `docker-compose.yaml` et copier le contenu suivant dedans.
```yaml
version: '3.5'

services:
  postgres:
    container_name: local_database
    image: postgres:15
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: cine-circle-api
      PGDATA: /var/lib/postgresql/data/pgdata
    volumes:
       - postgres:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    restart: always
    
volumes:
    postgres:
```
Il suffit ensuite de faire `docker compose up -d` pour déployer la base de données PostgreSQL.

2. Installer Postgres sur la machine et configurer la base de données avec un outil de gestion de base comme `pgAdmin` ou `DBeaver`.

**Variables d'environnement**

Pour que l'API puisse fonctionner avec Postgres mais également d'autres services externes, il est important de configurer des variables d'environnement.

Pour ce faire, créer un fichier ``local.yaml`` et copier le contenu suivant dedans :
```yaml
DB_PASSWORD: "XXX"
DB_USER: "postgres"
DB_HOST: "localhost"
DB_PORT: "5432"
DB_NAME: "cine-circle-api"
THE_MOVIE_DATABASE_API_KEY: "XXX"
SMTP_SERVER: "smtp.gmail.com"
SMTP_USERNAME: "teasycinecircle@gmail.com"
SMTP_PASSWORD: "XXX"
SMTP_PORT_TLS: "587"
```

Remplacez les `XXX` par les valeurs correspondantes :
- `DB_PASSWORD` : mot de passe utiliser pour votre base [PostgreSQL](#postgresql)
- `THE_MOVIE_DATABASE_API_KEY` : la clé API associée à votre compte [TheMovieDb](https://developer.themoviedb.org/reference/intro/getting-started, qui peut sobtenir gratuitement, il suffit simplement de créer un compte.
- `SMTP_PASSWORD` : mot de passe du compte Gmail Teasy pour l'utiliser comme server SMTP.

### Lancement

Utilisez la configuration CineCircleApiapplication déjà créée dans Intellij qui utilise le fichier de variables d'environnement.

Rendez-vous à la section [Authentification & Swagger](#authentification--swagger) pour découvrir les webservices disponibles.

## Authentification & Swagger

L'API est désormais disponible via l'[URL](http://localhost:8080).
Les webservices disponibles sont visibles via le [Swagger](http://localhost:8080/swagger-ui/index.html).

L'API est protégé par un système de token JWT. Un endpoint ``api/v1/auth/sign-in`` permet de récupérer un token JWT à partir du username et mot de passe (Basic Authentication) d'un utilisateur enregistré en base.
Ce token permet ensuite d'accéder à l'ensemble des endpoints sécurisés.

Pour s'authentifier sur le swagger, il suffit de cliquer sur le bouton ``Authorize`` et de renseigner le username et 
mot de passe de l'utilisateur que l'on souhaite utiliser dans la partie `basic` puis cliquer sur `Authorize`.

On peut fermer la fenêtre et aller sur le webservice ``api/v1/auth/sign-in`` > `Try it out` > `Execute`.
Dans la réponse obtenue, il y a un JWT token que l'on peut maintenant utiliser pour les autres webservices.
Il faut de nouveau cliquer sur ``Authorize`` et copier la valeur de ce token dans le champ `value` de la partie `JWT` puis `Authorize`.
On peut désormais fermer la fenêtre et utiliser n'importe quel autres webservices qui nécessitent une authentification.

Répéter la procédure si vous voulez changer d'utilisateur authentifié.

## Système de notifications

Pour que les utilisateurs puissent recevoir des informations en temps réel, un système de notification a été mis en place.
Il permet notamment d'alerter lors de la réception d'une nouvelle recommendation.

Ce système de notification fonctionne via SockJS qui utilise la technologie de websocket.

Le webservice ``/api/v1/notifications/topic`` permet à un utilisateur authentifié de récupérer le nom du topic à utiliser pour sa connection SockJS.
Ce topic est unique par utilisateur, ce qui permet que seul un utilisateur authentifié peut connaître le nom du topic où recevoir ses propres notifications.

Pour chaque nouvelle recommendation reçue par un utilisateur, un message sera envoyé sur le topic correspondant afin qu'il puisse récupérer la notification.
Ce message contiendra toutes les informations d'une recommendations (media concerné, expéditeur, destinataires, commentaire et note).

Voici un exemple de [Frontend](https://github.com/SLFullStackers/SpringAngularWebSocket/blob/master/websocket-frontend/src/app/message.service.ts) qui utilise cette mécanique avec Angular.
