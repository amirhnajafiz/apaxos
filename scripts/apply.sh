#!/bin/bash

# set instance config file
config="configs/instance_$1.yaml"

echo "applying $config ..."

# execute the instance
./main node "$config"
