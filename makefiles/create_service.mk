MKFILE_DIR  := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

ROOT_PACKAGE := github.com/TylerAldrich814/MetaMovie/
ROOT_FILES := $(MKFILE_DIR)rootfiles/
ROOT_DIR := $(MKFILE_DIR)../
STRUCTURE := /{build,cmd,configs,internal/{controller,handler/{grpc,http}},pkg/model}

## ( Internal ): __create_new_service 
#  Creates a new Microservice Go Module.  
#   - creates direcotry ./$(NAME) 
#  Requires:
#   - NAME: Service Name
#   - PORT: The port for the service
__create_new_service:
	@if [ "$(PORT)" = "" ]; then                            \
		echo "PORT is missing, exiting...";                   \
		exit 1;                                               \
	fi
	@if [ -d "$(ROOT_DIR)/$(NAME)" ]; then                  \
		echo "Directory ${NAME} already exists, exiting..";   \
		exit 1;                                               \
	fi                                                     
	@mkdir -p $(ROOT_DIR)$(NAME)$(STRUCTURE)                    && \
	 echo " -- Created new service Directory @ ./${NAME}"       && \
   pushd $(ROOT_DIR)$(NAME)                                   && \
   echo "package main\n\nfunc() main{\n\n}" > ./cmd/main.go   && \
   pushd ./cmd && go mod init $(ROOT_PACKAGE)$(NAME) && popd  && \
	 echo "package grpc\n" > ./internal/handler/grpc/grpc.go    && \
	 echo "package http\n" > ./internal/handler/http/http.go    && \
	 echo "package handler\n\nimport \"errors\"" > ./internal/handler/error.go && \
	 echo "# $(NAME)/configs/base.yaml\n\napi:\n  port:$(PORT)" > ./configs/base.yaml && \
	 cp ${ROOT_FILES}service_config.go ./cmd/config.go
	 tree ./

b:
	@echo "pushd ${ROOT_DIR}${NAME}/"
	@echo "go mod init ${ROOT_PACKAGE}${NAME}"
	@echo " -- Created new go module ${ROOT_PACKAGE}/${service}"
	@echo "package main\n\nfunc main(){\n\n}" > ./cmd/main.go
	@echo "created ./cmd/main.go"
	@tree ./internal/handler/

t:
	@echo "created ./internal/handler/grpc/grpc.go"
	@echo "package http\n\n" > ./internal/handler/http/http.go
	@echo "created ./internal/handler/http/http.go"
	@echo "package handler\n\nimport \"errors\"\n\n" > ./internal/handler/error.go
	@echo "created ./internal/handler/errors.go"
	@echo "# $(service)/configs/base.yaml\n\napi:\n  post:$(PORT)" > ./configs/base.yaml
