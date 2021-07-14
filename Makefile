.PHONY: proto
proto:
	protoc \
		-I vendor/github.com/grpc-ecosystem/grpc-gateway/	\
		-I vendor/	\
		-I ./grpc/	\
		--go_out=./grpc/ --go_opt=paths=source_relative	\
		--go-grpc_out=./grpc/ --go-grpc_opt=paths=source_relative	\
		--grpc-gateway_out=./grpc/ --grpc-gateway_opt=allow_patch_feature=false,paths=source_relative	\
		--govalidators_out=./grpc/ --govalidators_opt=paths=source_relative	\
		./grpc/subjects/service.proto

.PHONY: build
build:
	mkdir -p bin
	go build -o bin/server ./cmd/server
	go build -o bin/client ./cmd/client
