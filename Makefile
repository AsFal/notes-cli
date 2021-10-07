.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint:
	golint ./...

.PHONY: lint
format:
	go fmt ./...

.PHONY: verify
verify: vet lint

.PHONY: test-unit
test-unit:
	@mkdir -p .coverage
	go test -coverprofile .coverage/cover.out ./...
	@go tool cover -html=.coverage/cover.out -o .coverage/cover.html

bin-%:
	@mkdir -p bin
	@GOOS=$* GOARCH=arm64 go build -o bin/notes-cli-$*-arm64 
	@echo "Generated bin/notes-cli-$*-arm64"

.PHONY: bin
bin: bin-linux bin-windows bin-darwin

.PHONY: install
install: bin-darwin
	@mv ./bin/notes-cli-darwin-arm64 /usr/local/bin/notes

.PHONY: run
run: bin-darwin
	@./bin/notes-cli-darwin-arm64 
