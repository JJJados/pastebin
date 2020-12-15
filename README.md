# Pastebin

This is text sharing web application complete with a RESTful API.

## Database Utilized

This API uses Docker and the Official Postgres DB image. The version is 13.0, the latest release.

## Building

Installing package dependencies 

```bash
# Go packages
go get -u ./...

# Node packages
npm install
```

Starting your Postgres DB container:

```bash
docker run --name postgres -p 5432:5432 -e POSTGRES_DB=pastebin -e POSTGRES_PASSWORD=password -e POSTGRES_USER=postgres postgres
```

Then applying Goose migrations:

```bash
cd migrations/

goose postgres "user=postgres dbname=pastebin sslmode=disable password=password" up
```

Finally starting the application:

```bash
go run *.go
```

You can then view this application at http://localhost:3333/

## Icons

Supplied by https://material.io/.

