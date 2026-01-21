#!/bin/bash
# Quick Test Deployment: 1 Lobby + 1 Game World
# Bash version for Linux/Mac

set -e

NAMESPACE="${NAMESPACE:-game-servers}"
CHART_PATH="${CHART_PATH:-./charts/game-servers}"

echo "═══════════════════════════════════════════════════════"
echo "  Hytale Test Network Deployment"
echo "  1 Lobby + 1 Game World"
echo "═══════════════════════════════════════════════════════"
echo ""

# Check if namespace exists
if ! kubectl get namespace "$NAMESPACE" &>/dev/null; then
    echo "→ Creating namespace: $NAMESPACE"
    kubectl create namespace "$NAMESPACE"
else
    echo "✓ Namespace exists: $NAMESPACE"
fi

echo ""
echo "────────────────────────────────────────────────────────"
echo "  Step 1: Deploying Lobby Server"
echo "────────────────────────────────────────────────────────"
echo ""

helm install hytale-lobby "$CHART_PATH" \
  -f ./charts/values/ht-lobby-values.yaml \
  --namespace "$NAMESPACE" \
  --timeout 10m

echo ""
echo "✓ Lobby server deployment started"

echo ""
echo "────────────────────────────────────────────────────────"
echo "  Step 2: Deploying Survival World Server"
echo "────────────────────────────────────────────────────────"
echo ""

helm install hytale-survival "$CHART_PATH" \
  -f ./charts/values/ht-world-values.yaml \
  --set nameOverride=hytale-survival \
  --set fullnameOverride=hytale-survival \
  --set serverConfig.SERVER_PORT=5521 \
  --set serverConfig.HYTALE_SERVER_NAME="🌲 Survival World" \
  --set serverConfig.HYTALE_WORLD=survival \
  --set podLabels.world-name=survival \
  --set gameService.ports[0].port=5521 \
  --set gameService.ports[0].targetPort=5521 \
  --set gameService.ports[1].port=5521 \
  --set gameService.ports[1].targetPort=5521 \
  --namespace "$NAMESPACE" \
  --timeout 10m

echo ""
echo "✓ Survival world deployment started"

echo ""
echo "════════════════════════════════════════════════════════"
echo "  Deployment Initiated!"
echo "════════════════════════════════════════════════════════"
echo ""

echo "Waiting for pods to start (this may take a few minutes)..."
sleep 10

echo ""
echo "Current Pod Status:"
kubectl get pods -n "$NAMESPACE" -l app.kubernetes.io/name=game-servers

echo ""
echo "Current Services:"
kubectl get svc -n "$NAMESPACE" -l app.kubernetes.io/name=game-servers

echo ""
echo "════════════════════════════════════════════════════════"
echo "  Next Steps:"
echo "════════════════════════════════════════════════════════"
echo ""
echo "1. Watch deployment progress:"
echo "   kubectl get pods -n $NAMESPACE -w"
echo ""
echo "2. Once pods are running, authenticate lobby server:"
echo "   kubectl attach -it \$(kubectl get pod -n $NAMESPACE -l nameOverride=hytale-lobby -o name) -n $NAMESPACE"
echo ""
echo "3. Authenticate survival server:"
echo "   kubectl attach -it \$(kubectl get pod -n $NAMESPACE -l nameOverride=hytale-survival -o name) -n $NAMESPACE"
echo ""
echo "4. Check logs:"
echo "   kubectl logs -f -n $NAMESPACE -l nameOverride=hytale-lobby"
echo "   kubectl logs -f -n $NAMESPACE -l nameOverride=hytale-survival"
echo ""
echo "5. Get lobby external IP:"
echo "   kubectl get svc hytale-lobby-game -n $NAMESPACE"
echo ""
echo "Server DNS Names for Lobby Plugin:"
echo "   Survival: hytale-survival-game.$NAMESPACE.svc.cluster.local:5521"
echo ""
