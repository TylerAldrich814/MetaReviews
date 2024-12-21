include Docker/docker.mk
include makefiles/create_service.mk
include makefiles/services.mk
include makefiles/consul.mk
include makefiles/kubernetes.mk
include metadata/metadata.mk
include rating/rating.mk
include movie/movie.mk

.PHONY: help

## Was gonna make a service creation tool. But, I feel task would be better suited for Go.
# new_service: ## new_service - Creates a New Go Service direcotry & file structure: requires name=<NAME> port=<XXXX>
# 	@echo "Creating a new Service $(name)"
# 	@$(MAKE) __create_new_service NAME=$(name) PORT=$(port)

help:
	@echo "Available Commands:"
	@echo $(MAKEFILE_LIST)
	@grep -E '^[a-zA-Z_-]+:.*##' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf " - %0s\n", $$2}'
