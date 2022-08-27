postgres:
	docker run --name postgres12 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root bank_mvp

dropdb:
	docker exec -it postgres12 dropdb bank_mvp

migrateup:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/bank_mvp?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/bank_mvp?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/bank_mvp?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/bank_mvp?sslmode=disable" -verbose down 1

migrateforce:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/bank_mvp?sslmode=disable" force $(N)

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go bank-mvp/db/sqlc Store

.PHONY: postgres, createdb, dropdb, migrateup, migrateup1, migratedown, migratedown1, migrateforce, test, server, mock