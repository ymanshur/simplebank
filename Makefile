POSTGRES_VERSION?=10
POSTGRES_USER?=postgres
POSTGRES_PASSWORD?=postgres

DB_NAME?=postgres
DB_HOST?=localhost
DB_PORT?=5432
DB_SSL?=disable

DB_URL=postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL}

postgres:
	docker run --name postgres${POSTGRES_VERSION} -p ${DB_PORT}:5432 -r POSTGRES_USER=${POSTGRES_USER} -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -d postgres:${POSTGRES_VERSION}-alpine

createdb:
	docker exec -it postgres${POSTGRES_VERSION} createdb --username=${POSTGRES_USER} --owner=${POSTGRES_USER} ${DB_NAME}

dropdb:
	docker exec -it postgres${POSTGRES_VERSION} dropdb ${DB_NAME}

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

migratecreate:
	migrate create -ext sql -dir db/migration -seq $(name)

dbdocs:
	dbdocs build docs/db.dbml

dbschema:
	dbml2sql --postgres -o docs/schema.sql docs/db.dbml

sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

server:
	clear
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/ymanshur/simplebank/db/sqlc Store

container:
	docker compose -f docker-compose.yaml up -d --build

# TOKEN_SYMMETRIC_KEY generator
randhex64:
	openssl rand -hex 64 | head -c $(head)

login-aws-ecr:
	aws ecr get-login-password | docker login --username AWS --password-stdin 009238256455.dkr.ecr.ap-southeast-2.amazonaws.com

.PHONY: proto
proto:
	rm -f pb/*.go
	rm -f docs/swagger/*.swagger.json
	rm -f docs/statik/*
	protoc \
	--proto_path=proto \
	--go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
    --openapiv2_out=docs/swagger --openapiv2_opt=allow_merge=true,json_names_for_fields=false \
    --experimental_allow_proto3_optional \
    proto/*.proto
	statik -src=./docs/swagger -dest=./docs

redis:
	docker run --name redis7 -p 6379:6379 -d redis:7-alpine

redis-ping:
	docker exec -it redis redis-cli ping
