build:
	@go build -o bin/GOLANG-CRUD

run: build
	@./bin/GOLANG-CRUD

test:
	@go test -v ./...
