# [gesquive/dyngo](https://github.com/gesquive/dyngo)
[![Software License](https://img.shields.io/badge/License-MIT-orange.svg?style=flat-square)](https://github.com/gesquive/dyngo/blob/master/LICENSE)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/gesquive/dyngo)
[![Pipeline Status](https://img.shields.io/gitlab/pipeline/gesquive/dyngo?style=flat-square)](https://gitlab.com/gesquive/dyngo/pipelines)
[![Coverage Report](https://gitlab.com/gesquive/dyngo/badges/master/coverage.svg?style=flat-square)](https://gesquive.gitlab.io/dyngo/coverage.html)
[![Github Release](https://img.shields.io/github/v/tag/gesquive/dyngo?style=flat-square)](https://github.com/gesquive/dyngo)

# Supported Architectures

This image supports multiple architectures:

- `amd64`, `x86-64`
- `armv7`, `armhf`
- `arm64`, `aarch64`

Docker images are uploaded with using Docker manifest lists to make multi-platform deployments easer. More info can be found from [Docker](https://github.com/docker/distribution/blob/master/docs/spec/manifest-v2-2.md#manifest-list)

You can simply pull the image using `gesquive/dyngo` and docker should retreive the correct image for your architecture.

# Supported Tags
If you want a specific version of `dyngo` you can pull it by specifying a version tag.

## Version Tags
This image provides versions that are available via tags. 

| Tag    | Description |
| ------ | ----------- |
| `latest` | Latest stable release |
| `0.9.0`  | Stable release v0.9.0 |
| `0.9.0-<git_hash>` | Development preview of version v0.9.0 at the given git hash |

# Usage

Here are some example snippets to help you get started creating a container.
docker

## Docker CLI

```shell
docker create \
  --name=dyngo \
  -v path/to/config:/config \
  --restart unless-stopped \
  gesquive/dyngo
```

## Docker Compose
Compatible with docker-compose v2 schemas.

```docker
---
version: "2"
services:
  dyngo:
    image: gesquive/dyngo
    container_name: dyngo
    volumes:
      - path/to/config:/config
    restart: unless-stopped
```
# Parameters
The container defines the following parameters that you can set:

| Parameter | Function |
| --------- | -------- |
| `-v /config`  | The dyngo config goes here |
