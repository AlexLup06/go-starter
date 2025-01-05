#!/usr/bin/env bash

git pull

BUILD_VERSION=$(git rev-parse HEAD)

echo "$(date): Releasing new server version. $BUILD_VERSION"

echo "$(date): Running build..."
make
docker compose -f docker-compose.prod.yml rm -f
docker compose -f docker-compose.prod.yml build

OLD_CONTAINER=$(docker ps -aqf "name=server")
echo "$(date): Scaling server up..."
BUILD_VERSION=$BUILD_VERSION docker compose -f docker-compose.prod.yml up -d --no-deps --scale server=2 --no-recreate server

sleep 60

echo "$(date): Scaling old server down..."
docker container rm -f $OLD_CONTAINER
docker compose -f docker-compose.prod.yml up -d --no-deps --scale server=1 --no-recreate server

echo "$(date): Reloading caddy..."
CADDY_CONTAINER=$(docker ps -aqf "name=caddy")
docker exec $CADDY_CONTAINER caddy reload -c /etc/caddy/Caddyfile
