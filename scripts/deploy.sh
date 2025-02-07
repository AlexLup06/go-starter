#!/usr/bin/env bash

git pull

BUILD_VERSION=$(git rev-parse HEAD)

echo "$(date): Releasing new server version. $BUILD_VERSION"

echo "$(date): Running build..."
make
sudo docker compose -f docker-compose.prod.yml rm -f
sudo docker compose -f docker-compose.prod.yml build

OLD_CONTAINER=$(sudo docker ps -aqf "name=server")
echo "$(date): Scaling server up..."
sudo docker compose -f docker-compose.prod.yml up -d --no-deps --scale server=2 --no-recreate server

sleep 30

echo "$(date): Scaling old server down..."
sudo docker container rm -f $OLD_CONTAINER
sudo docker compose -f docker-compose.prod.yml up -d --no-deps --scale server=1 --no-recreate server --remove-orphans

echo "$(date): Reloading caddy..."
CADDY_CONTAINER=$(sudo docker ps -aqf "name=caddy")
sudo docker exec $CADDY_CONTAINER caddy reload -c /etc/caddy/Caddyfile
