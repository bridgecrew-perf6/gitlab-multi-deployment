#!/bin/bash
cd $4
kubectl apply -f application/kubernetes/namespace.yaml
kubectl config set-context --current --namespace=$1
kubectl apply -f application/kubernetes/serviceaccount.yaml
kubectl apply -f application/kubernetes/service.yaml
kubectl apply -f application/kubernetes/prod/deployment.yaml
kubectl apply -f application/kubernetes/prod/ingress.yaml
if [ -e application/kubernetes/prod/cronJob.yaml ]; then kubectl apply -f application/kubernetes/prod/cronJob.yaml; fi
if [ -e application/kubernetes/autoscaler.yaml ]; then kubectl apply -f application/kubernetes/autoscaler.yaml; fi
if [ -e application/kubernetes/rolebinding.yaml ]; then kubectl apply -f application/kubernetes/rolebinding.yaml; fi
