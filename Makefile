
DB_URL=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable

make postgres:
	docker run --name postgres12 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

# migrateup:
#	migrate -path db/migration -database "postgresql://root:d0JvFkurQ7r8lM6GZHnx@simple-bank.csr6cvumpv2h.sa-east-1.rds.amazonaws.com:5432/simple_bank" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

migrations:
	migrate create -ext sql -dir db/migration -seq add_sessions

sqlc:
	docker run --rm -v "C:\Users\User\Desktop\projects\goschool\src\simple-bank:/src" -w /src kjconroy/sqlc generate

test:
	go test -v -cover -short ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/GustavoCielo/simple-bank/db/sqlc Store

docker:
	docker run --name simplebank -p 8080:8080 -e DB_SOURCE="$(DB_URL)" --network bank-network -e GIN_MODE=release simplebank:latest

network:
	docker network create bank-network

proto:
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
    proto/*.proto 
	statik -src=./doc/swagger -dest=./doc

evans:
	evans --host localhost --port 9090 -r repl

redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock migratedown1 migrateup1 docker migrations network proto evans redis