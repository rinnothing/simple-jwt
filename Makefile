.PHONY: codegen
codegen:
	go get -tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	
	go tool oapi-codegen -package=schema -generate=server -o=internal/api/schema/server.gen.go api/openapi.yaml
	go tool oapi-codegen -package=schema -generate=types -o=internal/api/schema/types.gen.go api/openapi.yaml
	go mod tidy