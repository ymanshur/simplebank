POSTGRES_VERSION?=12
POSTGRES_USER?=postgres
POSTGRES_PASSWORD?=postgres

DB_NAME?=simplebank
DB_HOST?=localhost
DB_PORT?=5432
DB_SSL?=disable

DB_SOURCE=postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL}

REDIS_VERSION?=8
REDIS_PORT?=6379

postgres:
	docker run --name postgres${POSTGRES_VERSION} --network simplebank-net -p ${DB_PORT}:5432 -r POSTGRES_USER=${POSTGRES_USER} -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -d postgres:${POSTGRES_VERSION}-alpine

createdb:
	docker exec -it postgres${POSTGRES_VERSION} createdb --username=${POSTGRES_USER} --owner=${POSTGRES_USER} ${DB_NAME}

dropdb:
	docker exec -it postgres${POSTGRES_VERSION} dropdb ${DB_NAME}

migrateup:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose down 1

migratecreate:
	migrate create -ext sql -dir db/migration -seq $(name)

sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

server:
	clear
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/ymanshur/simplebank/db/sqlc Store

dbschema:
	dbml2sql --postgres -o docs/schema.sql docs/db.dbml

dbdocs:
	dbdocs build docs/db.dbml

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

build:
	docker build -t simplebank:latest -f deployment/Dockerfile .

# run with necessary environment variables
# to override app.env.example as default values
run: build
	docker run --name simplebank --network simplebank-net \
	-p 8080:8080 \
	-p 9090:9090 \
	-e DEBUG=false \
	-e DB_SOURCE="postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres${POSTGRES_VERSION}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL}" \
	-e GRPC_GATEWAY_SERVER_ADDRESS=0.0.0.0:8080 \
	-e GRPC_SERVER_ADDRESS=0.0.0.0:9090 \
	-e TOKEN_SYMMETRIC_KEY=${shell openssl rand -hex 64 | head -c 32} \
	-e REDIS_ADDRESS=redis${REDIS_VERSION}:${REDIS_PORT} \
	-e EMAIL_SENDER_ADDRESS=${EMAIL_SENDER_ADDRESS} \
	-e EMAIL_SENDER_PASSWORD=${EMAIL_SENDER_PASSWORD} \
	simplebank:latest

containers:
	docker compose -f docker-compose.yaml up -d --build

# TOKEN_SYMMETRIC_KEY generator
randhex64:
	openssl rand -hex 64 | head -c $(head)

login-aws-ecr:
	aws ecr get-login-password | docker login --username AWS --password-stdin 009238256455.dkr.ecr.ap-southeast-2.amazonaws.com

redis:
	docker run --name redis${REDIS_VERSION} --network simplebank-net -p ${REDIS_PORT}:6379 -d redis:${REDIS_VERSION}-alpine

redis-ping:
	docker exec -it redis${REDIS_VERSION} redis-cli ping

# Qodo AI Agent commands
pr-agent:
	@mkdir -p .qodo
	qodo run agents/gh_pull_request.toml create_pr --from-branch=$(from) --target-branch=$(target) 2>&1 | tee .qodo/gh_pull_request_$(shell date +%Y%m%d_%H%M%S).log
