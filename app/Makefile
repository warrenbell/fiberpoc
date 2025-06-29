BINARY_NAME=fiberpoc

build:
	go build -o ./bin/${BINARY_NAME}_app ./cmd/.
	go build -o ./bin/${BINARY_NAME}_migration ./scripts/migration/.
	go build -o ./bin/${BINARY_NAME}_build_test_postgres ./scripts/build_test_postgres/.

run: build
	./bin/${BINARY_NAME}_app

migup: build
	./bin/${BINARY_NAME}_migration -direct=up

migdown: build
	./bin/${BINARY_NAME}_migration -direct=down

unittest:
	go test -v ./handlers/...

inttest:
	go test -v ./int_testing/...

buildtestpostgres: build
	./bin/${BINARY_NAME}_build_test_postgres


