MAKE_HOME=${PWD}

.PHONY: test

test:
	go test ./... -v
