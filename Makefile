NAME := overseer
VERSION ?= 0.0.2

# run program
run: build
	@echo "run: running..."
	cd bin && ./${NAME}

run-%:
	@echo "run-$*: running..."
	@cd $* && GOOS=darwin GOARCH=amd64 go build -buildmode=pie -ldflags="-X main.Version=${VERSION} -s -w" -o ../bin/$*
	cd bin && ./$*

# build for local OS
build:
	cd overseer && GOOS=darwin GOARCH=amd64 go build -buildmode=pie -ldflags="-X main.Version=${VERSION} -s -w" -o ../bin/overseer main.go
	
build-mock:
	@echo "build-mock: building to bin/${NAME}..."
	@mkdir -p bin
	go build -o bin/mock mock/mock.go


run-test: build
	@echo "run-test: running..."
	cp bin/${NAME} ../server/build/bin/${NAME}
	cd ../server/build/bin && ./${NAME}
	
# bundle program with windows icon
bundle:
	@echo "if go-winres is not found, run go install github.com/tc-hib/go-winres@latest"
	@echo "bundle: setting ${NAME} icon"
	go-winres simply --icon ${NAME}.png

# run tests that aren't flagged for SINGLE_TEST
test:
	@echo "test: running tests..."
	@go test ./...

# build all supported os's
build-all: build-overseer-darwin build-overseer-windows build-overseer-linux build-diagnose-darwin build-diagnose-windows build-diagnose-linux build-bootstrap-darwin build-bootstrap-windows build-bootstrap-linux

build-%-darwin:
	@echo "build-$*-darwin: ${VERSION}"
	@cd $* && GOOS=darwin GOARCH=amd64 go build -buildmode=pie -ldflags="-X main.Version=${VERSION} -s -w" -o ../bin/$*
	cd bin && zip -r overseer-darwin-${VERSION}.zip $*
	@rm bin/$*

build-%-linux:
	@echo "build-$*-linux: ${VERSION}"
	@cd $* && GOOS=linux GOARCH=amd64 go build -ldflags="-X main.Version=${VERSION} -s -w" -o ../bin/$*
	cd bin && zip -r overseer-linux-${VERSION}.zip $*
	@rm bin/$*

build-%-windows:
	@echo "build-$*-windows: ${VERSION}"
	@cd $* && GOOS=windows GOARCH=amd64 go build -buildmode=pie -ldflags="-X main.Version=${VERSION} -s -w" -o ../bin/$*.exe
	cd bin && zip -r overseer-windows-${VERSION}.zip $*.exe
	@rm bin/$*.exe

# used by xackery, build darwin copy and move to blender path
build-copy: build-darwin
	@echo "copying to ${NAME}-addons..."
	cp bin/${NAME}-darwin "/Users/xackery/Library/Application Support/Blender/3.4/scripts/addons/${NAME}-addon/${NAME}-darwin"

# run pprof and dump 3 snapshots of heap
profile-heap:
	@echo "profile-heap: running pprof watcher for 2 minutes with snapshots 0 to 3..."
	@-mkdir -p bin
	curl http://localhost:8082/debug/pprof/heap > bin/heap.0.pprof
	sleep 30
	curl http://localhost:8082/debug/pprof/heap > bin/heap.1.pprof
	sleep 30
	curl http://localhost:8082/debug/pprof/heap > bin/heap.2.pprof
	sleep 30
	curl http://localhost:8082/debug/pprof/heap > bin/heap.3.pprof

# peek at a heap
profile-heap-%:
	@echo "profile-heap-$*: use top20, svg, or list *word* for pprof commands, ctrl+c when done"
	go tool pprof bin/heap.$*.pprof

# run a trace on ${NAME}
profile-trace:
	@echo "profile-trace: getting trace data, this can show memory leaks and other issues..."
	curl http://localhost:8082/debug/pprof/trace > bin/trace.out
	go tool trace bin/trace.out

# run sanitization against golang
sanitize:
	@echo "sanitize: checking for errors"
	rm -rf vendor/
	go vet -tags ci ./...
	test -z $(goimports -e -d . | tee /dev/stderr)
	-go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	gocyclo -over 99 .
	golint -set_exit_status $(go list -tags ci ./...)
	staticcheck -go 1.14 ./...
	go test -tags ci -covermode=atomic -coverprofile=coverage.out ./...
    coverage=`go tool cover -func coverage.out | grep total | tr -s '\t' | cut -f 3 | grep -o '[^%]*'`

# CICD triggers this
set-version-%:
	@echo "VERSION=${VERSION}.$*" >> $$GITHUB_ENV
