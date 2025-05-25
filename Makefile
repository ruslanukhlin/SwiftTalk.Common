PROTO_DIR=proto
PROTO_OUT=gen

PROTO_FILES=$(shell find $(PROTO_DIR) -name "*.proto")

all: generate

generate:
	@protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(PROTO_OUT) \
		--go-grpc_out=$(PROTO_OUT) \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)

clean:
	@rm -rf $(PROTO_OUT)

.PHONY: all generate clean