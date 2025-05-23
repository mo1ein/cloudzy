#!/bin/bash
#
# todo: fix this and write it
## remove old protoc if any
#sudo rm -rf $HOME/ptotoc $HOME/protoc.zip
#sudo rm -rf /usr/local/include/google
#sudo rm -rf /usr/local/bin/protoc
#
## install protoc compiler
#sudo curl -o $HOME/protoc.zip -L https://github.com/protocolbuffers/protobuf/releases/download/v27.2/protoc-27.2-linux-x86_64.zip
#sudo unzip $HOME/protoc.zip -d $HOME/protoc
#sudo mv $HOME/protoc/include/* /usr/local/include/
#sudo mv $HOME/protoc/bin/protoc /usr/local/bin/
#sudo rm -rf $HOME/protoc $HOME/protoc.zip


#
#PROTO_DIR := proto
#PROTO_FILES := $(shell find $(PROTO_DIR) -name "*.proto")
#
#PROTOC_CMD = protoc --go_out=. --go-grpc_out=.
#
#generate: $(PROTO_FILES)
#	@echo "Generation complete."
#
#$(PROTO_FILES):
#	@echo "Compiling $@"
#	protoc --go_out=module=git.o3social.app/backend/packages/golang/protocol-buffers:. --go-grpc_out=module=git.o3social.app/backend/packages/golang/protocol-buffers:.  --proto_path=./proto $(dir $@)$(notdir $@)
#
#install:
#	#bash ./install.sh
#	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31
#	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
#	export PATH="\$$PATH:$(shell go env GOPATH)/bin"
#
#.PHONY: generate $(PROTO_FILES) install

# ───────────────────────────────────────────────────────────
# Project-wide dirs
# ───────────────────────────────────────────────────────────
PROTO_DIR    := proto
PROTO_OUT    := $(PROTO_DIR)/gen

# ───────────────────────────────────────────────────────────
# All .proto inputs:
# ───────────────────────────────────────────────────────────
PROTO_FILES  := $(wildcard $(PROTO_DIR)/*.proto)

# ───────────────────────────────────────────────────────────
# Default target: generate all .pb.go
# ───────────────────────────────────────────────────────────
.PHONY: generate
generate: $(PROTO_OUT) $(PROTO_FILES)
	@echo "→ Generating gRPC/Protobuf code into $(PROTO_OUT)"
	protoc \
	  --proto_path=$(PROTO_DIR) \
	  --go_out=$(PROTO_OUT) \
	  --go_opt=paths=source_relative \
	  --go-grpc_out=$(PROTO_OUT) \
	  --go-grpc_opt=paths=source_relative \
	  $(PROTO_FILES)
	@echo "✓ Generation complete"

# ───────────────────────────────────────────────────────────
# Ensure output directory exists
# ───────────────────────────────────────────────────────────
$(PROTO_OUT):
	mkdir -p $(PROTO_OUT)

# ───────────────────────────────────────────────────────────
# Handy install for protoc plugins
# ───────────────────────────────────────────────────────────
.PHONY: install-plugins
install-plugins:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2


