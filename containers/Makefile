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
REGISTRY ?= kubelize/game-servers
PLATFORM ?= linux/amd64
BUILDER ?= docker
TAG_SUFFIX ?= -alpha

# Version scheme: BASE_VERSION-{base_type}
BASE_VERSION ?= 0.3.7
# Version scheme: GAME_VERSION-{game_ext}
SDTD_VERSION ?= 0.4.0
PW_VERSION ?= 0.2.0
VH_VERSION ?= 0.2.0
CE_VERSION ?= 0.2.0
MC_VERSION ?= 0.2.0
HT_VERSION ?= 0.2.8

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
	go run generate-images.go -base-version=$(BASE_VERSION) -tag-suffix=$(TAG_SUFFIX)
	@echo "✓ Dockerfiles generated"

# Build all images (base + games)
build: build-base build-games
	@echo "✓ All images built"

# Build base images
build-base: $(BASE_TARGETS)
	@echo "✓ Base images built"

build-base-ln:
	@echo "→ Building base image: ln (Linux)"
	$(BUILDER) build --platform $(PLATFORM) \
		-t $(REGISTRY):$(BASE_VERSION)-ln$(TAG_SUFFIX) \
		-f images/Dockerfile.ln \
		.

build-base-we:
	@echo "→ Building base image: we (Wine)"
	$(BUILDER) build --platform $(PLATFORM) \
		-t $(REGISTRY):$(BASE_VERSION)-we$(TAG_SUFFIX) \
		-f images/Dockerfile.we \
		.

build-base-jv:
	@echo "→ Building base image: jv (Java)"
	$(BUILDER) build --platform $(PLATFORM) \
		-t $(REGISTRY):$(BASE_VERSION)-jv$(TAG_SUFFIX) \
		-f images/Dockerfile.jv \
		.

# Build game server images
build-games: $(GAME_TARGETS)
	@echo "✓ Game server images built"

build-game-sdtd:
	@echo "→ Building game: Seven Days to Die (sdtd)"
	$(BUILDER) build --platform $(PLATFORM) \
		-t $(REGISTRY):$(SDTD_VERSION)-sdtd$(TAG_SUFFIX) \
		-f images/Dockerfile.sdtd \
		.

build-game-pw:
	@echo "→ Building game: Palworld (pw)"
	$(BUILDER) build --platform $(PLATFORM) \
		-t $(REGISTRY):$(PW_VERSION)-pw$(TAG_SUFFIX) \
		-f images/Dockerfile.pw \
		.

build-game-vh:
	@echo "→ Building game: Valheim (vh)"
	$(BUILDER) build --platform $(PLATFORM) \
		-t $(REGISTRY):$(VH_VERSION)-vh$(TAG_SUFFIX) \
		-f images/Dockerfile.vh \
		.

build-game-ce:
	@echo "→ Building game: Conan Exiles (ce)"
	$(BUILDER) build --platform $(PLATFORM) \
		-t $(REGISTRY):$(CE_VERSION)-ce$(TAG_SUFFIX) \
		-f images/Dockerfile.ce \
		.

build-game-mc:
	@echo "→ Building game: Minecraft (mc)"
	$(BUILDER) build --platform $(PLATFORM) \
		-t $(REGISTRY):$(MC_VERSION)-mc$(TAG_SUFFIX) \
		-f images/Dockerfile.mc \
		.

build-game-ht:
	@echo "→ Building game: Hytale (ht)"
	$(BUILDER) build --platform $(PLATFORM) \
		-t $(REGISTRY):$(HT_VERSION)-ht$(TAG_SUFFIX) \
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
	$(BUILDER) push $(REGISTRY):$(BASE_VERSION)-ln$(TAG_SUFFIX)

push-base-we:
	@echo "→ Pushing base image: we"
	$(BUILDER) push $(REGISTRY):$(BASE_VERSION)-we$(TAG_SUFFIX)

push-base-jv:
	@echo "→ Pushing base image: jv"
	$(BUILDER) push $(REGISTRY):$(BASE_VERSION)-jv$(TAG_SUFFIX)

# Push game server images
push-games: $(PUSH_GAME_TARGETS)
	@echo "✓ Game server images pushed"

push-game-sdtd:
	@echo "→ Pushing game: Seven Days to Die"
	$(BUILDER) push $(REGISTRY):$(SDTD_VERSION)-sdtd$(TAG_SUFFIX)

push-game-pw:
	@echo "→ Pushing game: Palworld"
	$(BUILDER) push $(REGISTRY):$(PW_VERSION)-pw$(TAG_SUFFIX)

push-game-vh:
	@echo "→ Pushing game: Valheim"
	$(BUILDER) push $(REGISTRY):$(VH_VERSION)-vh$(TAG_SUFFIX)

push-game-ce:
	@echo "→ Pushing game: Conan Exiles"
	$(BUILDER) push $(REGISTRY):$(CE_VERSION)-ce$(TAG_SUFFIX)

push-game-mc:
	@echo "→ Pushing game: Minecraft"
	$(BUILDER) push $(REGISTRY):$(MC_VERSION)-mc$(TAG_SUFFIX)

push-game-ht:
	@echo "→ Pushing game: Hytale"
	$(BUILDER) push $(REGISTRY):$(HT_VERSION)-ht$(TAG_SUFFIX)

# Clean generated files
clean:
	@echo "→ Removing generated Dockerfiles..."
	rm -f images/Dockerfile.sdtd images/Dockerfile.pw images/Dockerfile.vh \
		images/Dockerfile.ce images/Dockerfile.mc images/Dockerfile.ht
	@echo "✓ Cleaned"

# Update version in generate-images.go
update-base-version:
	@echo "→ Updating base version to $(BASE_VERSION) in generate-images.go..."
	sed -i 's/kubelize\/game-servers:[0-9.]*-%s/kubelize\/game-servers:$(BASE_VERSION)-%s/g' ../generate-images.go
	@echo "✓ Base version updated"

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
	@echo "  BASE_VERSION=$(BASE_VERSION)"
	@echo ""
	@echo "Examples:"
	@echo "  make build BUILDER=docker           - Build with Docker instead of Podman"
	@echo "  make build-game-mc                  - Build only Minecraft image"
	@echo "  make BASE_VERSION=0.3.0 build-base  - Build base images with custom version"
