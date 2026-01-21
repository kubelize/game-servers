# Quick Test Deployment: 1 Lobby + 1 Game World
# Run this script to deploy a basic Hytale network for testing

# Configuration
$NAMESPACE = "game-servers"
$CHART_PATH = "./charts/game-servers"

Write-Host "═══════════════════════════════════════════════════════" -ForegroundColor Cyan
Write-Host "  Hytale Test Network Deployment" -ForegroundColor Cyan
Write-Host "  1 Lobby + 1 Game World" -ForegroundColor Cyan
Write-Host "═══════════════════════════════════════════════════════" -ForegroundColor Cyan
Write-Host ""

# Create namespace if it doesn't exist
Write-Host "→ Checking namespace: $NAMESPACE" -ForegroundColor Yellow
$namespaceExists = kubectl get namespace $NAMESPACE 2>$null
if (-not $namespaceExists) {
    Write-Host "  Creating namespace..." -ForegroundColor Gray
    kubectl create namespace $NAMESPACE
    Write-Host "✓ Namespace created" -ForegroundColor Green
} else {
    Write-Host "✓ Namespace exists" -ForegroundColor Green
}

Write-Host ""
Write-Host "────────────────────────────────────────────────────────" -ForegroundColor Gray
Write-Host "  Step 1: Deploying Lobby Server" -ForegroundColor Cyan
Write-Host "────────────────────────────────────────────────────────" -ForegroundColor Gray
Write-Host ""

helm install hytale-lobby $CHART_PATH `
  -f ./charts/values/ht-lobby-values.yaml `
  --namespace $NAMESPACE `
  --timeout 10m

if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Lobby server deployment started" -ForegroundColor Green
} else {
    Write-Host "✗ Lobby deployment failed" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "────────────────────────────────────────────────────────" -ForegroundColor Gray
Write-Host "  Step 2: Deploying Survival World Server" -ForegroundColor Cyan
Write-Host "────────────────────────────────────────────────────────" -ForegroundColor Gray
Write-Host ""

helm install hytale-survival $CHART_PATH `
  -f ./charts/values/ht-world-values.yaml `
  --set nameOverride=hytale-survival `
  --set fullnameOverride=hytale-survival `
  --set serverConfig.SERVER_PORT=5521 `
  --set "serverConfig.HYTALE_SERVER_NAME=🌲 Survival World" `
  --set serverConfig.HYTALE_WORLD=survival `
  --set podLabels.world-name=survival `
  --set gameService.ports[0].port=5521 `
  --set gameService.ports[0].targetPort=5521 `
  --set gameService.ports[1].port=5521 `
  --set gameService.ports[1].targetPort=5521 `
  --namespace $NAMESPACE `
  --timeout 10m

if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Survival world deployment started" -ForegroundColor Green
} else {
    Write-Host "✗ Survival world deployment failed" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "════════════════════════════════════════════════════════" -ForegroundColor Cyan
Write-Host "  Deployment Initiated!" -ForegroundColor Cyan
Write-Host "════════════════════════════════════════════════════════" -ForegroundColor Cyan
Write-Host ""

Write-Host "Waiting for pods to start (this may take a few minutes)..." -ForegroundColor Yellow
Write-Host ""

# Wait a bit for pods to be created
Start-Sleep -Seconds 10

# Show pod status
Write-Host "Current Pod Status:" -ForegroundColor Cyan
kubectl get pods -n $NAMESPACE -l app.kubernetes.io/name=game-servers

Write-Host ""
Write-Host "Current Services:" -ForegroundColor Cyan
kubectl get svc -n $NAMESPACE -l app.kubernetes.io/name=game-servers

Write-Host ""
Write-Host "════════════════════════════════════════════════════════" -ForegroundColor Cyan
Write-Host "  Next Steps:" -ForegroundColor Cyan
Write-Host "════════════════════════════════════════════════════════" -ForegroundColor Cyan
Write-Host ""
Write-Host "1. Watch deployment progress:" -ForegroundColor White
Write-Host "   kubectl get pods -n $NAMESPACE -w" -ForegroundColor Gray
Write-Host ""
Write-Host "2. Once pods are running, authenticate lobby server:" -ForegroundColor White
Write-Host "   kubectl attach -it (kubectl get pod -n $NAMESPACE -l nameOverride=hytale-lobby -o name) -n $NAMESPACE" -ForegroundColor Gray
Write-Host ""
Write-Host "3. Authenticate survival server:" -ForegroundColor White
Write-Host "   kubectl attach -it (kubectl get pod -n $NAMESPACE -l nameOverride=hytale-survival -o name) -n $NAMESPACE" -ForegroundColor Gray
Write-Host ""
Write-Host "4. Check logs:" -ForegroundColor White
Write-Host "   kubectl logs -f -n $NAMESPACE -l nameOverride=hytale-lobby" -ForegroundColor Gray
Write-Host "   kubectl logs -f -n $NAMESPACE -l nameOverride=hytale-survival" -ForegroundColor Gray
Write-Host ""
Write-Host "5. Get lobby external IP:" -ForegroundColor White
Write-Host "   kubectl get svc hytale-lobby-game -n $NAMESPACE" -ForegroundColor Gray
Write-Host ""
Write-Host "Server DNS Names for Lobby Plugin:" -ForegroundColor White
Write-Host "   Survival: hytale-survival-game.$NAMESPACE.svc.cluster.local:5521" -ForegroundColor Yellow
Write-Host ""
