build:
	go build -o bin/user-management-service
	
serve: build
	ENV=development ./bin/user-management-service

test: build
	ENV=development_test go test -v ./...