# Pastebin

This is text sharing web application complete with a RESTful API.

The implementation was done with TypeScript and Go for the class CMPT315 at MacEwan University.

## Database Utilized

This API uses Docker and the Official Postgres DB image. The version is 13.0, the latest release.

## Building

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

