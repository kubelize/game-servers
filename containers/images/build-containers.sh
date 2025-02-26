
# base
podman build --platform linux/amd64 -t kubelize/game-servers:0.2.2-ln -f images/Dockerfile.ln
podman build --platform linux/amd64 -t kubelize/game-servers:0.2.2-we -f images/Dockerfile.we
podman push kubelize/game-servers:0.2.2-ln
podman push kubelize/game-servers:0.2.2-we

# sdtd
podman build --platform linux/amd64 -t kubelize/game-servers:0.2.9-sdtd -f images/Dockerfile.sdtd
podman push kubelize/game-servers:0.2.9-sdtd

# conan exiles
podman build --platform linux/amd64 -t kubelize/game-servers:0.1.0-ce -f images/Dockerfile.ce
podman push kubelize/game-servers:0.1.0-ce