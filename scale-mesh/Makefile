# One of: breadth, breadth-sink, circle, circle-callback, depth, depth-sink, hourglass
MESH_TYPE ?= circle
NUMBER_NAMESPACES ?= 1
NUMBER_APPS ?= 5
NUMBER_VERSIONS ?= 2
NUMBER_SERVICES ?= 5

# Only used when creating SMM resources within a Maistra environment
CONTROL_PLANE_NAMESPACE ?= istio-system

DORP ?= docker
KUBECONFIG ?= ${HOME}/.kube/config
KUBECTL ?= kubectl

build:
	docker build -t quay.io/kiali/scale-mesh-demo:latest $(CURDIR)

install: .set-state-install .execute

uninstall: .set-state-uninstall .execute

.execute:
	@$(eval API_HOSTNAME ?= $(shell ${KUBECTL} config view --minify | grep 'server:' | sed -E 's/ *server: +https:\/\/(.*):.+/\1/'))
	@$(eval API_HOSTIP ?= $(shell getent hosts ${API_HOSTNAME} | head -n 1 | awk '{ print $$1 }'))
	@echo K8s API Hostname: ${API_HOSTNAME}
	@echo K8s API Host IP: ${API_HOSTIP}
	${DORP} run --rm -it \
	  --network="host" \
	  --add-host="${API_HOSTNAME}:${API_HOSTIP}" \
	  -v "${KUBECONFIG}":/root/.kube/config:ro \
	  quay.io/kiali/scale-mesh-demo:latest \
	  --extra-vars "control_plane_namespace=${CONTROL_PLANE_NAMESPACE}" \
	  --extra-vars "mesh_type=${MESH_TYPE}" \
	  --extra-vars "number_namespaces=${NUMBER_NAMESPACES}" \
	  --extra-vars "number_apps=${NUMBER_APPS}" \
	  --extra-vars "number_versions=${NUMBER_VERSIONS}" \
	  --extra-vars "number_services=${NUMBER_SERVICES}" \
	  --extra-vars "state=${STATE}"

.set-state-install:
	@$(eval STATE = present)

.set-state-uninstall:
	@$(eval STATE = absent)
