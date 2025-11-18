#!/bin/bash

set -e  # Detiene el script si ocurre un error

echo "🧨 Eliminando clúster anterior (si existe)..."
kind delete cluster --name sistema-personas

echo "🚀 Creando nuevo clúster desde kind-config.yaml..."
kind create cluster --config k8s/kind-config.yaml

echo "🌐 Instalando NGINX Ingress Controller..."
kubectl apply -f k8s/ingress-nginx.yaml

echo "⏳ Esperando a que el Ingress Controller esté 'Running'..."
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controllers \
  --timeout=180s

echo "✅ Ingress listo."

echo "📂 Aplicando recursos de MongoDB..."

kubectl apply -f k8s/mongo/namespace-mongo.yaml
kubectl apply -f k8s/mongo/pv-mongo.yaml
kubectl apply -f k8s/mongo/pvc-mongo.yaml
kubectl apply -f k8s/mongo/secrets-mongo.yaml
kubectl apply -f k8s/mongo/deployment-mongo.yaml
kubectl apply -f k8s/mongo/service-mongo.yaml

echo "⏳ Esperando a que Mongo esté listo..."
kubectl wait --namespace=mongo-ns \
  --for=condition=available deployment/mongo-deployment \
  --timeout=180s

echo "✅ MongoDB desplegado y funcionando."

echo "📂 Aplicando recursos del microservicio p-go-create..."

kubectl apply -f k8s/create/namespace-create.yaml
kubectl apply -f k8s/create/secrets-create.yaml
kubectl apply -f k8s/create/deployment-create.yaml
kubectl apply -f k8s/create/service-create.yaml
kubectl apply -f k8s/create/ingress.yaml

echo "⏳ Esperando a que p-go-create esté listo..."
kubectl wait --namespace=p-go-create \
  --for=condition=available deployment/create-deployment \
  --timeout=90s

echo "✅ Microservicio p-go-create desplegado y listo."

echo -e "\n🔍 Estado final por namespace:\n"

echo "🔵 Namespace mongo-ns:"
kubectl get all -n mongo-ns

echo -e "\n🟢 Namespace p-go-create:"
kubectl get all -n p-go-create
kubectl get ingress -n p-go-create

echo -e "\n🎉 ¡Clúster desplegado y operativo correctamente!"
