#!/bin/bash

set -e  # Detiene el script si hay un error

echo "📁 Desplegando microservicio p-go-create..."

# Namespace
kubectl apply -f k8s/create/namespace-create.yaml

# Secret
kubectl apply -f k8s/create/secrets-create.yaml
kubectl apply -f k8s/create/dockerhub-secret.yaml

# Deployment
kubectl apply -f k8s/create/deployment-create.yaml

# Esperar a que el deployment esté disponible
echo "⏳ Esperando a que p-go-create esté listo..."
kubectl wait --namespace=p-go-create \
  --for=condition=available deployment/create-deployment \
  --timeout=90s

# Service
kubectl apply -f k8s/create/service-create.yaml

# Ingress
kubectl apply -f k8s/create/ingress.yaml

echo "✅ p-go-create desplegado correctamente."

echo -e "\n🔍 Estado actual:"
kubectl get all -n p-go-create
kubectl get ingress -n p-go-create