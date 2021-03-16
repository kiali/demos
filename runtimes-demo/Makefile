# OCI CLI (docker or podman)
OCI_BIN ?= $(shell which podman 2>/dev/null || which docker 2>/dev/null)
OCI_BIN_SHORT = $(shell if [[ ${OCI_BIN} == *"podman" ]]; then echo "podman"; else echo "docker"; fi)
# Tag for docker images
TAG ?= dev
OCI_USER ?= jotak
# Set REMOTE=true if you want to use remote (quay.io) images
REMOTE ?= false
NAMESPACE ?= default
# List of images
SERVICES ?= quarkus vertx nodejs springboot
ISTIO ?= true

LATEST = 0.1.0

ifeq ($(REMOTE),true)
OCI_DOMAIN ?= quay.io/
OCI_DOMAIN_IN_CLUSTER ?= quay.io/
PULL_POLICY ?= "IfNotPresent"
else ifeq ($(OCI_BIN_SHORT),podman)
OCI_DOMAIN ?= "$(shell minikube ip):5000/"
OCI_DOMAIN_IN_CLUSTER ?= localhost:5000/
PULL_POLICY ?= "Always"
TAG_MINIKUBE = true
else
OCI_DOMAIN ?= ""
OCI_DOMAIN_IN_CLUSTER ?= ""
PULL_POLICY ?= "Never"
endif

ifeq ($(OCI_BIN_SHORT),podman)
PUSH_OPTS ?= --tls-verify=false
endif

ifeq ($(ISTIO),true)
ISTIO_LABEL = "istio-injection=enabled"
else
ISTIO_LABEL = "istio-injection-"
endif

bold := $(shell tput bold)
smul := $(shell tput smul)
sgr0 := $(shell tput sgr0)

man: help

help:
	@echo "$(bold)Runtimes Demo$(sgr0)"
	@echo "🏃 🕔 📊"
	@echo ""
	@echo "Note that you may need a couple of tools depending of which target you run: kubectl, npm, mvn, docker/podman"
	@echo ""
	@echo "$(smul)Main options$(sgr0):"
	@echo "- $(bold)TAG$(sgr0):          Image tag to be used (current: $(TAG))"
	@echo "- $(bold)OCI_USER$(sgr0):     User to use for image tagging / pulling (current: $(OCI_USER))"
	@echo "- $(bold)REMOTE$(sgr0):       When true, options are pre-defined to use quay.io; when false it would default to Minikube (current: $(REMOTE))"
	@echo "- $(bold)NAMESPACE$(sgr0):    Namespace where to deploy everything (current: $(NAMESPACE))"
	@echo ""
	@echo "$(smul)Main targets$(sgr0):"
	@echo "- $(bold)build$(sgr0):               Builds sources (maven/npm)"
	@echo "- $(bold)images$(sgr0):              Builds OCI images, eventually prepare them for Minikube (pushing to local registry)"
	@echo "- $(bold)push-minikube$(sgr0):       Push built images to Minikube and tag for it"
	@echo "- $(bold)push$(sgr0):                Push to registry (local or remote)"
	@echo "- $(bold)deploy$(sgr0):              Deploy to Kubernetes"
	@echo "- $(bold)kill$(sgr0):                Kill (restart) all pods"
	@echo "- $(bold)reload-minikube$(sgr0):     Push images to minikube and restart pods"
	@echo "- $(bold)undeploy$(sgr0):            Remove everything deployed"
	@echo "- $(bold)deploy-latest$(sgr0):       Deploys by pulling last known images from quay.io"
	@echo ""
	@echo "$(smul)Special targets$(sgr0):"
	@echo "- $(bold)gen-quickstart$(sgr0):           Generate the quickstart templates with latest tag"
	@echo "- $(bold)istio-enable$(sgr0):             Enable istio on namespace '${NAMESPACE}'"
	@echo ""
	@echo "$(smul)Examples$(sgr0):"
	@echo ""
	@echo "$(bold)make deploy-latest$(sgr0)"
	@echo "  Quick start, using latest images. Runtime metrics are enabled."
	@echo ""
	@echo "$(bold)NAMESPACE=runtimes TAG=$(LATEST) REMOTE=true make deploy$(sgr0)"
	@echo "  Deploy images from quay.io, tagged $(LATEST), to namespace 'runtimes'"
	@echo ""
	@echo "$(bold)make build images deploy$(sgr0)"
	@echo "  Build everything, deploy to Minikube."
	@echo ""
	@echo "$(bold)OCI_USER=$(OCI_USER) TAG=$(LATEST) REMOTE=true make build images push gen-quickstart$(sgr0)"
	@echo "  Build everything, push to quay.io and generate quickstart templates"
	@echo ""
	@echo "Enjoy! 🚀"

build:
	@echo "🏃🕔 Building services..."
	for svc in ${SERVICES} ; do \
		echo "Building $$svc:" ; \
		cd $$svc && ./build.sh ; \
		cd .. ; \
	done

ifeq ($(TAG_MINIKUBE),true)
images: images-build push-minikube
else
images: images-build
endif

images-build:
	@echo "🏃🕔 Building images..."
	for svc in ${SERVICES} ; do \
		echo "Building $$svc:" ; \
		${OCI_BIN_SHORT} build -t ${OCI_DOMAIN}${OCI_USER}/runtimes-$$svc:${TAG} -f ./kube/$$svc.dockerfile . ; \
	done

push-minikube:
	@echo "🏃🕔 Tagging for Minikube..."
	for svc in ${SERVICES} ; do \
		${OCI_BIN_SHORT} push ${PUSH_OPTS} ${OCI_DOMAIN}${OCI_USER}/runtimes-$$svc:${TAG} ; \
		${OCI_BIN_SHORT} tag ${OCI_DOMAIN}${OCI_USER}/runtimes-$$svc:${TAG} ${OCI_DOMAIN_IN_CLUSTER}${OCI_USER}/runtimes-$$svc:${TAG} ; \
	done

push:
	@echo "🏃🕔 Pushing images..."
	for svc in ${SERVICES} ; do \
		${OCI_BIN_SHORT} push ${PUSH_OPTS} ${OCI_DOMAIN}${OCI_USER}/runtimes-$$svc:${TAG} ; \
	done

deploy:
	kubectl label namespace ${NAMESPACE} ${ISTIO_LABEL} 2> /dev/null ; \
	kubectl apply -n ${NAMESPACE} -f ./kube/loadgen.yml ; \
	for svc in ${SERVICES} ; do \
		./gentpl.sh $$svc -pp ${PULL_POLICY} -d "${OCI_DOMAIN_IN_CLUSTER}" -u ${OCI_USER} -t ${TAG} ${GENTPL_OPTS} | kubectl -n ${NAMESPACE} apply -f - ; \
	done

kill:
	kubectl -n ${NAMESPACE} delete pods -l "project=runtimes"

undeploy:
	kubectl -n ${NAMESPACE} delete all -l "project=runtimes" ; \
	kubectl -n ${NAMESPACE} delete destinationrule -l "project=runtimes" ; \
	kubectl -n ${NAMESPACE} delete virtualservice -l "project=runtimes" ; \
	kubectl -n ${NAMESPACE} delete gateway -l "project=runtimes" ; \
	kubectl -n ${NAMESPACE} delete envoyfilter -l "project=runtimes"

deploy-latest: TAG=${LATEST}
deploy-latest: OCI_DOMAIN_IN_CLUSTER=quay.io/
deploy-latest: PULL_POLICY=IfNotPresent
deploy-latest: TAG_MINIKUBE=false
deploy-latest: deploy

istio-enable:
	kubectl label namespace ${NAMESPACE} ${ISTIO_LABEL} ; \
	kubectl -n ${NAMESPACE} delete pods -l "project=runtimes"

gen-quickstart:
	@echo "🏃🕔 Generating quickstart template..."
	rm quickstart.yml ; \
	for svc in ${SERVICES} ; do \
		./gentpl.sh $$svc -pp IfNotPresent -d "quay.io/" -u jotak -t ${LATEST} >> quickstart.yml ; \
	done ; \
	cat kube/loadgen.yml >> quickstart.yml ; \

reload-minikube: push-minikube kill
