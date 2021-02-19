build:
	go build -o docsync .

addlicense:
	# install with `go get github.com/google/addlicense`
	addlicense -c 'rnrch' -l apache -v .

clean:
	rm -f docsync

test: build
	./docsync -t test/test.tmpl -i ignore -o test/output.md -d test/test-folder
