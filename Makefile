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
	  --go_opt=module=github.com/mo1ein/cloudzy/proto \
	  --go-grpc_out=$(PROTO_OUT) \
	  --go-grpc_opt=module=github.com/mo1ein/cloudzy/proto \
	  $(PROTO_FILES)
	@echo "✓ Generation complete"

# ───────────────────────────────────────────────────────────
# Ensure output directory exists
# ───────────────────────────────────────────────────────────
$(PROTO_OUT):
	mkdir -p $(PROTO_OUT)

# ───────────────────────────────────────────────────────────
# Clean generated files
# ───────────────────────────────────────────────────────────
.PHONY: clean
clean:
	rm -rf $(PROTO_OUT)/*

# ───────────────────────────────────────────────────────────
# Install dependencies
# ───────────────────────────────────────────────────────────
.PHONY: install-plugins
install-plugins:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0