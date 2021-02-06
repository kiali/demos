#/bin/bash

##########################################################################
# scale-mesh.sh
#
# Installs a mesh with customized size and topology.
##########################################################################

set -u

DEFAULT_MESH_TYPE="circle"
DEFAULT_NUMBER_NAMESPACES="1"
DEFAULT_NUMBER_APPS="5"
DEFAULT_NUMBER_VERSIONS="2"
DEFAULT_NUMBER_SERVICES="5"

DEFAULT_TG_DURATION="0s"
DEFAULT_TG_RATE="1"

# Only used when creating SMM resources within a Maistra environment
DEFAULT_CONTROL_PLANE_NAMESPACE="istio-system"

DEFAULT_KUBECONFIG="${HOME}/.kube/config"
DEFAULT_MINIKUBECONFIG="${HOME}/.minikube"

DEFAULT_DORP="docker"
DEFAULT_KUBECTL="kubectl"

discover_api_host_ip() {
  if ! which ${KUBECTL} > /dev/null 2>&1; then
    echo "Cannot auto-discover the k8s API host/IP - invalid kubectl command [${KUBECTL}]. See the --kubectl option."
    exit 111
  fi
  API_HOSTNAME="$(${KUBECTL} config view --minify | grep 'server:' | sed -E 's/ *server: +https:\/\/(.*):.+/\1/')"
  if [[ "${API_HOSTNAME}" =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    API_HOSTIP="${API_HOSTNAME}" # the hostname is the IP
    echo "k8s API Host: ${API_HOSTIP}"
  else
    API_HOSTIP="$(getent hosts ${API_HOSTNAME} | head -n 1 | awk '{ print $1 }')"
    echo "k8s API Host: ${API_HOSTNAME} -> ${API_HOSTIP}"
  fi
}

execute() {
  if [ -z "${API_HOSTNAME}" -a -z "${API_HOSTIP}" ]; then
    discover_api_host_ip
  fi

  # if the host and IP are different, it means we need to add that host:ip to the docker process
  if [ "${API_HOSTNAME}" != "${API_HOSTIP}" ]; then
    if [ -z "${API_HOSTNAME}" ]; then
      echo "Missing the API hostname - see the --api-hostname option"
      exit 1
    fi
    if [ -z "${API_HOSTIP}" ]; then
      echo "Missing the API IP - see the --api-ip option"
      exit 1
    fi
    ADD_HOST_ARG="--add-host=${API_HOSTNAME}:${API_HOSTIP}"
  fi

  if [ -d "${MINIKUBECONFIG}" ]; then
    MINIKUBE_CONFIG_VOL_ARG="-v ${MINIKUBECONFIG}:${MINIKUBECONFIG}"
  fi

  ${DORP} run --rm -it \
    --network="host" \
    -v "${KUBECONFIG}":/root/.kube/config:ro \
    ${ADD_HOST_ARG:-} \
    ${MINIKUBE_CONFIG_VOL_ARG:-} \
    quay.io/kiali/scale-mesh-demo:latest \
    --extra-vars "mesh_type=${MESH_TYPE}" \
    --extra-vars "number_namespaces=${NUMBER_NAMESPACES}" \
    --extra-vars "number_apps=${NUMBER_APPS}" \
    --extra-vars "number_versions=${NUMBER_VERSIONS}" \
    --extra-vars "number_services=${NUMBER_SERVICES}" \
    --extra-vars "tg_duration=${TG_DURATION}" \
    --extra-vars "tg_rate=${TG_RATE}" \
    --extra-vars "control_plane_namespace=${CONTROL_PLANE_NAMESPACE}" \
    --extra-vars "state=${STATE}"
}

# Change to the directory where this script is and set our env
cd "$(dirname "${BASH_SOURCE[0]}")"

_CMD=""
while [[ $# -gt 0 ]]; do
  key="$1"
  case $key in
    install)                        _CMD="install"; shift ;;
    uninstall)                      _CMD="uninstall"; shift ;;
    -a|--apps)                      NUMBER_APPS="${2}"; shift; shift ;;
    -ah|--api-hostname)             API_HOSTNAME="${2}"; shift; shift ;;
    -ai|--api-ip)                   API_HOSTIP="${2}"; shift; shift ;;
    -dorp|--docker-or-podman)       DORP="${2}"; shift; shift ;;
    -k|--kubectl)                   KUBECTL="${2}"; shift; shift ;;
    -kc|--kube-config)              KUBECONFIG="${2}"; shift; shift ;;
    -mc|--minikube-config)          MINIKUBECONFIG="${2}"; shift; shift ;;
    -mt|--mesh-type)                MESH_TYPE="${2}"; shift; shift ;;
    -cpn|--control-plane-namespace) CONTROL_PLANE_NAMESPACE="${2}"; shift; shift ;;
    -n|--namespaces)                NUMBER_NAMESPACES="${2}"; shift; shift ;;
    -s|--services)                  NUMBER_SERVICES="${2}"; shift; shift ;;
    -tgd|--traffic-gen-duration)    TG_DURATION="${2}"; shift; shift ;;
    -tgr|--traffic-gen-rate)        TG_RATE="${2}"; shift; shift ;;
    -v|--versions)                  NUMBER_VERSIONS="${2}"; shift; shift ;;
    -h|--help)
      cat <<HELPMSG

$0 [option...] command

Valid options:
  -a|--apps
      The number of apps in the mesh
      Default: ${DEFAULT_NUMBER_APPS}
  -ah|--api-hostname
      The k8s API host (if bound to a hostname).
      If not specified, an attempt to auto-discover it will be made.
  -ai|--api-ip
      The k8s API IP - if you pass in --api-hostname option, you must provide it's IP using this option.
      If not specified, an attempt to auto-discover it will be made.
  -dorp|--docker-or-podman
      Indicates if you want to use 'docker' or 'podman'.
      Default: ${DEFAULT_DORP}
  -k|--kubectl
      The kubectl executable - if not an absolute path then it must be in PATH.
      This is only used when attempting to auto-discover the k8s API host/IP.
      Default: ${DEFAULT_KUBECTL}
  -kc|--kube-config
      Where the kubectl configuration directory is located
      Default: ${DEFAULT_KUBECONFIG}
  -mc|--minikube-config
      Where the Minikube configuration directory is located.
      This is only needed if you are accessing k8s within a Minikube environment.
      Default: ${DEFAULT_MINIKUBECONFIG}
  -mt|--mesh-type
      Determines the mesh topology to be used
      Must be one of: breadth, breadth-sink, circle, circle-callback, depth, depth-sink, hourglass
      Default: ${DEFAULT_MESH_TYPE}
  -cpn|--control-plane-namespace
      Where Istio's control plane components are installed.
      This is only used when creating SMM resources within a Maistra environment.
      Default: ${DEFAULT_CONTROL_PLANE_NAMESPACE}
  -n|--namespaces
      The number of namespaces in the mesh
      Default: ${DEFAULT_NUMBER_NAMESPACES}
  -s|--services
      The number of services in the mesh.
      Default: ${DEFAULT_NUMBER_SERVICES}
  -tgd|--traffic-gen-duration
      The duration of requests generated by the traffic generator
      Default: ${DEFAULT_TG_DURATION}
  -tgr|--traffic-gen-rate
      The rate of requests generated by the traffic generator
      Default: ${DEFAULT_TG_RATE}
  -v|--versions
      The number of versions per app
      Default: ${DEFAULT_NUMBER_VERSIONS}

The command must be either:
  install:   installs the scale-mesh components
  uninstall: uninstalls the scale-mesh components

HELPMSG
      exit 1
      ;;
    *)
      echo "Unknown argument [$key]. Aborting."
      exit 1
      ;;
  esac
done

# Prepare env vars
: ${API_HOSTNAME:=}
: ${API_HOSTIP:=}
: ${MESH_TYPE:=${DEFAULT_MESH_TYPE}}
: ${NUMBER_NAMESPACES:=${DEFAULT_NUMBER_NAMESPACES}}
: ${NUMBER_APPS:=${DEFAULT_NUMBER_APPS}}
: ${NUMBER_VERSIONS:=${DEFAULT_NUMBER_VERSIONS}}
: ${NUMBER_SERVICES:=${DEFAULT_NUMBER_SERVICES}}
: ${CONTROL_PLANE_NAMESPACE:=${DEFAULT_CONTROL_PLANE_NAMESPACE}}
: ${DORP:=${DEFAULT_DORP}}
: ${KUBECONFIG:=${DEFAULT_KUBECONFIG}}
: ${MINIKUBECONFIG:=${DEFAULT_MINIKUBECONFIG}}
: ${KUBECTL:=${DEFAULT_KUBECTL}}
: ${TG_DURATION:=${DEFAULT_TG_DURATION}}
: ${TG_RATE:=${DEFAULT_TG_RATE}}

if [ "$_CMD" == "install" ]; then
  STATE="present"
  execute
elif [ "$_CMD" == "uninstall" ]; then
  STATE="absent"
  execute
else
  echo "ERROR: Missing required command. See --help for available commands and options."
  exit 1
fi

