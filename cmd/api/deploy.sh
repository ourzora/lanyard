#!/usr/bin/env sh

set -eux

flyctl auth docker

go run github.com/tailscale/mkctr@latest \
  --base="ghcr.io/tailscale/alpine-base:3.16" \
  --gopaths="github.com/contextart/al/cmd/api:/usr/local/bin/api" \
  --tags="latest" \
  --repos="registry.fly.io/al-prod" \
  --target=flyio \
  --push \
/usr/local/bin/api

flyctl deploy -i registry.fly.io/al-prod:latest
