#!/bin/bash

# setting variables which are used
# as protoc input arguments.
SRC_DIR="proto/"
DST_DIR="pkg/"
PROTO="transactions"

# make sure to have protoc installed on your system
protoc -I=$SRC_DIR --go_out=$DST_DIR $SRC_DIR/$PROTO.proto || echo "protoc not installed on this machine!"
protoc -I=$SRC_DIR --go-grpc_out=$DST_DIR $SRC_DIR/$PROTO.proto || echo "protoc not installed on this machine!"
