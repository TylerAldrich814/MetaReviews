MKFILE_DIR  := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

.PHONY: compile_movie docker_movie_build docker_movie_push

compile_movie: ## compile_movie - Compiles a Movie Service Linux binary under ./build
	@echo "${MKFILE_DIR}"
	@$(MAKE) __compile  S=movie DIR=$(MKFILE_DIR)

docker_movie_build: ## docker_movie_build - Build Movie Service Docker Image.
	@$(MAKE) __dk_build S=movie DIR=$(MKFILE_DIR)

docker_movie_run:  ## docker_movie_run - Runs the Movie Service Docker Image. Optional flag=< -it | -d >:: -it = default
	@$(MAKE) __dk_run S=movie DIR=$(MKFILE_DIR) PORT=8081 F=$(flag)

docker_movie_tag: ## docker_movie_tag - Creates a new Docker Tag: Required VIR=<X.X.X>
	@$(MAKE) __dk_tag  S=movie  VIR=$(VIR)

docker_movie_push: ## docker_movie_push - Updates Movie Docker Repo: Required: VER=<x.x.x>
	@$(MAKE) __dk_push S=movie VIR=$(VIR)
