GHCR_URL=ghcr.io
EXPLODING_KITTENS_REPO=sweetloveinyourheart

push-to-ghcr:
	@docker tag $(IMAGE_NAME) $(REPOSITORY_URI):$(IMAGE_TAG)
	@docker push $(REPOSITORY_URI):$(IMAGE_TAG)


# Generate a new image tag
image-tag:
	@echo $(IMAGE_TAG)

kittens-docker-push:
	@make push-to-ghcr IMAGE_TAG=$(IMAGE_TAG) IMAGE_NAME=kittens:latest REPOSITORY_URI=$(GHCR_URL)/$(EXPLODING_KITTENS_REPO)