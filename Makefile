DB_URL=postgresql://postgres:postgres@localhost:5432/simplebank?sslmode=disable
#DB_URL=postgresql://postgres:3ail1lQcOm37MhCOBmoU@simplebank.cr6urtxwxgp1.ap-southeast-2.rds.amazonaws.com:5432/simplebank

postgres:
	docker run --name postgres10 -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres:10-alpine

createdb:
	docker exec -it simplebank_postgres createdb --username=postgres --owner=postgres simplebank

dropdb:
	docker exec -it simplebank_postgres dropdb simplebank

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/ymanshur/simplebank/db/sqlc Store

container:
	docker compose -f deployment/docker-compose.yaml up -d --build

release: container

# TOKEN_SYMMETRIC_KEY generator
randhex64:
	openssl rand -hex 64 | head -c $(head)

login-container-registry:
	aws ecr get-login-password | docker login --username AWS --password-stdin 009238256455.dkr.ecr.ap-southeast-2.amazonaws.com

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc test server mock release randhex64