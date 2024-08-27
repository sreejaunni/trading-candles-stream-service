
build:
	@echo "Building Docker image..."
	docker build --platform linux/amd64 --no-cache -t marketplace-app .
start:
	@echo "Streaming Candle sticks.."
	go run cmd/*.go server

reset-db:
	@echo "Resetting db..."
	go run cmd/*.go reset-database


APPLICATION_CONTAINER_NAME := marketplace-application-container
APPLICATION_CONTAINER_ID := $(shell docker ps -a -q -f name=$(APPLICATION_CONTAINER_NAME))

run-server:
	@if [ "$(APPLICATION_CONTAINER_ID)" ]; then \
        		docker rm -f $(APPLICATION_CONTAINER_NAME); \
        	fi
	docker run --name $(APPLICATION_CONTAINER_NAME) -e DB_HOST=host.docker.internal marketplace-app ./app server

MIGRATION_CONTAINER_NAME := marketplace-migration-container
MIGRATION_CONTAINER_ID := $(shell docker ps -a -q -f name=$(MIGRATION_CONTAINER_NAME))

run-migrate:
	@if [ "$(MIGRATION_CONTAINER_ID)" ]; then \
    		docker rm -f $(MIGRATION_CONTAINER_NAME); \
    	fi
	docker run --name $(MIGRATION_CONTAINER_NAME) -e DB_HOST=host.docker.internal marketplace-app ./app reset-database





