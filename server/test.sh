#!/usr/bin/env bash
docker build -t deploy-server:latest .
docker run --rm -p 8888:8000 -e PUSHER_APP_ID="497217" \
    -e PUSHER_KEY="65e885d06c3606bfce03" \
    -e PUSHER_SECRET="0a0135c2586d11a901a3" \
    -e PUSHER_CLUSTER="us2" \
    deploy-server:latest