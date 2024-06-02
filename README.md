# Books Management Application

This is a Go-based RESTful API service for managing a collection of books, containerized using Docker.

## Prerequisites

- Docker installed on your machine.

## Build the Docker Image

To build the Docker image, run the following command in the project root directory where the Dockerfile is located:

```sh
docker build -t book-management-app:latest .
```

## Running the Docker Container
```sh
docker run -d -p 1323:1323 --name book-management-app-containe -e PORT=1323 -e DATABASE_DSN="username:password@tcp(127.0.0.1:3306)/database_name?charset=utf8mb4&parseTime=True&loc=Local" --network=host book-management-app:latest
```
(192.168.0.4:6033
### Notes

1. Make sure to replace placeholder values in the `.env` file or the Docker run command with your actual database credentials.
2. The `--network=host` option in the Docker run command allows the container to use the host's network stack, which is useful for local development. In a production environment, it's recommended to configure proper networking and security settings.

By following these instructions, you can containerize and run your Go API service efficiently using Docker.


### API Reference

For detailed API reference, please refer to the [OpenAPI Specification](./openapi.yaml).