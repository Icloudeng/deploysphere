BINARY_NAME=platform-installer
MAIN_PACKAGE_PATH := ./cmd/${BINARY_NAME}


.PHONY: build
build:
	GOARCH=amd64 GOOS=darwin go build -o ./bin/${BINARY_NAME}-darwin ${MAIN_PACKAGE_PATH}
	GOARCH=amd64 GOOS=linux go build -o ./bin/${BINARY_NAME}-linux ${MAIN_PACKAGE_PATH}
	GOARCH=amd64 GOOS=windows go build -o ./bin/${BINARY_NAME}-windows ${MAIN_PACKAGE_PATH}


.PHONY: build-linux
build-linux:
	GOARCH=amd64 GOOS=linux go build -o ./bin/${BINARY_NAME}-linux ${MAIN_PACKAGE_PATH}


.PHONY: dev
dev:
	go run ${MAIN_PACKAGE_PATH}



.PHONY: service-restart
service-restart:
	sudo systemctl restart platform-installer



.PHONY: journal
journal:
	sudo journalctl -fu platform-installer



.PHONY: run
run: build
	./bin/${BINARY_NAME}


.PHONY: clean
clean:
	go clean
	rm ./bin/${BINARY_NAME}-darwin
	rm ./bin/${BINARY_NAME}-linux
	rm ./bin/${BINARY_NAME}-windows



.PHONY: pull
pull:
	git pull origin main




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
prod-build: build-linux service-restart


.PHONY: prod-build-linux
prod-build-linux: build-linux service-restart


.PHONY: queue-status
queue-status:
	go run ./cmd/queue-status


.PHONY: ansible-logs
ansible-logs:
	tail -n -0 -f infrastructure/provisioner/logs/ansible_log.txt


.PHONY: yamllint
yamllint:
	yamllint "infrastructure/provisioner/scripts/**/*.(yaml|yml)"



