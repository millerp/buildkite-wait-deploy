#!/usr/bin/env bash
docker build -t deploy-client:latest .
docker run --rm -e PUSHER_APP_KEY="65e885d06c3606bfce03" -e BUILDKITE_COMMIT="0a0135c2586d11a901a3" deploy-client:latest