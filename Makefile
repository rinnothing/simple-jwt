.PHONY: codegen
codegen:
	go get -tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	
	go tool oapi-codegen -package=schema -generate=server -o=internal/api/schema/server.gen.go api/openapi.yaml
	go tool oapi-codegen -package=schema -generate=types -o=internal/api/schema/types.gen.go api/openapi.yaml
	go mod tidy

.PHONY: generate-key
generate-key:
	go run cmd/generate_key/main.go

.PHONY: test
test:
	go test ./...

.PHONY: migrate
migrate:
	docker-compose down
	docker volume rm simple-jwt_postgres_data 

	docker-compose build
	docker-compose up -d db
	go run cmd/migrate/main.go
	docker-compose down

.PHONY: build
build: codegen
	go build -o server cmd/server/main.go

.PHONY: run
run: codegen
	docker-compose build
	docker-compose up -d

.PHONY: stop
stop:
	docker-compose down
