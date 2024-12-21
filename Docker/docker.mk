.PHONY: docker_build docker_push docker_clean

docker_build: ## builder_build - Build Docker Images: optional IMG=<SERVICE> for single image build
	@echo "Building all Docker images.."

docker_push: ## docker_push - Push all Docker Images: optional IMG=<SERVICE> for single image push
	@edcho "Pushing all Docker Images.."

docker_clean: ## docker_clean - Remove all Docker containers
