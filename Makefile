# protogen runs a script to export all proto files
protogen:
	scripts/proto_gen.sh

# running unit tests
test:
	scripts/tests.sh

# compile the project into a main file
compile:
	go build -o main
	chmod +x main

# the following commands are used to start the cmds
CONFIG = configs/instance_1.yaml
CTL_C=configs/controller.yaml
DB_C=configs/database.yaml

node:
	./main node $(CONFIG)

controller:
	./main controller $(CTL_C)

db:
	./main database $(DB_C)
