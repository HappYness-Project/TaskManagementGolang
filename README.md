# TaskManagementAPI
Task Management API is used for handling tasks and containers for the specific user groups.
This project is designed to be run using Docker and Docker Compose. Follow the instructions below to get started.

## Prerequisites

Make sure you have the following installed on your system:

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- Makefile setup

## Getting Started


### Run Docker Container
Inside the root of the project, You should run
```sh
make start
```
Above comment will create the postgres database within the docker and create the tables and sample data.

If you want to stop/remove the containers,
```sh
make down
```

Rebuild command
```sh
make rebuild-docker
```

