.PHONY: build_consul run_consul
MKFILE_DIR  := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
NETWORK_NAME := metamovie-net
CONSUL_CONTAINER := consul-container
CONSUL_NAME := metamovie-consul

run_network: ## Creates a Network for Docker communication
	@docker network create $(NETWORK_NAME) 

build_consul: ## build_consul:  Builds a Consul Docker Image
	@docker build -t $(CONSUL_CONTAINER) -f $(MKFILE_DIR)../Docker/consul.Dockerfile .

run_consul: ## run_consul: Runs the Consul Docker Image
	@docker ps -a --filter "name=$(CONSUL_NAME)" --format '{{.ID}}' | xargs -r docker rm -f
	@docker run -d              \
    --network $(NETWORK_NAME) \
		-p 8500:8500              \
	  -p 8600:8600/udp          \
		--name $(CONSUL_NAME)     \
		$(CONSUL_CONTAINER)
