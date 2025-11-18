#!/bin/bash

set -e

# 🔍 Validar que el clúster esté corriendo
if ! kubectl cluster-info > /dev/null 2>&1; then
  echo "❌ El clúster de Kubernetes no está activo. Por favor, ejecútalo con kind antes de continuar."
  exit 1
fi

# -------- CONFIGURACIÓN --------
REPO_BASE="danysoftdev/p-go"
TAG="v$(date +%Y%m%d%H%M%S)"
DOCUMENTO="9999"

echo "🔄 Iniciando construcción, despliegue y verificación..."

# -------- CREATE --------
echo -e "\n🧩 p-go-create"
docker build -t "$REPO_BASE-create:$TAG" ../p-go-create
docker push "$REPO_BASE-create:$TAG"
kubectl set image deployment/create-deployment create-container="$REPO_BASE-create:$TAG" -n p-go-create
kubectl rollout status deployment/create-deployment -n p-go-create

# ➕ Crear persona
echo "📤 POST /crear-personas"
curl -s -o /dev/null -w "🔗 HTTP %{http_code}\n" -X POST http://localhost/crear/crear-personas \
  -H "Content-Type: application/json" \
  -d '{"documento":"'"$DOCUMENTO"'", "nombre":"Ana", "apellido":"Díaz", "edad":28, "correo":"ana@test.com", "telefono":"1234567890", "direccion":"Calle 123"}'

# -------- LIST --------
echo -e "\n📋 p-go-list"
docker build -t "$REPO_BASE-list:$TAG" ../p-go-list
docker push "$REPO_BASE-list:$TAG"
kubectl set image deployment/list-deployment list-container="$REPO_BASE-list:$TAG" -n p-go-list
kubectl rollout status deployment/list-deployment -n p-go-list

# 📄 Listar personas
echo "📄 GET /listar-personas"
curl -s -o /dev/null -w "🔗 HTTP %{http_code}\n" http://localhost/listar/listar-personas

# -------- SEARCH --------
echo -e "\n🔍 p-go-search"
docker build -t "$REPO_BASE-search:$TAG" ../p-go-search
docker push "$REPO_BASE-search:$TAG"
kubectl set image deployment/search-deployment search-container="$REPO_BASE-search:$TAG" -n p-go-search
kubectl rollout status deployment/search-deployment -n p-go-search

# 🔎 Buscar persona por documento
echo "🔎 GET /buscar-personas/$DOCUMENTO"
curl -s -o /dev/null -w "🔗 HTTP %{http_code}\n" http://localhost/buscar/buscar-personas/$DOCUMENTO

# -------- UPDATE --------
echo -e "\n✏️ p-go-update"
docker build -t "$REPO_BASE-update:$TAG" ../p-go-update
docker push "$REPO_BASE-update:$TAG"
kubectl set image deployment/update-deployment update-container="$REPO_BASE-update:$TAG" -n p-go-update
kubectl rollout status deployment/update-deployment -n p-go-update

# ✏️ Actualizar persona
echo "✏️ PUT /actualizar-personas/$DOCUMENTO"
curl -s -o /dev/null -w "🔗 HTTP %{http_code}\n" -X PUT http://localhost/actualizar/actualizar-personas/$DOCUMENTO \
  -H "Content-Type: application/json" \
  -d '{"documento":"'"$DOCUMENTO"'", "nombre":"Ana Actualizada", "apellido":"Díaz", "edad":29, "correo":"ana_actualizada@test.com", "telefono":"0987654321", "direccion":"Nueva dirección"}'

# -------- DELETE --------
echo -e "\n❌ p-go-delete"
docker build -t "$REPO_BASE-delete:$TAG" ../p-go-delete
docker push "$REPO_BASE-delete:$TAG"
kubectl set image deployment/delete-deployment delete-container="$REPO_BASE-delete:$TAG" -n p-go-delete
kubectl rollout status deployment/delete-deployment -n p-go-delete

# 🗑️ Eliminar persona
echo "🗑️ DELETE /eliminar-personas/$DOCUMENTO"
curl -s -o /dev/null -w "🔗 HTTP %{http_code}\n" -X DELETE http://localhost/eliminar/eliminar-personas/$DOCUMENTO

echo -e "\n✅ Todo verificado correctamente. Automatización finalizada 🎯"
