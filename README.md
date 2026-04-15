# Vehicle Dealership Management System Inventory Microservice

[![Go](https://img.shields.io/badge/go-00ADD8.svg?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)

Vehicle inventory microservice for vehicle dealership management.

## Development

To get started, clone this repository and follow these steps to run the Go application in your local environment.

Start all of the Docker containers in the background, you may start in "detached" mode:

```bash
docker-compose up -d
```

The application is executing within a Docker container and is isolated from your local computer. To run various commands against your application use:

```bash
docker-compose exec app {CMD}
```

## Database Migrations

Run all of your outstanding migrations:

```bash
make migrate action=up
```

Roll back the latest migration operation, you may use the rollback Artisan command. This command rolls back the last "batch" of migrations, which may include multiple migration files:

```bash
make migrate action=down
```
