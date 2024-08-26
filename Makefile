
build:
	@echo "Building Docker image..."
	docker build --platform linux/amd64 -t marketplace-app .
start:
	@echo "Streaming Candle sticks.."
	go run cmd/*.go server

reset-db:
	@echo "Resetting db..."
	go run cmd/*.go reset-database


