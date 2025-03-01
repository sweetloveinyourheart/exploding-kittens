SHELL=bash

ALPINE_CONTAINER_IMAGE=alpine:3.20.3
GO_CONTAINER_IMAGE=golang:1.24.0-alpine

POSTGRES_CONTAINER_IMAGE=postgres:17.2-alpine3.19
NATS_CONTAINER_IMAGE=nats:2.10.25-alpine3.21

COMPOSE_PROJECT_NAME=kittens

include $(PWD)/scripts/_makefiles/build.mk
include $(PWD)/scripts/_makefiles/develop.mk
include $(PWD)/scripts/_makefiles/generate.mk
include $(PWD)/scripts/_makefiles/tests_unit.mk

export
