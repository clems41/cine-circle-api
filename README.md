# cine-circle-api
**(!) WARNING :** you need postgres database for running this API

## PostgresSQL
This project is using PostgresSQL Database.  
For getting one on your machine, you can use docker.  
**(!) (Don't forget to replace value in `<>`)**
```bash
docker run --rm -d --name pg-cine-circle \
  -e POSTGRES_PASSWORD=<postgres_password> \
  -e POSTGRES_USER=<postgres_user> \
  -e PGDATA=/var/lib/postgresql/data/pgdata \
  -v /custom/mount:/var/lib/postgresql/data \
  -p <port_to_use>:5432 \
  postgres:10
```

## Docker
### Build
Simple way of building API on local
```bash
docker build . -t cine-circle-api
```

### Run
Simple way to run API on local (to be used with PostgresSQl instance !)  
**(!) (Don't forget to replace value in `<>`)**
```bash
docker run -d --rm --name cine-circle-api \
  --link pg-cine-circle \
  -e DB_HOST=pg-cine-circle \
  -e DB_USER=<postgres_user> \
  -e DB_PASSWORD=<postgres_password> \
  -e DB_NAME=<database_name> \
  -e DB_PORT=<postgres_port> \
  -p 8080:8080 \
  cine-circle-api
```
If you want to use default values (see below), simply run :
```bash
docker run -d --rm --name cine-circle-api \
  --link pg-cine-circle \
  -e DB_HOST=pg-cine-circle \
  -p 8080:8080 \
  cine-circle-api
```
Default values :
- `DB_HOST = localhost`
- `DB_USER = postgres`
- `DB_PASSWORD = postgres`
- `DB_NAME = cine-circle`
- `DB_PORT = 5432`

## API endpoints
Postman collection can be found into `resources` directory.
Some authentication is needed for this kind of application.  
We'll be using OAuthv2 with JWT token.  
But, until we get these, we are doing authentication by simply adding header into the request with the username.
```bash
--header "username: ${USERNAME}" \
```
### Movies
#### Search for a movie (by title)
```bash
MOVIE_TITLE="inception"
curl --location --request GET "http://localhost:8080/v1/movies?title=${MOVIE_TITLE}"
```
#### Get movie by ID
```bash
MOVIE_ID="tt9484998"
curl --location --request GET "http://localhost:8080/v1/movies/${MOVIE_ID}"
```

### Users (used for app authentication)
#### Create user
**(!) Mandatory fields :**
- fullname (type: string)
- username (type: string) (SQL unique index constraints --> will be used for log in the application)
- email (type: string)
```bash
curl --location --request POST 'http://localhost:8080/v1/users' \
--header 'Content-Type: application/json' \
--data-raw '{
    "FullName": "first last",
    "Email": "test@mail.com",
    "Username": "user1"
}'
```

#### Update existing user
**(!) Mandatory fields :**
- fullname (type: string)
- username (type: string) (SQL unique index constraints --> will be used for log in the application)
- email (type: string)
```bash
USER_ID=2
USERNAME="user1"
curl --location --request PUT "http://localhost:8080/v1/users/${USER_ID}" \
--header "username: ${USERNAME}" \
--header 'Content-Type: application/json' \
--data-raw '{
        "FullName": "fullName2",
        "Username": "username2",
        "Email": "mail2@mail.com"
}'
```

#### Search for users
Not all fields are mandatory, you can use combination with some of thse or none or all at once.
```bash
USERNAME="user1"
FULLNAME="full"
EMAIL="mail"
curl --location --request GET "http://localhost:8080/v1/users?fullname=${FULLNAME}&username=${USERNAME}&email=${EMAIL}"
```

#### Get user info (using ID)
```bash
USER_ID="10"
curl --location --request GET "http://localhost:8080/v1/users/${USER_ID}"
```

#### Check if username already exists (using username)
```bash
USERNAME="user1"
curl --location --request GET "http://localhost:8080/v1/users/${USERNAME}/exists"
```

#### Get all movies rated by user (using username)
```bash
USER_ID="10"
curl --location --request GET "http://localhost:8080/v1/users/${USER_ID}/movies"
```

### Ratings
#### Rate movie for specific (need to be authenticated)
**Fields :**
- Rating (type: float)
- Comment (type: string)
```bash
MOVIE_ID="tt9484998"
USERNAME="user1"
curl --location --request POST "http://localhost:8080/v1/ratings/${MOVIE_ID}" \
--header "username: ${USERNAME}" \
--header 'Content-Type: application/json' \
--data-raw '{
	"Rating": 10,
	"Comment": "Meilleur film de tous les temps !"
}'
```

### Circles
#### Create new circle
```bash
USERNAME="user1"
curl --location --request POST 'http://localhost:8080/v1/circles/' \
--header 'username: ${USERNAME}' \
--header 'Content-Type: application/json' \
--data-raw '{
    "Name": "circle name",
    "Description": "cercle pour les films dactions du dimanche soir"
}'
```

#### Update existing circle
```bash
USERNAME="user1"
CIRCLE_ID="1"
curl --location --request PUT "http://localhost:8080/v1/circles/${CIRCLE_ID}" \
--header "username: ${USERNAME}" \
--header 'Content-Type: application/json' \
--data-raw '{
    "Name": "circle name",
    "Description": "cercle pour les films dactions du dimanche soir"
}'
```

#### Delete existing circle
```bash
USERNAME="user1"
CIRCLE_ID="1"
curl --location --request DELETE "http://localhost:8080/v1/circles/${CIRCLE_ID}" \
--header "username: ${USERNAME}"
```

#### Search for existing circles
```bash
USERNAME="user1"
NAME_TO_SEARCH="my_name"
curl --location --request GET "http://localhost:8080/v1/circles/?name=${NAME_TO_SEARCH}" \
--header "username: ${USERNAME}"
```

#### Add user to existing circle
```bash
USERNAME="user1"
CIRCLE_ID="1"
USER_ID="7"
curl --location --request PUT "http://localhost:8080/v1/circles/${CIRCLE_ID}/${USER_ID}" \
--header "username: ${USERNAME}"
```

#### Remove user from existing circle
```bash
USERNAME="user1"
CIRCLE_ID="1"
USER_ID="7"
curl --location --request DELETE "http://localhost:8080/v1/circles/${CIRCLE_ID}/${USER_ID}" \
--header "username: ${USERNAME}"
```

#### Get all movies rated by users from a circle
```bash
USERNAME="user1"
CIRCLE_ID="1"
SORT_PARAMETER="title:asc"
curl --location --request GET "http://localhost:8080/v1/circles/${}CIRCLE_ID/movies?sort=${SORT_PARAMETER}" \
--header "username: ${USERNAME}"
```

## Update swagger.yaml
(!) API should be running for updating swagger
```bash
cd resources
rm swagger.yaml
wget http://localhost:8080/api/swagger.yaml
```