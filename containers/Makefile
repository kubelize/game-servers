# Makefile for building and pushing game server container images
# Usage:
#   make all              - Generate Dockerfiles, build all images, and push to registry
#   make generate         - Generate Dockerfiles from template
#   make build-base       - Build base images only
#   make build-games      - Build game server images only
#   make build            - Build all images
#   make push             - Push all images to registry
#   make clean            - Remove generated Dockerfiles

# Configuration
REGISTRY ?= kubelize/games
PLATFORM ?= linux/amd64
BUILDER ?= docker
TAG_SUFFIX ?= -alpha

# Single release version for all images
# Image tags become:
#   $(RELEASE_VERSION)-ln$(TAG_SUFFIX), $(RELEASE_VERSION)-we$(TAG_SUFFIX), $(RELEASE_VERSION)-jv$(TAG_SUFFIX)
#   $(RELEASE_VERSION)-sdtd$(TAG_SUFFIX), $(RELEASE_VERSION)-ce$(TAG_SUFFIX), ...
RELEASE_VERSION ?= 0.0.1

# Base images
BASE_IMAGES = ln we jv
BASE_TARGETS = $(addprefix build-base-,$(BASE_IMAGES))
PUSH_BASE_TARGETS = $(addprefix push-base-,$(BASE_IMAGES))

# Game images
GAME_IMAGES = sdtd pw vh ce mc ht
GAME_TARGETS = $(addprefix build-game-,$(GAME_IMAGES))
PUSH_GAME_TARGETS = $(addprefix push-game-,$(GAME_IMAGES))

.PHONY: all generate build build-base build-games push push-base push-games clean help build-gamekeeper

# Default target
all: build-gamekeeper generate build push
	@echo "✓ All images built and pushed successfully"

# Build GameKeeper binary for Linux
build-gamekeeper:
	@echo "→ Building GameKeeper binary for Linux..."
	cd gamekeeper && $(MAKE) build
	@echo "✓ GameKeeper binary built"

# Generate Dockerfiles from template
generate:
	@echo "→ Generating Dockerfiles from template..."
	go run generate-images.go -registry=$(REGISTRY) -base-version=$(RELEASE_VERSION) -tag-suffix=$(TAG_SUFFIX)
	@echo "✓ Dockerfiles generated"

# Build all images (base + games)
build: build-base build-games
	@echo "✓ All images built"

# Build base images
build-base: build-gamekeeper $(BASE_TARGETS)
	@echo "✓ Base images built"

build-base-ln: build-gamekeeper
	@echo "→ Building base image: ln (Linux)"
	$(BUILDER) build --platform $(PLATFORM) \
		-t $(REGISTRY):$(RELEASE_VERSION)-ln$(TAG_SUFFIX) \
		-f images/Dockerfile.ln \
		.

build-base-we: build-gamekeeper
	@echo "→ Building base image: we (Wine)"
	$(BUILDER) build --platform $(PLATFORM) \
		-t $(REGISTRY):$(RELEASE_VERSION)-we$(TAG_SUFFIX) \
		-f images/Dockerfile.we \
		.

build-base-jv: build-gamekeeper
	@echo "→ Building base image: jv (Java)"
	$(BUILDER) build --platform $(PLATFORM) \
		-t $(REGISTRY):$(RELEASE_VERSION)-jv$(TAG_SUFFIX) \
		-f images/Dockerfile.jv \
		.

# Build game server images
build-games: $(GAME_TARGETS)
	@echo "✓ Game server images built"

build-game-sdtd:
	@echo "→ Building game: Seven Days to Die (sdtd)"
	$(BUILDER) build --platform $(PLATFORM) \
		-t $(REGISTRY):$(RELEASE_VERSION)-sdtd$(TAG_SUFFIX) \
		-f images/Dockerfile.sdtd \
		.

build-game-pw:
	@echo "→ Building game: Palworld (pw)"
	$(BUILDER) build --platform $(PLATFORM) \
		-t $(REGISTRY):$(RELEASE_VERSION)-pw$(TAG_SUFFIX) \
		-f images/Dockerfile.pw \
		.

build-game-vh:
	@echo "→ Building game: Valheim (vh)"
	$(BUILDER) build --platform $(PLATFORM) \
		-t $(REGISTRY):$(RELEASE_VERSION)-vh$(TAG_SUFFIX) \
		-f images/Dockerfile.vh \
		.

build-game-ce:
	@echo "→ Building game: Conan Exiles (ce)"
	$(BUILDER) build --platform $(PLATFORM) \
		-t $(REGISTRY):$(RELEASE_VERSION)-ce$(TAG_SUFFIX) \
		-f images/Dockerfile.ce \
		.

build-game-mc:
	@echo "→ Building game: Minecraft (mc)"
	$(BUILDER) build --platform $(PLATFORM) \
		-t $(REGISTRY):$(RELEASE_VERSION)-mc$(TAG_SUFFIX) \
		-f images/Dockerfile.mc \
		.

build-game-ht:
	@echo "→ Building game: Hytale (ht)"
	$(BUILDER) build --platform $(PLATFORM) \
		-t $(REGISTRY):$(RELEASE_VERSION)-ht$(TAG_SUFFIX) \
		-f images/Dockerfile.ht \
		.

# Push all images
push: push-base push-games
	@echo "✓ All images pushed"

# Push base images
push-base: $(PUSH_BASE_TARGETS)
	@echo "✓ Base images pushed"

push-base-ln:
	@echo "→ Pushing base image: ln"
	$(BUILDER) push $(REGISTRY):$(RELEASE_VERSION)-ln$(TAG_SUFFIX)

push-base-we:
	@echo "→ Pushing base image: we"
	$(BUILDER) push $(REGISTRY):$(RELEASE_VERSION)-we$(TAG_SUFFIX)

push-base-jv:
	@echo "→ Pushing base image: jv"
	$(BUILDER) push $(REGISTRY):$(RELEASE_VERSION)-jv$(TAG_SUFFIX)

# Push game server images
push-games: $(PUSH_GAME_TARGETS)
	@echo "✓ Game server images pushed"

push-game-sdtd:
	@echo "→ Pushing game: Seven Days to Die"
	$(BUILDER) push $(REGISTRY):$(RELEASE_VERSION)-sdtd$(TAG_SUFFIX)

push-game-pw:
	@echo "→ Pushing game: Palworld"
	$(BUILDER) push $(REGISTRY):$(RELEASE_VERSION)-pw$(TAG_SUFFIX)

push-game-vh:
	@echo "→ Pushing game: Valheim"
	$(BUILDER) push $(REGISTRY):$(RELEASE_VERSION)-vh$(TAG_SUFFIX)

push-game-ce:
	@echo "→ Pushing game: Conan Exiles"
	$(BUILDER) push $(REGISTRY):$(RELEASE_VERSION)-ce$(TAG_SUFFIX)

push-game-mc:
	@echo "→ Pushing game: Minecraft"
	$(BUILDER) push $(REGISTRY):$(RELEASE_VERSION)-mc$(TAG_SUFFIX)

push-game-ht:
	@echo "→ Pushing game: Hytale"
	$(BUILDER) push $(REGISTRY):$(RELEASE_VERSION)-ht$(TAG_SUFFIX)

# Clean generated files
clean:
	@echo "→ Removing generated Dockerfiles..."
	rm -f images/Dockerfile.sdtd images/Dockerfile.pw images/Dockerfile.vh \
		images/Dockerfile.ce images/Dockerfile.mc images/Dockerfile.ht
	@echo "✓ Cleaned"

# Help target
help:
	@echo "Game Server Container Build System"
	@echo ""
	@echo "Targets:"
	@echo "  make all              - Generate Dockerfiles, build all images, and push"
	@echo "  make generate         - Generate Dockerfiles from template"
	@echo "  make build            - Build all images (base + games)"
	@echo "  make build-base       - Build base images only"
	@echo "  make build-games      - Build game server images only"
	@echo "  make push             - Push all images to registry"
	@echo "  make push-base        - Push base images only"
	@echo "  make push-games       - Push game server images only"
	@echo "  make clean            - Remove generated Dockerfiles"
	@echo "  make help             - Show this help message"
	@echo ""
	@echo "Configuration (can be overridden):"
	@echo "  REGISTRY=$(REGISTRY)"
	@echo "  PLATFORM=$(PLATFORM)"
	@echo "  BUILDER=$(BUILDER)"
	@echo "  RELEASE_VERSION=$(RELEASE_VERSION)"
	@echo "  TAG_SUFFIX=$(TAG_SUFFIX)"
	@echo ""
	@echo "Examples:"
	@echo "  make build BUILDER=docker           - Build with Docker instead of Podman"
	@echo "  make build-game-mc                  - Build only Minecraft image"
	@echo "  make RELEASE_VERSION=0.3.0 build    - Build all images with a shared version"
