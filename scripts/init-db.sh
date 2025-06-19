#!/bin/bash
# Wait for YugabyteDB to be ready and create database

echo "Waiting for YugabyteDB to be ready..."
until docker exec golang-boilerplate-db bin/ysqlsh -h localhost -p 5433 -U yugabyte -c "SELECT 1" > /dev/null 2>&1; do
  echo "YugabyteDB is not ready yet..."
  sleep 5
done

echo "YugabyteDB is ready! Creating database..."
docker exec golang-boilerplate-db bin/ysqlsh -h localhost -p 5433 -U yugabyte -c "CREATE DATABASE golang_boilerplate;"
echo "Database created successfully!"