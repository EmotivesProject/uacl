SLEEP_TIME=0

.PHONY: lint test integration start_test_db

lint:
	golangci-lint run

test:
	go test -v ./...

integration:
	docker-compose -f docker/uacl/docker-compose.test.yml down
	docker-compose -f docker/uacl/docker-compose.test.yml up -d
	sleep $(SLEEP_TIME)
	go test -v -tags=integration ./...
	docker-compose -f docker/uacl/docker-compose.test.yml down

start_test_db:
	docker-compose -f docker/uacl/docker-compose.test.yml down --remove-orphans
	docker-compose -f docker/uacl/docker-compose.test.yml up -d --remove-orphans

integration_dry:
	go test -v -tags=integration ./...