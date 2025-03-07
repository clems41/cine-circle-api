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
DATABASE_USER="postgres"
DATABASE_PASSWORD="postgres"
DATABASE_NAME="huco-db"
THE_MOVIE_DATABASE_API_KEY="XXX"
SMTP_SERVER="smtp.hostinger.com"
SMTP_USERNAME="noreply@hucoapp.io"
SMTP_PASSWORD="XXX"
```

Remplacez les `XXX` par les valeurs correspondantes :
- `THE_MOVIE_DATABASE_API_KEY` : la clé API associée à votre compte [TheMovieDb](https://developer.themoviedb.org/reference/intro/getting-started), qui peut s'obtenir gratuitement. 
Il suffit de se créer un compte sur le site.
- `SMTP_PASSWORD` : mot de passe du compte mail noreply@hucoapp.io.

## Lancement

### Démarrer le projet en local avec Docker (recommandé)

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

### Démarrer le projet en local avec Java (debug)

#### Pré-requis

**PostgreSQL**

Pour fonctionner, l'API a besoin d'une base de données PostgresSQL. 
Vous pouvez utiliser celle paramétrée dans le `docker-compose.yaml`.
Pour ce faire :

```bash
docker compose up -d database
```

**Variables d'environnement**

Pour que l'API puisse fonctionner avec Postgres mais également d'autres services externes, il est important de configurer des variables d'environnement.

Pour ce faire, créer un fichier ``local.yaml`` et copier le contenu suivant dedans :
```yaml
DB_USER: "postgres"
DB_PASSWORD: "postgres"
DB_HOST: "localhost"
DB_NAME: "huco-db"
THE_MOVIE_DATABASE_API_KEY: "XXX"
SMTP_SERVER: "smtp.hostinger.com"
SMTP_USERNAME: "noreply@hucoapp.io"
SMTP_PASSWORD: "XXX"
SMTP_PORT_TLS: "465"
```

Remplacez les `XXX` par les valeurs correspondantes :
- `THE_MOVIE_DATABASE_API_KEY` : la clé API associée à votre compte [TheMovieDb](https://developer.themoviedb.org/reference/intro/getting-started, qui peut sobtenir gratuitement, il suffit simplement de créer un compte.
- `SMTP_PASSWORD` : mot de passe du compte mail noreply@hucoapp.io.

**Plugin EnvFile pour Intellij**

Pour que la configuration fonctionne avec le fichier d'environnement précédemment créé, il faut installer le plugin `EnvFile` dans votre IDE.

### Lancement

Utilisez la configuration CineCircleApiapplication déjà créée dans Intellij qui utilise le fichier de variables d'environnement.
Si vous n'utilisez pas Intellij, il faut démarrer le projet en ajoutant les variables du fichier `local.yaml` en tant que variables d'environnement avant de démarrer l'API.

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

Côté front, il faut utiliser le websocket mis en place sur l'url ``ws://localhost:5000/api/v1/notifications/websocket`` et en précisant le topic reçu précédemment ``/topic/XXX`` lors de la subscription.

Pour chaque nouvelle recommendation reçue par un utilisateur, un message sera envoyé sur le topic correspondant afin qu'il puisse récupérer la notification.
Ce message contiendra toutes les informations d'une recommendations (media concerné, expéditeur, destinataires, commentaire et note).

Voici un exemple de [Frontend](https://github.com/SLFullStackers/SpringAngularWebSocket/blob/master/websocket-frontend/src/app/message.service.ts) qui utilise cette mécanique avec Angular.
