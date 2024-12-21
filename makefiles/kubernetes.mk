.PHONY: kube_apply kube_delete_pod 

ROOT_DIR := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

## Doesn't work..?
# kube_apply: ## kube_apply - Taking in the service name, runs kubectl apply on ./<service>/kubernetes-deployment.yaml: Required pod=<service name>
# 	@if [ "${pod}" = "" ]; then                                        \
# 		echo "pod name is required";                                     \
# 		exit 1;                                                          \
# 	fi
# 	@if [ ! -d "$(pod)" ]; then                                        \
# 		echo "$(pod) isn't a valid Service";                             \
# 	  exit 1;                                                          \
# 	fi
# 	@if [ ! -f "$(ROOT_DIR)/$(pod)/kubernetes-deployment.yaml" ]; then \
# 		echo "$(ROOT_DIR)/$(pod)/kubernetes-deployment.yml is missing";  \
# 	  exit 1;                                                          \
# 	fi;
# 	@kubctl apply -f $(ROOT_DIR)/$(pod)/kubernetes-deployment.yaml


kube_delete_pod: ## kube_delete_pod - Deletes app=<SERVICE> Kuberenetes Pod.
	@if [ "${pod}" = "" ]; then     \
		echo "pod name is required";  \
		exit 1;                       \
	fi
	@kubectl delete pod -l app=${pod}


