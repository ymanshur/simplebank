# Simple Bank Service

[![CI Test](https://github.com/ymanshur/simplebank/actions/workflows/test.yml/badge.svg)](https://github.com/ymanshur/simplebank/actions/workflows/test.yml)

Simple Bank Service is perhaps the first project I've undertaken outside of my primary professional focus.

I am committed to maintaining this repository as a resource for my professional development in Go. It is my intention that this repository will serve as a valuable asset for anyone seeking to learn how to develop robust software products using Go best practices.

Please let me know if you have any request or probles about this project by create the issue

Thank you for watch!

## About

The service that I'm going to build is a simple bank. It will provide APIs for the frontend to do following things:

1. Create and manage bank accounts, which are composed of owner’s name, balance, and currency.
2. Record all balance changes to each of the account. So every time some money is added to or subtracted from the account, an account entry record will be created.
3. Perform a money transfer between 2 accounts. This should happen within a transaction, so that either both account's balance are updated successfully or none of them are.

TODO features including:

1. Top-up a balance account through a payment gateway such as Midtrans.
2. Release the balance from an account in booking-action schema.

## Running The Services

1. Clone the repository

    ```shell
    git clone https://github.com/ymanshur/simplebank.git
    ```

2. [Dependencies installation](#dependencies)
3. [Setup infrastructure](#setup-infrastructure)
4. [Run and test your services](#run-your-servies-in-local-machine)

### Dependencies

- [Go](https://golang.org/) v1.23
- [Migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) is database migrations written in Go

    ```shell
    # Versioned installation
    mkdir -p $GOPATH/src/github.com/golang-migrate
    git clone github.com/golang-migrate/migrate $GOPATH/src/github.com/golang-migrate/
    cd $GOPATH/src/github.com/golang-migrate/migrate/cmd/migrate
    git checkout $TAG  # e.g. v4.15.0

    # Replace the postgres build tag with the appropriate database tag(s) for the databases desired
    go build -tags 'postgres' -ldflags="-X main.Version=$(git describe --tags)" -o $GOPATH/bin/migrate $GOPATH/src/github.com/golang-migrate/migrate/cmd/migrate/
    ```

- [SQL Compiler](https://docs.sqlc.dev/en/latest/overview/install.html) that generates type-safe code from SQL

    ```shell
    sudo snap install sqlc
    ```

- [GoMock](https://github.com/golang/mock) is a mocking framework for the Go programming language

    ```shell
    go install github.com/golang/mock/mockgen@v1.6.0
    ```

    Alternatively, use a [maintained fork](https://github.com/uber-go/mock?tab=readme-ov-file#installation) instead

    ```shell
    go install go.uber.org/mock/mockgen@latest
    mockgen -version
    ```

- [DB Docs](https://dbdocs.io/docs) is a simple tool to create web-based documentation for your database.

    ```shell
    npm install -g dbdocs
    dbdocs login
    ```

- [DBML CLI](https://www.dbml.org/cli/#installation)

    DBML (Database Markup Language) is an open-source DSL language designed to define and document database schemas and structures.

    ```shell
    npm install -g @dbml/cli
    ```

    `dbml2sql` is used to convert a DBML file to SQL

    ```shell
    dbml2sql --version
    ```

### Setup infrastructure

Start database PostgreSQL container service:

```shell
make postgres
```

Create `simplebank` database:

```shell
make createdb
```

Run db migration up all versions:

```shell
make migrateup
```

Run db migration down all versions:

```shell
make migratedown
```

### Run your servies in local machine

```shell
make server
```

Test your services:

```shell
make test
```

### Run your servies and the insfrastructures in Docker containers

Following command will run [docker-compose.yaml](docker-compose.yaml) file

```shell
make containes
```

Alternatively, if you already create PostgreSQL dan Redis containers, you just have to run the following command to create only the application container

```shell
make run EMAIL_SENDER_ADDRESS=ymanshur@gmail.com EMAIL_SENDER_PASSWORD=***juc***yrs***f
```

Note: EMAIL_SENDER_ADDRESS and EMAIL_SENDER_PASSWORD are the environment variables that needed to create user

## Documentation

### Database

Update your database design in [docs/db.dbml](docs/db.dbml) the build DB documentation:

```shell
make dbdocs
```

You can access my DB documentation for this project at [this address](https://dbdocs.io/ymanshur/simplebank)

### OpenAPI

Open <http://localhost:8080/swagger> to see APIs documentation based on gRPC Gateway proto definition, see my own at [this address](https://ymanshur.github.io/simplebank/docs/swagger/)

## Code Generation

Generate schema SQL file with DBML CLI:

```shell
make dbschema
```

Generate SQL CRUD with `sqlc`:

```shell
make sqlc
```

Generate DB mock with GoMock:

```shell
make mock
```

Create a new DB migration:

```shell
make migratecreate name=<migration_name>
```

Generate [protobuf](https://grpc.io/docs/languages/go/quickstart/#regenerate-grpc-code) files and update the [API documentation](#openapi)

```shell
make proto
```

## Tips

### How to hit the endpoint using endpoints.http as a playground

1. Install [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) extension
2. To control environment variables add following lines to .vscode/settings.json

    ```json
    "rest-client.environmentVariables": {
        "local": {
            "authority": "localhost:8000",
            "accessToken": "",
            "refreshToken": ""
        },
    },
    ```

3. Copy the returning access and refresh token into environment variables
4. Run the HTTP or Gateway server and follow REST Client documentation to [making request](https://github.com/Huachao/vscode-restclient?tab=readme-ov-file#making-request)

### Control Workspace environment variables

Add following line into .vscode/settings.json

```json
{
    "terminal.integrated.env.linux": {
        "POSTGRES_USER": "",
        "POSTGRES_PASSWORD": "",
        "DB_NAME": "simplebank"
    }
}
```
