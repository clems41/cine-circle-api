# cine-circle-api
WARNING : you need postgres database for running this API

## PostgresSQL
This project is using PostgresSQl Database. For getting one on your machine, you can use docker :
(Don't forget to replace value in `<value>`)
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
(Don't forget to replace value in `<value>`)
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
Default values :
- `DB_HOST = localhost`
- `DB_USER = postgres`
- `DB_PASSWORD = postgres`
- `DB_NAME = cine-circle`
- `DB_PORT = 5432`

## API endpoints
Postman collection can be found into `resources` directory.
### Movies
#### Search for a movie (by title)
```bash
MOVIE_TITLE="inception"
curl --location --request GET 'http://localhost:8080/v1/movies?title=${MOVIE_TITLE}'
```
#### Get movie by ID
```bash
MOVIE_ID="tt9484998"
curl --location --request GET 'http://localhost:8080/v1/movies/${MOVIE_ID}'
```

### Users
#### Create user (used for app authentication)
**Mandatory fields :**
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

#### Get user info (using username)
```bash
USERNAME="user1"
curl --location --request GET 'http://localhost:8080/v1/users/${USERNAME}'
```

#### Check if user already exists (using username)
```bash
USERNAME="user1"
curl --location --request GET 'http://localhost:8080/v1/users/${USERNAME}/exists'
```

#### Get all movies rated by user (using username)
```bash
USERNAME="user1"
curl --location --request GET 'http://localhost:8080/v1/users/${USERNAME}/movies'
```

### Ratings
#### Rate movie for specific (need to be authenticated)
**Fields :**
- Rating (type: float)
- Comment (type: string)
```bash
MOVIE_ID="tt9484998"
USERNAME="user1"
curl --location --request POST 'http://localhost:8080/v1/ratings/${MOVIE_ID}' \
--header 'username: ${USERNAME}' \
--header 'Content-Type: application/json' \
--data-raw '{
	"Rating": 10,
	"Comment": "Meilleur film de tous les temps !"
}'
```
