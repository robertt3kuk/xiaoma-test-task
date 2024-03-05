# Test task

## Description

at first we need to start the postgresq in docker or in local

```bash
make post-up
```

it will start the postgresql in docker

then we need to run the migrate command to create the tables in the database

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
make local-up
```

and then you can build and run the app

```bash
go mod tidy
make run
```
