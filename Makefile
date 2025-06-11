# dir struct will be like:
# proto/
#		serviceA/
#			v1/
#				file.proto
#			v2/
#				file.proto
#		serviceB/
#			v1/
#				file.proto
#			v2/
#				file.proto
# pkg/
#		serviceA/
#			v1/
#				*.pb.go
# cmd/
# 		serviceA/
# 			v1/
# 				main.go
# 				<serviceA binary>


# name = like dir name - 'serviceA' or 'serviceB' etc.
# version = number only - '1' or '2' etc
# ver = shortened name for 'version'

# Path for local binaries with fixed version for protoc
LOCAL_BIN=$(CURDIR)/bin
GEN_GO_VER=1.36.6
GEN_GO_GRPC_VER=1.5.1
# Path for store all proto files
DEFAULT_API_PATH=proto
# Path for store generated output
DEFAULT_GENERATED_PATH=internal/generated

_DEFAULT: help

install_deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v$(GEN_GO_VER)
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v$(GEN_GO_GRPC_VER)

# get_deps:
# 	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
# 	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

# build - build all .cmd/<name>/<version>/*.go files
# require name and ver 
build: check_name
	go build -o cmd/$(name)/v$(version)/$(name) cmd/$(name)/v$(version)/*.go

version=$(ver)
# Directory with proto files for one service
PROTO_INPUT_DIR=$(DEFAULT_API_PATH)/$(name)/v$(version)
# Concrete input proto file
PROTO_INPUT_FILE=$(PROTO_INPUT_DIR)/$(name).proto

# Directory for output generated files
PROTO_OUTPUT_DIR=$(DEFAULT_GENERATED_PATH)/$(name)/v$(version)

# All proto files in single input dir for compilation
PROTO_INPUT_FILES=$(wildcard $(PROTO_INPUT_DIR)/*.proto)

# generate - generate a gRPC stubs from given $(name)/v$(version) folder
generate: check_name check_proto_input_dir_file
	mkdir -p $(PROTO_OUTPUT_DIR)
	protoc -I $(PROTO_INPUT_DIR) \
		--go_out=$(PROTO_OUTPUT_DIR) --go_opt=paths=source_relative \
		--plugin=protoc-gen-go=bin/protoc-gen-go \
		--go-grpc_out=$(PROTO_OUTPUT_DIR) --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
		$(PROTO_INPUT_FILES) 

# check_name - check is vars 'name' and 'version' exists. Exit with 1 if not exists
check_name:
	@if [ -z "$(name)" ]; then echo "[!] Please provide a 'name'"; exit 1 ; fi
	@if [ -z "$(version)" ]; then echo "[!] Please provide a 'version' number"; exit 1 ; fi

# check_proto_input_dir_file - check is vars 'PROTO_INPUT_DIR' and 'PROTO_INPUT_FILE' exists. Exit with 1 if not exists
check_proto_input_dir_file:
	@if [ ! -d "$(PROTO_INPUT_DIR)" ]; then echo "[!] Dir $(PROTO_INPUT_DIR) is not exist"; exit 1 ; fi
	@if [ ! -e "$(PROTO_INPUT_FILE)" ]; then echo "[!] File '$(PROTO_INPUT_FILE)' is not exist"; exit 1; fi

# new - create new dir and first '.proto' file
new:
	@if [ -z "$(name)" ]; then echo "[!] Please provide a name for create"; exit 1 ; fi
	@if [ -z "$(version)" ]; then echo "[!] Please provide a version number for create"; exit 1 ; fi
	mkdir -p $(PROTO_INPUT_DIR)
	touch $(PROTO_INPUT_FILE)

help:
	@echo "'make new name=<string> version=<int>'      : For create a new version of a gRPC proto files"
	@echo "'make generate name=<string> version=<int>' : For generate a gRPC stub"
	@echo "'make build name=<string> version=<int>'    : For compile gRPC implementation"


# This targets may be useless in other project

# generate_sqlc - generate go code from sqlc .yaml config.
# @sqlc_path: path to dir with sqlc .yaml config file. Default "db/sqlc_conf"

generate_sqlc: sqlc_path=db/sqlc_conf
generate_sqlc:
	cd $(sqlc_path) && sqlc generate

build_migrator:
	go build -o migrator cmd/migrator/main.go
