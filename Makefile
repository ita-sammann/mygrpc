.PHONY: proto
proto: vendor
	protoc \
		-I vendor/github.com/grpc-ecosystem/grpc-gateway/	\
		-I vendor/	\
		-I ./grpc/subjects	\
		--go_out=./grpc/subjects --go_opt=paths=source_relative	\
		--go-grpc_out=./grpc/subjects --go-grpc_opt=paths=source_relative	\
		--grpc-gateway_out=./grpc/subjects --grpc-gateway_opt=allow_patch_feature=false,paths=source_relative	\
		--govalidators_out=./grpc/subjects --govalidators_opt=paths=source_relative	\
		./grpc/subjects/service.proto
