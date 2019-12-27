#!/usr/bin/env bash

build(){
local DOCKER_DIR=docker
local COMPONENT=cart
local TAG=local
local DOCKER_REGISTRY=registry.local/kamilkoduo/digicart/${COMPONENT}
local DOCKER_PATH=${DOCKER_DIR}/${COMPONENT}/Dockerfile
local CONTEXT=.
docker build -t ${DOCKER_REGISTRY}:${TAG} -f ${DOCKER_PATH} ${CONTEXT} --no-cache
}
build