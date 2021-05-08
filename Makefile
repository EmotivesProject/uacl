SLEEP_TIME=0

.PHONY: lint test integration start_test_db get_latest_db

lint:
	golangci-lint run

test:
	go test -v ./...

integration:
	make start_test_db
	sleep $(SLEEP_TIME)
	go test -v -tags=integration ./...
	docker-compose -f docker/uacl/docker-compose.test.yml down

start_test_db:
	make get_latest_db
	docker-compose -f docker/uacl/docker-compose.test.yml down --remove-orphans
	docker-compose -f docker/uacl/docker-compose.test.yml up -d --remove-orphans

integration_dry:
	go test -v -tags=integration ./...

get_latest_db:
	echo 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 | docker login ghcr.io -u imthetom --password-stdin
	docker pull ghcr.io/tombowyerresearchproject/postgres_db:latest
