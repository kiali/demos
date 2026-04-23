#!/usr/bin/env bash
set -e

CLUSTER_NAME=${CLUSTER_NAME:-ci}
NEW_TAG="v1-$(date +%Y%m%d-%H%M%S)"

echo "==================================================================="
echo "Reloading images in kind cluster '${CLUSTER_NAME}' with tag: ${NEW_TAG}"
echo "==================================================================="

# 1. Retag images
echo ""
echo "1. Retagging local images..."
for service in cars discounts flights hotels insurances travels mysqldb control portal loadtester; do
  echo "  - ${service}"
  podman tag quay.io/kiali/demo_travels_${service}:v1 quay.io/kiali/demo_travels_${service}:${NEW_TAG} 2>/dev/null || echo "    (image not found, skipping)"
done

# 2. Load images into kind
echo ""
echo "2. Loading images into kind cluster '${CLUSTER_NAME}'..."
for service in cars discounts flights hotels insurances travels mysqldb control portal loadtester; do
  if podman image exists quay.io/kiali/demo_travels_${service}:${NEW_TAG}; then
    echo "  - ${service}"
    # Remove any existing tar file first
    rm -f /tmp/${service}.tar
    # Remove image from kind if it exists
    docker exec ${CLUSTER_NAME}-control-plane crictl rmi quay.io/kiali/demo_travels_${service}:${NEW_TAG} 2>/dev/null || true
    # Load new image
    podman save quay.io/kiali/demo_travels_${service}:${NEW_TAG} -o /tmp/${service}.tar
    kind load image-archive /tmp/${service}.tar --name ${CLUSTER_NAME}
    rm -f /tmp/${service}.tar
  fi
done

# 3. Update deployments
echo ""
echo "3. Updating deployments..."

echo "  - travel-agency namespace"
for service in cars discounts flights hotels insurances travels mysqldb; do
  kubectl set image "deployment/${service}-v1" -n travel-agency "${service}=quay.io/kiali/demo_travels_${service}:${NEW_TAG}" 2>/dev/null || echo "    ${service}-v1 not found"
done

echo "  - travel-portal namespace"
for dep in travels viaggi voyages; do
  kubectl set image "deployment/${dep}" -n travel-portal "control=quay.io/kiali/demo_travels_portal:${NEW_TAG}" 2>/dev/null || echo "    ${dep} not found"
done

echo "  - travel-control namespace"
kubectl set image "deployment/control" -n travel-control "control=quay.io/kiali/demo_travels_control:${NEW_TAG}" 2>/dev/null || echo "    control not found"

# 4. Wait for pods to be ready
echo ""
echo "4. Waiting for pods to be ready..."
kubectl wait --for=condition=ready pod --all -n travel-agency --timeout=120s 2>&1 | grep -v waypoint || true
kubectl wait --for=condition=ready pod --all -n travel-portal --timeout=120s 2>&1 | grep -v waypoint || true
kubectl wait --for=condition=ready pod --all -n travel-control --timeout=120s 2>&1 | grep -v waypoint || true

echo ""
echo "==================================================================="
echo "✓ Images reloaded with tag: ${NEW_TAG}"
echo "==================================================================="
