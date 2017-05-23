build:
	go build 

install: 
	go install

go-test:
	go test -cover -v $$(go list ./... | grep -v -E 'vendor' )

fmt:
	find -name '*.go' | grep -v .escape | grep -v vendor | xargs -n 1 go fmt
