BINARY_NAME=platform-installer

.PHONY: build
build:
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin .
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux .
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-windows .


.PHONY: run
run: build
	./${BINARY_NAME}


.PHONY: clean
clean:
	go clean
	rm ${BINARY_NAME}-darwin
	rm ${BINARY_NAME}-linux
	rm ${BINARY_NAME}-windows


.PHONY: test
test:
	go test ./...


.PHONY: test-coverage
test-coverage:
	go test ./... -coverprofile=coverage.out


.PHONY: dep
dep:
	go mod download


.PHONY: vet
vet:
	go vet


.PHONY: lint
lint:
	golangci-lint run --enable-all


.PHONY: prod-build
prod-build:
	make build
	sudo systemctl restart platform-installer
