.PHONY: codegen
codegen:
	go get -tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	npx swagger2openapi -y -o /dev/stdout api/swagger.yaml | \
		go tool oapi-codegen -package=schema -generate=server -o=internal/api/schema/server.gen.go /dev/stdin
	npx swagger2openapi -y -o /dev/stdout api/swagger.yaml | \
		go tool oapi-codegen -package=schema -generate=types -o=internal/api/schema/types.gen.go /dev/stdin
	go mod tidy