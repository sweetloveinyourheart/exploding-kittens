
ROOT_DIR=$(PWD)

build: # Build everything
	@make build-containers IMAGE_TAG=$(IMAGE_TAG)

build-containers:
	@make pocker-docker optionalReproFlag=$(optionalReproFlag)

# Base makefile target for building a binary
GOOS_OVERRIDE ?= GOOS=linux
build-binary:
	@echo "Building $(executablePath) with tag: $(IMAGE_TAG)"
	@cd $(directory) && \
	CGO_ENABLED=0 $(GOOS_OVERRIDE) $(extraArgs) go build -buildvcs=false -asmflags= -trimpath -ldflags "-buildid= -s -w $(LINKER_VERSION_FLAGS) -extldflags "-static"" && \
	cd $(ROOT_DIR) && \
	sha256sum $(executablePath)

# Base makefile target for building a docker image
build-docker:
	@DOCKER_BUILDKIT=1 docker build $(buildPlatform) \
	--ssh default \
	--target $(target) \
	--quiet \
	. \
	-t $(target):latest \
	--build-arg ALPINE_CONTAINER_IMAGE=$(ALPINE_CONTAINER_IMAGE) \
	--build-arg GO_CONTAINER_IMAGE=$(GO_CONTAINER_IMAGE) \
	$(additionalDockerArgs)

pocker-build:
	@make build-binary extraArgs=$(extraArgs) directory=cmd/planning-pocker executablePath=cmd/planning-pocker/planning-pocker

pocker-docker:
	@make pocker-build $(optionalReproFlag) extraArgs=$(extraArgs)
	@make build-docker buildPlatform=$(buildPlatorm) target=pocker