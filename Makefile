NAME := overseer
VERSION ?= 0.0.11

# run program
run: build
	@echo "run: running..."
	cd bin && ./${NAME}

run-%:
	@echo "run-$*: running..."
	@cd $* && go build -buildmode=pie -ldflags="-X main.Version=${VERSION} -s -w" -o ../bin/$*
	cd bin && ./$*

windows-run-%:
	@echo "windows-run-$*: running..."
	@cd $* && GOOS=windows go build -buildmode=pie -ldflags="-X main.Version=${VERSION} -s -w" -o ../bin/$*.exe
	cd bin && wine $*.exe

# build for local OS
build:
	cd overseer && go build -buildmode=pie -ldflags="-X main.Version=${VERSION} -s -w" -o ../bin/overseer main.go
	
build-mock:
	@echo "build-mock: building to bin/${NAME}..."
	@mkdir -p bin
	go build -o bin/mock mock/mock.go

docker-sim:
	@echo "docker-sim: running simulator of debian"
	docker run --rm -it -v ${PWD}/bin:/eqemu debian:stable-slim bash

run-test: build
	@echo "run-test: running..."
	cp bin/${NAME} ../server/build/bin/${NAME}
	cd ../server/build/bin && ./${NAME}
	
deq-%:
	@echo "deq-$*: running..."
	@cd $* && GOOS=linux GOARCH=amd64 go build -buildmode=pie -ldflags="-X main.Version=${VERSION} -s -w" -o ../bin/$*
	scp bin/$* deq@deq:/eqemu/
	ssh -t deq@deq "cd /eqemu && ./$*"

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
build-all:
	-rm bin/overseer-*.zip
	make build-all-overseer
	make build-all-diagnose
	make build-all-start
	make build-all-stop
	make build-all-install
	make build-all-update
	make build-all-verify

build-all-%:
	make build-$*-windows 
	make build-$*-linux

build-%-linux:
	@echo "build-$*-linux: ${VERSION}"
	@cd $* && GOOS=linux GOARCH=amd64 go build -ldflags="-X main.Version=${VERSION} -s -w" -o ../bin/$*
	cd bin && zip -r overseer-linux.zip $*
	@rm bin/$*

build-%-windows:
	@echo "build-$*-windows: ${VERSION}"
	@cd $* && GOOS=windows GOARCH=amd64 go build -ldflags -H=windowsgui -buildmode=pie -ldflags="-X main.Version=${VERSION} -s -w" -o ../bin/$*.exe
	cd bin && zip -r overseer-windows.zip $*.exe
	@rm bin/$*.exe

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
	go vet ./...
	staticcheck -go 1.14 ./...
	go test -tags ci -covermode=atomic -coverprofile=coverage.out ./...
    coverage=`go tool cover -func coverage.out | grep total | tr -s '\t' | cut -f 3 | grep -o '[^%]*'`

# CICD triggers this
set-version-%:
	@echo "VERSION=${VERSION}.$*" >> $$GITHUB_ENV

dev-copy-%:
	@echo "dev-copy-$*: building..."
	@cd $* && GOOS=linux GOARCH=amd64 go build -ldflags="-X main.Version=${VERSION} -s -w" -o ../bin/$*
	@echo "dev-copy-$*: copying to deq..."
	@scp bin/$* deq@deq-dev:/eqemu/