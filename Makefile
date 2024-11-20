install: go.sum
	go install -mod=readonly ./cmd/demod

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

docker-image:
	docker build -t demod:local .

test:
	@go test -mod=readonly -race ./...

ictest:
	cd e2e && go test -race -v .