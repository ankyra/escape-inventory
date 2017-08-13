build: update-schema
	go build 

install: update-schema
	go install

go-test: update-schema
	go test -cover -v $$(go list ./... | grep -v -E 'vendor' )

fmt:
	find -name '*.go' | grep -v .escape | grep -v vendor | xargs -n 1 go fmt

update-schema:
	go-bindata -prefix dao/sqlite/schemas/ -o dao/sqlite/schema.go -pkg sqlite dao/sqlite/schemas
	go-bindata -prefix dao/postgres/schemas/ -o dao/postgres/schema.go -pkg postgres dao/postgres/schemas
