test-database:
	@ echo starting tests on database 
	@ go test ./internals/database/... -cover -fullpath -v -benchmem 

test-handlers:
	@ echo starting tests on handlers 
	@ go test ./internals/handlers/... -cover -fullpath -benchmem 

build:
	@ echo building API
	@ go build -o ./dist/app ./cmd/main.go 
	@ echo API built

start: build 
	@ echo starting API
	@ ./dist/app