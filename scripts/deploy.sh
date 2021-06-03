#!/bin/bash

COMMIT_TAG=$(git rev-parse --short HEAD)
REGISTRY=dkr.isi.nc
GROUP=incubator
APP_NAME=cine-circle-api

DOCKER_TAG=${REGISTRY}/${GROUP}/${APP_NAME}:${COMMIT_TAG}
IMAGE_NAME="${REGISTRY}\/${GROUP}\/${APP_NAME}:${COMMIT_TAG}"

#docker build . -t "${DOCKER_TAG}"
#docker push "${DOCKER_TAG}"
kubectl --context=bureau get deploy cine-circle-api -n incubator-isi -o yaml | tee deploy.yaml \
&& sed -i "s/image:.*$/image: ${IMAGE_NAME}/g" deploy.yaml \
&& kubectl replace -f deploy.yaml

rm deploy.yaml
