build: clean
	./hack/build.sh
	mv bin/docsync ~/.local/go/bin/
	docsync -v
.PHONY: build

addlicense:
	@# install with `go install github.com/google/addlicense`
	addlicense -c 'rnrch' -l apache -v .
.PHONY: addlicense

clean:
	rm -rf bin
.PHONY: clean

test: build
	./hack/test.sh
.PHONY: test
