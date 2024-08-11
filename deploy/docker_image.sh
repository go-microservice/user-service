#!/bin/bash

tag=$1
IMAGE_NAME="user-service"
NAMESPACE="go-microservices"
REGISTRY="registry.cn-hangzhou.aliyuncs.com"

if [[ -z "$tag" ]]; then
  echo "tag is empty"
  echo "usage: sh $0 <tag>"
  exit 1
fi

# build image
echo "1. build docker image"
docker build -t ${REGISTRY}/${NAMESPACE}/${IMAGE_NAME}:${tag} -t ${REGISTRY}/${NAMESPACE}/${IMAGE_NAME}:latest -f deploy/docker/Dockerfile .

# docker push new-repo:tagname
echo "2. push docker image to remote hub"
docker push ${REGISTRY}/${NAMESPACE}/${IMAGE_NAME}:${tag}

echo "Done. push docker image success."
