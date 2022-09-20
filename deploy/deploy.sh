#!/bin/bash

ver=$1
serviceName="user-service"
repoAddr="qloog"
# aliyun: repoAddr="registry.cn-shanghai.aliyuncs.com/{username}"

# 指定docker hub

if [[ -z "$ver" ]]; then
  echo "version is empty"
  echo "usage: sh $0 <version>"
  exit 1
fi

# build image
echo "1. build docker image"
docker build -t ${serviceName}:${ver} -f docker/Dockerfile .

# docker tag local-image:tagname new-repo:tagname
echo "2. gen docker tag"
docker tag ${serviceName}:${ver} ${repoAddr}/${serviceName}:${ver}

# docker push new-repo:tagname
echo "3. push docker image to hub"
docker push ${repoAddr}/${serviceName}:${ver}

# deploy k8s deployment
echo "4. deploy k8s deployment"
kubectl apply -f k8s/go-deployment.yaml

# deploy k8s service
echo "5. deploy k8s service"
kubectl apply -f k8s/go-service.yaml

echo "Done. deploy success."