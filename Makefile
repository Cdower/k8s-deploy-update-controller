#!make
.PHONY: build all
.DEFAULT_GOAL := all
include .env
include secret_envfile

all: build

build: buildx-build

buildx-build:
	docker buildx use mybuilder
	docker buildx build --push --platform "linux/arm64,linux/amd64" --tag ${IMAGEFULLNAME} -f ./build/Dockerfile .

# DEPLOY #cdower/deploy-update-controller:v0.0.1

plan:
	terraform -chdir=./deployments/terraform plan -var="image=${IMAGEFULLNAME}" -var="registry_user=${REGISTRY_USER}" -var="registry_pass=${REGISTRY_PASS}"

apply:
	terraform -chdir=./deployments/terraform apply -var="image=${IMAGEFULLNAME}" -var="registry_user=${REGISTRY_USER}" -var="registry_pass=${REGISTRY_PASS}"

# SETUP

create-context:
	docker context create --docker host=unix:///var/run/docker.sock --kubernetes config-file=${HOME}/.kube/config node-arm64
	docker context create --docker host=unix:///var/run/docker.sock --kubernetes config-file=${HOME}/.kube/config node-amd64

buildx-create:
	docker buildx create --platform linux/arm64,linux/arm/v8 --name mybuilder node-arm64
	docker buildx create --append --platform linux/amd64 --name mybuilder node-amd64

tag: tag-create tag-push

tag-create:
	git tag -a ${VERSION}

tag-push:
	git push origin ${VERSION}

# DEPRECATED
manif-push: manif-build push-arm64 push-amd64 manifest

manif-build: build-arm64 build-amd64

push-amd64:
	docker push ${IMAGEFULLNAME}-amd64

build-amd64:
	docker build -f ./build/Dockerfile -t ${IMAGEFULLNAME}-manif-amd64 --build-arg TARGETPLATFORM=linux/amd64 --build-arg BUILDPLATFORM=linux/amd64 .

push-arm64:
	docker push ${IMAGEFULLNAME}-arm64

build-arm64:
	docker build -f ./build/Dockerfile -t ${IMAGEFULLNAME}-manif-arm64 --build-arg TARGETPLATFORM=linux/arm64 --build-arg BUILDPLATFORM=linux/arm64 .

manifest:
	docker manifest rm ${IMAGEFULLNAME}-manif
	docker manifest create ${IMAGEFULLNAME} --amend ${IMAGEFULLNAME}-manif-amd64 --amend ${IMAGEFULLNAME}-manif-arm64
	docker manifest push ${IMAGEFULLNAME}-manif