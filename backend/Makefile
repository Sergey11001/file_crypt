sqlc := docker run --rm -v .:/src -w /src sqlc/sqlc:1.27.0

oapi_codegen := go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.3.0

.PHONY: gen-internal-dat-pgrepository
gen-internal-dat-pgrepository:
	$(sqlc) -f internal/pgrepository/sqlc.yml generate

.PHONY: gen-httpengine-openapi
gen-httpengine-openapi:
	$(oapi_codegen) -config internal/httpengine/openapi/oapi-codegen.yml internal/httpengine/openapi/spec.yml
