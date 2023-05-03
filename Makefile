DB_URL=postgresql://postgres:postgres@localhost:5432/simplebank?sslmode=disable
#DB_URL=postgresql://postgres:3ail1lQcOm37MhCOBmoU@simplebank.cr6urtxwxgp1.ap-southeast-2.rds.amazonaws.com:5432/simplebank

PG_VERSION ?= 10

.PHONY: postgres
postgres:
	docker run --name postgres${PG_VERSION} -p 5432:5432 -e POSTGRES_PASSWORD=postgres -d postgres:${PG_VERSION}-alpine

.PHONY: createdb
createdb:
	docker exec -it simplebank-postgres createdb --username=postgres --owner=postgres simplebank

.PHONY: dropdb
dropdb:
	docker exec -it simplebank-postgres dropdb simplebank

.PHONY: migrateup
migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

.PHONY: migrateup1
migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

.PHONY: migratedown
migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

.PHONY: migratedown1
migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

.PHONY: new_migration
new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

db_schema:
	dbml2sql --postgres -o docs/schema.sql docs/db.dbml

.PHONY: sqlc
sqlc:
	sqlc generate

.PHONY: test
test:
	go test -v -cover -short ./...

.PHONY: server
server:
	go run main.go

.PHONY: mock
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/ymanshur/simplebank/db/sqlc Store

.PHONY: container
container:
	docker compose -f docker-compose.yaml up -d --build

.PHONY: randhex64
# TOKEN_SYMMETRIC_KEY generator
randhex64:
	openssl rand -hex 64 | head -c $(head)

.PHONY: login-container-registry
login-container-registry:
	aws ecr get-login-password | docker login --username AWS --password-stdin 009238256455.dkr.ecr.ap-southeast-2.amazonaws.com

.PHONY: proto
proto:
	rm -f pb/*.go
	rm -f docs/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
    --openapiv2_out=docs/swagger --openapiv2_opt=allow_merge=true,json_names_for_fields=false \
    --experimental_allow_proto3_optional \
    proto/*.proto
	statik -src=./docs/swagger -dest=./docs

.PHONY: redis
redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine

.PHONY: ping-redis
ping-redis:
	docker exec -it redis redis-cli ping