postgres:
	docker run --name postgres_simplebank -p 5431:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:latest

startpostgres:
	docker start postgres_simplebank

createdb:
	docker exec -it postgres_simplebank createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres_simplebank dropdb simple_bank

migrateup:
	 migrate -path db/migration -database "postgresql://root:secret@localhost:5431/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5431/simple_bank?sslmode=disable" -verbose down

sqlc: 
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres dropdb createdb migrateup migratedown sqlc test

