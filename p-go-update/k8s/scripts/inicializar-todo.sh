#!/bin/bash

set -e  # Detener si algo falla

BACKUP_DIR=./backups/mongodump

# ------------------ BACKUP ANTES DE ELIMINAR ------------------
echo "💾 Realizando respaldo de datos Mongo (si existe clúster)..."

POD_NAME=$(kubectl get pods -n mongo-ns -l app=mongo -o jsonpath="{.items[0].metadata.name}" 2>/dev/null || echo "")

if [ -z "$POD_NAME" ]; then
  echo "ℹ️ No se encontró ningún pod Mongo. Se omite respaldo."
else
  echo "📦 Ejecutando mongodump con autenticación segura..."

  MONGO_USER=$(kubectl get secret mongo-db-secret -n mongo-ns -o jsonpath="{.data.MONGO_ROOT_USER}" | base64 --decode)
  MONGO_PASS=$(kubectl get secret mongo-db-secret -n mongo-ns -o jsonpath="{.data.MONGO_ROOT_PASS}" | base64 --decode)

  mkdir -p ./backups

  kubectl exec -n mongo-ns "$POD_NAME" -- \
    mongodump --out /tmp/mongodump \
    --username="$MONGO_USER" \
    --password="$MONGO_PASS" \
    --authenticationDatabase=admin || echo "⚠️ mongodump falló, puede que Mongo no esté listo."

  kubectl cp mongo-ns/"$POD_NAME":/tmp/mongodump "$BACKUP_DIR" || echo "⚠️ No se pudo copiar el respaldo."
fi

# ------------------ ELIMINAR CLUSTER ------------------
echo "💨 Eliminando clúster anterior (si existe)..."
kind delete cluster --name sistema-personas

# ------------------ CREAR CLUSTER ------------------
echo "🚀 Creando clúster nuevo desde kind-config.yaml..."
kind create cluster --config k8s/kind-config.yaml

# ------------------ PERMISOS ------------------
echo "🔧 Asegurando permisos de ejecución para scripts..."
chmod +x k8s/scripts/inicializar-create.sh
chmod +x ../p-go-search/k8s/inicializar-search.sh
chmod +x ../p-go-list/k8s/inicializar-list.sh
chmod +x ../p-go-update/k8s/inicializar-update.sh
chmod +x ../p-go-delete/k8s/inicializar-delete.sh

# ------------------ INGRESS ------------------
echo "🌐 Instalando Ingress NGINX..."
kubectl apply -f k8s/ingress-nginx.yaml

kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=180s || echo "⚠️ Ingress puede no estar listo."

# ------------------ MONGO ------------------
echo "🛠️ Desplegando base de datos MongoDB..."
kubectl apply -f k8s/mongo/namespace-mongo.yaml
kubectl apply -f k8s/mongo/pv-mongo.yaml
kubectl apply -f k8s/mongo/pvc-mongo.yaml
kubectl apply -f k8s/mongo/secrets-mongo.yaml
kubectl apply -f k8s/mongo/dockerhub-secret.yaml
kubectl apply -f k8s/mongo/deployment-mongo.yaml
kubectl apply -f k8s/mongo/service-mongo.yaml

kubectl wait --namespace=mongo-ns \
  --for=condition=available deployment/mongo-deployment \
  --timeout=180s || echo "⚠️ Mongo puede no estar disponible aún."

# ------------------ RESTAURAR BACKUP ------------------
if [ -d "$BACKUP_DIR" ]; then
  echo "♻️ Restaurando respaldo de Mongo..."

  POD_NAME=$(kubectl get pods -n mongo-ns -l app=mongo -o jsonpath="{.items[0].metadata.name}")
  MONGO_USER=$(kubectl get secret mongo-db-secret -n mongo-ns -o jsonpath="{.data.MONGO_ROOT_USER}" | base64 --decode)
  MONGO_PASS=$(kubectl get secret mongo-db-secret -n mongo-ns -o jsonpath="{.data.MONGO_ROOT_PASS}" | base64 --decode)

  kubectl cp "$BACKUP_DIR" mongo-ns/"$POD_NAME":/tmp/mongodump

  echo "⏳ Esperando que Mongo acepte conexiones..."
  for i in {1..10}; do
    if kubectl exec -n mongo-ns "$POD_NAME" -- \
      mongosh --username="$MONGO_USER" --password="$MONGO_PASS" --authenticationDatabase=admin --eval "db.runCommand({ ping: 1 })" > /dev/null 2>&1; then
      echo "✅ Mongo respondió."
      break
    else
      echo "⏰ Esperando 5 segundos..."
      sleep 5
    fi
  done

  kubectl exec -n mongo-ns "$POD_NAME" -- \
    mongorestore /tmp/mongodump \
    --username="$MONGO_USER" \
    --password="$MONGO_PASS" \
    --authenticationDatabase=admin

  echo "✅ Respaldo restaurado correctamente."
else
  echo "ℹ️ No se encontró respaldo previo. Se omite restauración."
fi

# ------------------ MICROSERVICIOS ------------------
echo "🧩 Desplegando p-go-create..."
./k8s/scripts/inicializar-create.sh

echo "🔍 Desplegando p-go-search..."
(cd ../p-go-search && ./k8s/inicializar-search.sh)

echo "📋 Desplegando p-go-list..."
(cd ../p-go-list && ./k8s/inicializar-list.sh)

echo "✏️ Desplegando p-go-update..."
(cd ../p-go-update && ./k8s/inicializar-update.sh)

echo "❌ Desplegando p-go-delete..."
(cd ../p-go-delete && ./k8s/inicializar-delete.sh)

# ------------------ FINAL ------------------
echo -e "\n🎉 Todo el clúster ha sido desplegado correctamente 🎉"
kubectl get all -A
