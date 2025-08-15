GHCR_URL=ghcr.io
EXPLODING_KITTENS_REPO=sweetloveinyourheart/exploding-kittens

BRANCHTAG=$(shell git rev-parse --abbrev-ref HEAD)
BRANCHTAG_SAFE=$(shell echo $(BRANCHTAG) | tr '[:upper:]' '[:lower:]' | tr '/' '_' | cut -c 1-32)
GIT_COMMIT_HASH=$(shell git rev-parse --short HEAD | tr '[:upper:]' '[:lower:]')
ifeq ("$(BRANCHTAG)", "HEAD")
	BRANCHTAG=$(shell git branch | grep HEAD | cut -d' ' -f 5 | tr -d '()' | tr '/' '_' | xargs)
	IMAGE_TAG=$(BRANCHTAG)
else
	IMAGE_TAG=$(shell date '+%Y-%m-%d-%H%M')_$(BRANCHTAG_SAFE)_$(GIT_COMMIT_HASH)
endif

auth-to-ghcr:
	@./scripts/github/auth.sh

push-to-ghcr:
	@docker tag $(IMAGE_NAME) $(REPOSITORY_URI):$(IMAGE_TAG)
	@docker push $(REPOSITORY_URI):$(IMAGE_TAG)

# Generate a new image tag
image-tag:
	@echo $(IMAGE_TAG)

kittens-docker-push:
	@make push-to-ghcr IMAGE_TAG=$(IMAGE_TAG) IMAGE_NAME=kittens:latest REPOSITORY_URI=$(GHCR_URL)/$(EXPLODING_KITTENS_REPO)
