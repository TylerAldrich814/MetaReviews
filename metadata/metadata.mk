MKFILE_DIR  := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

.PHONY: compile_metadata docker_metadata_build docker_metadata_push

compile_metadata: ## compile_metadata - Compiles a Metadata Service Linux binary under ./build
	@echo "${MKFILE_DIR}"
	@$(MAKE) __compile  S=metadata DIR=$(MKFILE_DIR)

docker_metadata_build: ## docker_metadata_build - Build Metadata Service Docker Image.
	@$(MAKE) __dk_build S=metadata DIR=$(MKFILE_DIR)

docker_metadata_run:  ## docker_metadata_run - Runs the Metadata Service Docker Image. Optional flag=< -it | -d >:: -it = default
	@$(MAKE) __dk_run S=metadata DIR=$(MKFILE_DIR) PORT=8081 F=$(flag)

docker_metadata_tag: ## docker_metadata_tag - Creates a new Docker Tag: Required VER=<X.X.X>
	@$(MAKE) __dk_tag  S=metadata  VER=$(VER)

docker_metadata_push: ## docker_metadata_push - Updates Metadata Docker Repo: Required: VER=<x.x.x>
	@$(MAKE) __dk_push S=metadata  VER=$(VER)
