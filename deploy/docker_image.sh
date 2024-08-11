#!/bin/bash

tag=$1
IMAGE_NAME="user-service"
NAMESPACE="go-microservices"
REGISTRY="registry.cn-shanghai.aliyuncs.com"

if [[ -z "$tag" ]]; then
  echo "tag is empty"
  echo "usage: sh $0 <tag>"
  exit 1
fi

# build image
echo "1. build docker image"
docker build -t ${REGISTRY}/${NAMESPACE}/${IMAGE_NAME}:${tag} -f deploy/docker/Dockerfile .

# docker tag local-image:tagname new-repo:tagname
echo "2. gen docker tag"
docker tag ${IMAGE_NAME}:${tag} ${REGISTRY}/${NAMESPACE}/${IMAGE_NAME}:${tag}

# docker push new-repo:tagname
echo "3. push docker image to hub"
docker push ${REGISTRY}/${NAMESPACE}/${IMAGE_NAME}:${tag}

echo "Done. push docker image success."