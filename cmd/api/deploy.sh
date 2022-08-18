#!/usr/bin/env sh

set -eux

flyctl auth docker

go run github.com/tailscale/mkctr@latest \
  --base="ghcr.io/tailscale/alpine-base:3.16" \
  --gopaths="github.com/contextwtf/lanyard/cmd/api:/usr/local/bin/api" \
  --ldflags="-X main.GitSha=`git rev-parse --short HEAD`" \
  --tags="latest" \
  --repos="registry.fly.io/al-prod" \
  --target=flyio \
  --push \
/usr/local/bin/api

flyctl deploy --detach \
	-i registry.fly.io/al-prod:latest \
	-a al-prod \
	-e "DD_ENV=production" \
	-e "DD_SERVICE=al-prod" \
	-e "DD_AGENT_HOST=datadog-agent.internal"
