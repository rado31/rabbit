proto:
	@mkdir -p proto/gen
	protoc \
		--proto_path=proto \
		--go_out=proto/gen      --go_opt=paths=source_relative \
		--go-grpc_out=proto/gen --go-grpc_opt=paths=source_relative \
		client.proto
	@echo "proto generated → proto/gen/"

up:
	docker compose up --build

down:
	docker compose down -v

.PHONY: proto up down
