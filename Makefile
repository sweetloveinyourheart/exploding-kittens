SHELL=bash

ALPINE_CONTAINER_IMAGE=alpine:3.20.3
GO_CONTAINER_IMAGE=golang:1.23.2-alpine

COMPOSE_PROJECT_NAME=pocker

include $(PWD)/scripts/_makefiles/build.mk
include $(PWD)/scripts/_makefiles/develop.mk
include $(PWD)/scripts/_makefiles/generate.mk

export
