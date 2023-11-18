# golens-api

## Get Started

In the root directory, apply the docker compose file which will spin up a postgres db

```bash
docker compose up -d
```

Copy these environment variables to .env

```bash
ALLOW_API_KEY=68daa0a7-b8c0-4735-9328-f8c876aeb0b9
ALLOW_ORIGIN=http://localhost:5173
DATABASE_URL="host=localhost port=5432 user=golensdb password=golensdb dbname=golensdb sslmode=disable"
HOST_PORT=:3000
```

Run the app

```sh
godotenv -f .env go run main.go
```

Run tests suite

```bash
go test ./...
```

or

```bash
ginkgo ./...
```
