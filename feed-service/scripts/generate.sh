#!/bin/bash

set -euo pipefail

cd "$(dirname "$0")/.." || exit

mkdir -p ./internal/pb ./internal/pb/swagger

protoc -I ./api \
  -I ./api/google/api \
  --go_out=./internal/pb --go_opt=paths=source_relative \
  --go-grpc_out=./internal/pb --go-grpc_opt=paths=source_relative \
  ./api/feed_api/feed.proto ./api/models/feed_model.proto

protoc -I ./api \
  -I ./api/google/api \
  --grpc-gateway_out=./internal/pb \
  --grpc-gateway_opt paths=source_relative \
  --grpc-gateway_opt logtostderr=true \
  ./api/feed_api/feed.proto

protoc -I ./api \
  -I ./api/google/api \
  --openapiv2_out=./internal/pb/swagger \
  --openapiv2_opt logtostderr=true \
  ./api/feed_api/feed.proto


