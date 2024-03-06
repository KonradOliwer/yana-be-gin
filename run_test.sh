#!/bin/bash

container_id=$(docker run -d --name postgres -e POSTGRES_PASSWORD=password -e POSTGRES_USER=user -e POSTGRES_DB=yana -p 5432:5432 postgres)

# Wait for PostgreSQL to become operational
until docker exec $container_id pg_isready -U user; do
  echo "Waiting for PostgreSQL..."
  sleep 1
done

go test ./...

docker kill $container_id
docker rm $container_id