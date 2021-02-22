build: clean
	@echo "Begin to build binary"
	./hack/build.sh
.PHONY: build

addlicense:
	# install with `go get github.com/google/addlicense`
	addlicense -c 'rnrch' -l apache -v .
.PHONY: addlicense

clean:
	rm -rf bin
.PHONY: clean

test: build
	@echo "Begin to run test"
	./hack/test.sh
.PHONY: test
