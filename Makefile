.PHONY: lint test integration

lint:
	golangci-lint run

test:
	go test -v ./...

integration:
	docker-compose -f docker/uacl/docker-compose.test.yml up -d
	go test -v -tags=integration ./...
	docker-compose -f docker/uacl/docker-compose.test.yml down