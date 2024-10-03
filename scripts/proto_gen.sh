#!/bin/bash

# setting env variables which are used
# as protoc input arguments.
export $SRC_DIR=proto/
export $DST_DIR=pkg/
export $PROTO=transactions

# make sure to have protoc installed on your system
protoc -I=$SRC_DIR --go_out=$DST_DIR $SRC_DIR/$PROTO.proto || echo "protoc not installed on this machine!"
