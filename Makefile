.PHONY: network postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 new_migration sqlc api

DB_URL=postgresql://root:secret@localhost:5432/core?sslmode=disable

network:
	docker network create core-network

postgres:
	docker run --name postgres --network core-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

createdb:
	docker exec -it postgres createdb --username=root --owner=root core

dropdb:
	docker exec -it postgres dropdb core

migrateup:
	migrate -path migrations -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path migrations -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path migrations -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path migrations -database "$(DB_URL)" -verbose down 1

new_migration:
	migrate create -ext sql -dir migrations -seq $(name)

sqlc:
	sqlc generate

api:
	go run cmd/main.go
