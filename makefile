.PHONY: up

MIGRATE_CMD=migrate
MIGRATE_DIR=./migrations
DB_DSN=postgres://user-svc:user-svc@localhost:5433/user-svc?sslmode=disable
DATE=$(shell date +%Y%m%d_%H%M%S)

# Generates mocks for interfaces
INTERFACES_GO_FILES := $(shell find internal -name "interfaces.go")
INTERFACES_GEN_GO_FILES := $(INTERFACES_GO_FILES:%.go=%.mock.gen.go)

generate_mocks: $(INTERFACES_GEN_GO_FILES)
$(INTERFACES_GEN_GO_FILES): %.mock.gen.go: %.go
	@echo "Generating mocks $@ for $<"
	mockgen -source=$< -destination=$@ -package=$(shell basename $(dir $<))

# generate protobuf
PROTOBUF_FILES := $(shell find api/proto -name "*.proto")
PROTOBUF_GEN_FILES := $(PROTOBUF_FILES:api/proto/%.proto=generated/proto/%.pb.go)
generate_protobuf: $(PROTOBUF_GEN_FILES)
	@echo "Generating protobuf files"

$(PROTOBUF_GEN_FILES): generated/proto/%.pb.go: api/proto/%.proto
	@echo "Generating protobuf files $@ for $<"
	buf generate --path $<

generate: api/api.yml generate_mocks generate_protobuf
	mkdir -p generated/api
	oapi-codegen --package api -generate types $< > generated/api/api-types.gen.go

create:
	@$(MIGRATE_CMD) create -ext sql -dir $(MIGRATE_DIR) $(NAME)

up:
	@$(MIGRATE_CMD) -source file://$(MIGRATE_DIR) -database=$(DB_DSN) up

reset:
	@$(MIGRATE_CMD) reset -dir $(MIGRATE_DIR)

refresh: reset up

down:
	@$(MIGRATE_CMD) -source file://$(MIGRATE_DIR) -database=$(DB_DSN) down

status:
	@$(MIGRATE_CMD) status -dir $(MIGRATE_DIR)