#!/bin/bash

# setting variables which are used
# as protoc input arguments.
SRC_DIR="proto/"
DST_DIR="pkg/"
PROTO_FILES=("transactions" "liveness" "apaxos")

# make sure to have protoc installed on your system
for PROTO in "${PROTO_FILES[@]}"; do
  echo "generating Go code for $PROTO.proto ..."
  protoc -I=$SRC_DIR --go_out=$DST_DIR $SRC_DIR/$PROTO.proto || echo "protoc not installed on this machine!"
  protoc -I=$SRC_DIR --go-grpc_out=$DST_DIR $SRC_DIR/$PROTO.proto || echo "protoc not installed on this machine!"
done
