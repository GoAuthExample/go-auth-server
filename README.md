# Go Auth Server

A lightweight go server that handles OAuth login flows and database retrieval. I was motivated by my internship work to figure out how to create a fullstack approach to open authorization work flows. To help speed up development, I used [goblueprint](https://github.com/Melkeydev/go-blueprint) to create a boilerplate infrastructure. To see the client application visit this [repository](https://github.com/GoAuthExample/go-auth-client-app)


## Features

* Chi Router for routing capabilities
* Open Authorization with Google as an Identity Provider using the Gothic package
* Database service (MySQL) for storing user IdP data
* Dockerized development environement


## Project Structure

* cmd/main.go - entry point that creates the server and gothic store
* pkg/ - defines app packages
  * auth/ - initializes the gothic store for sessions & sets up the google provider
  * database/ - defines DB interactions and methods
  * responses/ - defines reusable types and functions for rendering HTTP responses
  * server/ - route handling and OAuth flows + callbacks
  * scripts/ - contains an SQL file for seeding the database
* docker-compose.yml - spins up the MySQL db container image
* .air.toml - builds the app by running `make build` and provides hot reloading/watching when changes are made
* Makefile - commands to build, run, watch, test the application


---

## Environment Configuration

Values in the `example.env` file are used to run the service. A select number of these variables can be overriden:

* `SESSION_SECRET`: Used to initialize the gothic cookie store (can be any string for now, but this should ideally be a unique randomly generated key)
* `GOOGLE_CLIENT_SECRET` and `GOOGLE_CLIENT_ID`  - These are generated when you create Google OAuth Credentials in Google Cloud Platform


> Make sure to rename example.env to .env



---

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.


### MakeFile

Run build make command with tests

```bash
make all
```

Build the application

```bash
make build
```

Run the application

```bash
make run
```

Create DB container

```bash
make docker-run
```

Shutdown DB Container

```bash
make docker-down
```

DB Integrations Test:

```bash
make itest
```

Live reload the application:

```bash
make watch
```

Run the test suite:

```bash
make test
```

Clean up binary from the last build:

```bash
make clean
```



### Manual Set up

This is if you want to avoid using the commands found in the make file

```bash
docker compose down -v # remove volumes from mysql container
docker compose up      # spin up the mysql container and initialize the DB
```


```bash
# Make sure to set the GO environment variable 
export PATH=$PATH:$HOME/go/bin

# Start the app server
air
```


