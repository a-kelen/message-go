#!/bin/bash

docker-compose up --build -d
echo "Try to init cluster..."
docker exec -it db-central cockroach init --insecure
echo "The end!"

 