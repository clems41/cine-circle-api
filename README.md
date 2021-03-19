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
### localhost:8080/v1/movie
- `GET / ?id=<movieId>` get movie details by ID (based on IMDb IDs)
- `GET / ?title=<title>` get movie details based on title

### localhost:8080/v1/user
- `POST /` create user
- `GET / ?id=<userId>` get user details by ID