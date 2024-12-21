.PHONY: __compile __dk_build __dk_run __dk_tag __dk_push

REPO=aldrich814/metamovie_

# ( Internal ): __compile: 
#  compiles the service's linux binary.
#  Requires:
#    - S=<service_name>
#    - DIR=<service_dir>
__compile:
	@echo "->> Compiling $(S) Service Linux Binary @ ./build .."
	@pushd $(DIR)                                   &&  \
		GOOS=linux go build -o ./build/main ./cmd/*go &&  \
		popd                                          &&  \
		echo " -- Finished Compiling $(S)"

# ( Internal ): __dk_build: 
# Builds the Services Docker image.
# Requires:
#   - S=<service_name> 
#   - DIR=<service_Dir>
__dk_build:
	@echo "->> Building $(S) Docker Image.."
	@pushd $(DIR)                                   &&  \
		docker build -t $(S) .                        &&  \
		popd                                          &&  \
		echo " -- Finished Compiling $(S)"

# ( Internal ): __dk_run:
# Runs the Services Docker image.
# Requires:
#   - S=<service_name>
#   - DIR=<service_dir>
#   - PORT=<service_port>
# Optional:
#   - F=<docker run flag( -it | -d )>
__dk_run:
	@echo "->> Running $(S) Docker Image..." &&            \
	  flag="";                                             \
		if [ "$(F)" = "" ]; then                             \
		  flag="-it";                                        \
		elif [ "$(F)" != "-d" ] && [ "$(F)" != "-it" ]; then \
		  echo "Error: Wrong Run flag used '${F}'";          \
		  echo " -  Valid options: '-it' or '-d'";           \
		  exit 1;                                            \
		else                                                 \
		  flag=$(F);                                         \
		fi;                                                  \
		pushd $(DIR) && docker run $$flag -p $(PORT) $(S) && \
		echo " -- $(S) Finished running.."                && \
		popd


# ( Internal ): __dk_tag: 
#  Create a Docker Tag for a Service Docker Image.
#  Requires:
#   - S=<service_name>
#   - VIR=<VERSION>
__dk_tag:
	@if [ "$(VIR)" == "" ]; then                          \
		echo "Version is Mission: i.e. TAG=1.0.0";          \
		exit 1;                                             \
	fi;                                                   \
	docker tag $S $(REPO)$(S):$(VIR)


# ( Internal ): __dk_push: 
#  Pushes the Docker image to the Docker Repository
#  Requires:
#   - S=<service_name>
#   - DIR=<service_dir>
__dk_push:
	@if [ "$(VIR)" == "" ]; then                          \
		echo "Version is Mission: i.e. TAG=1.0.0";          \
		exit 1;                                             \
	fi;                                                   \
	docker push $S $(REPO)$(S):$(VIR)

