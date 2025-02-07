#!/bin/bash

# Run make
make

# Run docker compose
PWD=${PWD} docker-compose -f docker-compose.prod.yml build
PWD=${PWD} docker-compose -f docker-compose.prod.yml up -d