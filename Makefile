SHELL=bash

ALPINE_CONTAINER_IMAGE=alpine:3.20.3
GO_CONTAINER_IMAGE=golang:1.23.2-alpine

POSTGRES_CONTAINER_IMAGE=postgres:17.2-alpine3.19

COMPOSE_PROJECT_NAME=kittens

include $(PWD)/scripts/_makefiles/build.mk
include $(PWD)/scripts/_makefiles/develop.mk
include $(PWD)/scripts/_makefiles/generate.mk
include $(PWD)/scripts/_makefiles/tests_unit.mk

export
