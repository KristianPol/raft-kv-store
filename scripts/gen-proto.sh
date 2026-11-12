#!/usr/bin/env bash
# Needs: protoc, protoc-gen-go, protoc-gen-go-grpc on PATH.
set -euo pipefail
cd "$(dirname "$0")/.."
M=module=github.com/KristianPol/raft-kv-store
protoc --go_out=$M:. --go-grpc_out=$M:. proto/raft.proto
