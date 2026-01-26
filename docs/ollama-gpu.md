# Ollama GPU Acceleration (Optional)

This guide explains how to enable GPU acceleration for the Ollama container used by the AionApi DEV stack.
CPU is the default for compatibility and is sufficient for most local development workflows.

## Prerequisites

- NVIDIA GPU with recent drivers installed on the host.
- NVIDIA Container Toolkit installed.

## Install NVIDIA Container Toolkit (Ubuntu)

```bash
distribution=$(. /etc/os-release; echo $ID$VERSION_ID)
curl -s -L https://nvidia.github.io/libnvidia-container/gpgkey | sudo apt-key add -
curl -s -L https://nvidia.github.io/libnvidia-container/$distribution/libnvidia-container.list \
  | sudo tee /etc/apt/sources.list.d/nvidia-container-toolkit.list

sudo apt-get update
sudo apt-get install -y nvidia-container-toolkit
sudo systemctl restart docker
```

Validate GPU access:

```bash
docker run --rm --gpus all nvidia/cuda:11.0-base nvidia-smi
```

## Enable GPU in the DEV compose file

Edit `infrastructure/docker/environments/dev/docker-compose-dev.yaml` and add the `deploy` block under the `ollama` service:

```yaml
  ollama:
    image: ollama/ollama:latest
    # ...
    deploy:
      resources:
        reservations:
          devices:
            - driver: nvidia
              count: 1
              capabilities: [gpu]
```

## Restart the DEV stack

```bash
make dev-down
make dev
```

## Verify

```bash
docker exec ollama-dev nvidia-smi
docker logs ollama-dev | grep -i cuda
```

If CUDA is available, you should see it in the container logs.
