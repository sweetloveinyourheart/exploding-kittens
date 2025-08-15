#!/bin/bash
. ./scripts/util.sh

set -euo pipefail

REGISTRY="ghcr.io"
USERNAME="${GITHUB_ACTOR:-sweetloveinyourheart}"

# Use GITHUB_TOKEN in CI, GHCR_TOKEN locally
if [[ -n "${GITHUB_ACTIONS:-}" ]]; then
  echo "Logging in to GHCR using GITHUB_TOKEN..."
  echo "${GITHUB_TOKEN}" | docker login "$REGISTRY" -u "$USERNAME" --password-stdin
else
  if [[ -z "${GHCR_TOKEN:-}" ]]; then
    echo "Error: GHCR_TOKEN environment variable not set for local login" >&2
    exit 1
  fi
  echo "Logging in to GHCR locally..."
  echo "${GHCR_TOKEN}" | docker login "$REGISTRY" -u "$USERNAME" --password-stdin
fi