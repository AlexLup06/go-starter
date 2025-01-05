#!/bin/bash

# Run make
make

# Run docker compose
docker-compose -f docker-compose.prod.yml build
docker-compose -f docker-compose.prod.yml up -d