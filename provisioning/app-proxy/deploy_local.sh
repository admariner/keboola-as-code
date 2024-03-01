#!/usr/bin/env bash
set -Eeuo pipefail

# CD to the script directory
cd "$(dirname "$0")"

# Setup local environment
export CLOUD_PROVIDER=local

# Default values for the local deployment
export MINIKUBE_PROFILE="${MINIKUBE_PROFILE:=app-proxy}"
export BUILD_BUILDID="${BUILD_BUILDID:=dev}"
export RELEASE_RELEASENAME="${RELEASE_RELEASENAME:=my-release}"
export KEBOOLA_STACK="${KEBOOLA_STACK:=local-machine}"
export HOSTNAME_SUFFIX="${HOSTNAME_SUFFIX:=keboola.com}"
export APP_PROXY_REPOSITORY="${APP_PROXY_REPOSITORY:=docker.io/keboola/app-proxy}" # docker.io prefix is required
export APP_PROXY_IMAGE_TAG="${APP_PROXY_IMAGE_TAG:=$(git rev-parse --short HEAD)}"
export APP_PROXY_REPLICAS="${APP_PROXY_REPLICAS:=3}"

./deploy.sh