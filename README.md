# cine-circle-api
**(!) ATTENTION :** pour fonctionner, ce projet nécessite l'installation de Docker et Docker Compose. [Documentation](https://docs.docker.com/engine/install/).

## Fonctionnement

Ce projet fonctionne avec une base de données Postgres pour sauvegarder l'ensemble des données utilisateurs ainsi que les données relatives aux différents médias (films et séries).

L'API [TheMovieDatabase](https://developer.themoviedb.org/reference/intro/getting-started) est utilisée pour récupérer et recercher des informations sur les médias.

## Démarrer l'API en local

**Création des variables d'environnements**
1. Créer un fichier vide nommé `.env` à la racine du projet. Ce fichier sera ignoré par git (cf. `.gitignore`).
2. Copier coller le texte suivant dans le fichier `.env` :
```
DATABASE_PASSWORD=XXX
DATABASE_USER=postgres
DATABASE_NAME=cine-circle
RSA_PRIVATE_KEY=
RSA_PUBLIC_KEY=
```
3. Remplacer les `XXX` par le mot de passe que vous souhaitez utiliser.
4. Générer une paire de clé RSA256 utilisée pour l'authentification et la génération des JWT tokens. Pour cela rendez vous sur : `https://travistidwell.com/jsencrypt/demo/`.
5. Encoder les clés en base64 à l'aide de l'outil WEB [BASE64](https://www.base64encode.org/).
5. Copier les clés encodées dans le fichier `.env` au niveau des variables correspondates : `RSA_PUBLIC_KEY` pour la clé publique et `RSA_PRIVATE_KEY` pour la clé privée.

**Démarrer**
Pour démarrer l'API et l base de données Postgres :
```
docker compose up -d
```

L'API est désormais disponible via l'URL "http://localhost:8080".

Les webservices disponibles sont visibles via le [Swagger](https://petstore.swagger.io/?url=http://localhost:8080/swagger.json). 

## Webservices

Les webservices sont protégés par un système d'authentification utilisant des tokens JWT. Seuls 2 webservices sont accessibles sans authentification :
- `/v1/users/sign-in` : permet de générer un token JWT à partir d'un identifiant et un mot de passe fournis en Header (Basic Auth).
- `/v1/users/sign-up` : permet de se créer un compte pour ensuite pouvoir se connecter.

Tous les autres webservices nécessite l'utilisation d'un token JWT qui doit être ajouté en Header de la requête avec le nom `Authorization` et en contenu `Bearer <jwt_token>` où `<jwt_token>` est le token généré via `/v1/users/sign-in`.

Les exemples donnés ci-dessous utilise cURL, qui est un outil permettant de faire des requêtes HTTP. 
On peut l'installer sur sa machine pour lancer les requêtes ou alors utiliser directement une version WEB [Reqbin cURL](https://reqbin.com/curl).
Pour utiliser ce site WEB avec des requêtes en local, il faut installer le plugin Chrome [ReqBin HTTP Client
](https://chrome.google.com/webstore/detail/reqbin-http-client/gmmkjpcadciiokjpikmkkmapphbmdjok/related).

### Authentification

**Création de compte**

```bash
curl --location --request POST "http://localhost:8080/v1/users/sign-up" \
--header 'Content-Type: application/json' \
--data-raw '{
        "email": "monemail@gmail.com",
        "firstName": "John",
        "lastName": "Doe",
        "password": "password",
        "username": "johndoe"
}'
```

**Connexion/Génération de token**

La connexion utilise un système de [Basic Authentification](https://developer.mozilla.org/fr/docs/Web/HTTP/Headers/Authorization#directives).
Elle utilise un header `Authorization` dont le contenu est `Basic <credentials>` où `<credentials>` est le username et le mot de passe encodé en base64.
Exemple avec username `johndoe` et password `password`, la phrase à encoder est `johndoe:password`, ce qui donne `am9obmRvZTpwYXNzd29yZA==`.

Il est possible de générer ce Header directement via le site WEB [Blitter](https://www.blitter.se/utils/basic-authentication-header-generator/).

```bash
curl --location --request POST "http://localhost:8080/v1/users/sign-in" \
--header 'Authorization: Basic am9obmRvZTpwYXNzd29yZA=='
```

Vous obtenez ainsi un JWT token dans la réponse, par exemple `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhbnkiLCJleHAiOjE2OTc1NzU3MTUsImlzcyI6ImNpbmUtY2lyY2xlLWFwaSIsInN1YiI6IntcIklkXCI6MSxcIlJvbGVcIjpcIlwifSJ9.Kh3EhPRg1WDYLRqI4PWtFMWcYIJ7CSE2vgnDJaZWBcdh7LRY7BnKwv3U2Wf2dWoRaDnFZpnWilkg6tZ0mudCkoSuP29mWSq4CBr0kDxWk1FIr6Pnbm5Oap9Ylg89NZpuNGdZpt-wyaOt64SrGKm9LEzbVRFJC_TpMo9W4BmV6z4`.
Vous pouvez désormais requêter tous les autres webservices en utilisant ce token, voir les exemples après.
Il faudra ajouter un Header `Authorization` avec comme contenu `Bearer <jwt_token>`.

### Medias

**Rechercher un média**

```bash
curl --location --request GET "http://localhost:8080/v1/medias?keyword=inception" \
--header 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhbnkiLCJleHAiOjE2OTc1NzU3MTUsImlzcyI6ImNpbmUtY2lyY2xlLWFwaSIsInN1YiI6IntcIklkXCI6MSxcIlJvbGVcIjpcIlwifSJ9.Kh3EhPRg1WDYLRqI4PWtFMWcYIJ7CSE2vgnDJaZWBcdh7LRY7BnKwv3U2Wf2dWoRaDnFZpnWilkg6tZ0mudCkoSuP29mWSq4CBr0kDxWk1FIr6Pnbm5Oap9Ylg89NZpuNGdZpt-wyaOt64SrGKm9LEzbVRFJC_TpMo9W4BmV6z4'
```

**Récupérer les informations d'un média**

```bash
curl --location --request GET "http://localhost:8080/v1/medias/18" \
--header 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJhbnkiLCJleHAiOjE2OTc1NzU3MTUsImlzcyI6ImNpbmUtY2lyY2xlLWFwaSIsInN1YiI6IntcIklkXCI6MSxcIlJvbGVcIjpcIlwifSJ9.Kh3EhPRg1WDYLRqI4PWtFMWcYIJ7CSE2vgnDJaZWBcdh7LRY7BnKwv3U2Wf2dWoRaDnFZpnWilkg6tZ0mudCkoSuP29mWSq4CBr0kDxWk1FIr6Pnbm5Oap9Ylg89NZpuNGdZpt-wyaOt64SrGKm9LEzbVRFJC_TpMo9W4BmV6z4'
```