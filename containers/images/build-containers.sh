
podman build --platform linux/amd64 -t kubelize/game-servers:0.2.1-ln -f Dockerfile.ln
podman build --platform linux/amd64 -t kubelize/game-servers:0.2.0-we -f Dockerfile.we
podman push kubelize/game-servers:0.2.1-ln
podman push kubelize/game-servers:0.2.0-we
podman build --platform linux/amd64 -t kubelize/game-servers:0.2.9-sdtd -f images/Dockerfile.sdtd
podman push kubelize/game-servers:0.2.9-sdtd
