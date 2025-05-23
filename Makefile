# ───────────────────────────────────────────────────────────
# Project-wide dirs
# ───────────────────────────────────────────────────────────
PROTO_DIR    := proto
PROTO_OUT    := $(PROTO_DIR)/gen
ABS_PROTO_OUT := $(abspath $(PROTO_OUT))

# ───────────────────────────────────────────────────────────
# All .proto inputs:
# ───────────────────────────────────────────────────────────
PROTO_FILES  := $(wildcard $(PROTO_DIR)/*.proto)

# ───────────────────────────────────────────────────────────
# Default target: generate all .pb.go
# ───────────────────────────────────────────────────────────
.PHONY: generate
generate: | $(PROTO_OUT)
	@echo "→ Generating gRPC/Protobuf code into $(ABS_PROTO_OUT)"
	@echo "→ Proto files: $(PROTO_FILES)"
	protoc \
	  --proto_path=$(PROTO_DIR) \
	  --go_out=$(ABS_PROTO_OUT) \
	  --go_opt=paths=source_relative \
	  --go-grpc_out=$(ABS_PROTO_OUT) \
	  --go-grpc_opt=paths=source_relative \
	  $(PROTO_FILES)
	@echo "✓ Generation complete"

# ───────────────────────────────────────────────────────────
# Ensure output directory exists (order-only prerequisite)
# ───────────────────────────────────────────────────────────
$(PROTO_OUT):
	@echo "→ Creating output directory: $@"
	@mkdir -p $@

# ───────────────────────────────────────────────────────────
# Clean generated files
# ───────────────────────────────────────────────────────────
.PHONY: clean
clean:
	@echo "→ Cleaning generated files"
	@rm -rf $(PROTO_OUT)

# ───────────────────────────────────────────────────────────
# Install dependencies
# ───────────────────────────────────────────────────────────
.PHONY: install-plugins
install-plugins:
	@echo "→ Installing protoc plugins"
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
	@echo "✓ Plugins installed to $$(go env GOPATH)/bin"