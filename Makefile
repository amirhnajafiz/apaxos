protogen:
	scripts/proto_gen.sh

test:
	scripts/tests.sh

compile:
	go build -o main
	chmod +x main

CONFIG = config.yaml
CSV = testcase.csv

node:
	./main node --config $(CONFIG)

controller:
	./main controller --config $(CONFIG) --csv ignore

testcase:
	./main controller --config $(CONFIG) --csv $(CSV)
