MKFILE_DIR  := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

.PHONY: compile_rating docker_rating_build docker_rating_push

compile_rating: ## compile_rating - Compiles a Rating Service Linux binary under ./build
	@echo "${MKFILE_DIR}"
	@$(MAKE) __compile  S=rating DIR=$(MKFILE_DIR)

docker_rating_build: ## docker_rating_build - Build Rating Service Docker Image.
	@$(MAKE) __dk_build S=rating DIR=$(MKFILE_DIR)

docker_rating_run:  ## docker_rating_run - Runs the Rating Service Docker Image. Optional flag=< -it | -d >:: -it = default
	@$(MAKE) __dk_run S=rating DIR=$(MKFILE_DIR) PORT=8081 F=$(flag)

docker_rating_tag: ## docker_rating_tag - Creates a new Docker Tag: Required VIR=<X.X.X>
	@$(MAKE) __dk_tag  S=rating  VIR=$(VIR)

docker_rating_push: ## docker_rating_push - Updates Rating Docker Repo: Required: VER=<x.x.x>
	@$(MAKE) __dk_push S=rating VIR=$(VIR)
