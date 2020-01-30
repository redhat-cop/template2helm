LAST_TAG := $(shell git describe --abbrev=0 --tags 2>/dev/null || git rev-parse --short HEAD)

USER := redhat-cop
EXECUTABLE := template2helm

# only include the amd64 binaries, otherwise the github release will become
# too big
UNIX_EXECUTABLES := \
	darwin/amd64/$(EXECUTABLE) \
	freebsd/amd64/$(EXECUTABLE) \
	linux/amd64/$(EXECUTABLE)
WIN_EXECUTABLES := \
	windows/amd64/$(EXECUTABLE).exe

COMPRESSED_EXECUTABLES=$(UNIX_EXECUTABLES:%=%.bz2) $(WIN_EXECUTABLES:%.exe=%.zip)
COMPRESSED_EXECUTABLE_TARGETS=$(COMPRESSED_EXECUTABLES:%=bin/%)

UPLOAD_CMD = github-release upload -u $(USER) -r $(EXECUTABLE) -t $(LAST_TAG) -n $(subst /,-,$(FILE)) -f bin/$(FILE)

BUILD_ARGS = -ldflags "-X github.com/$(USER)/$(EXECUTABLE)/cmd.version=$(LAST_TAG)"

all: $(EXECUTABLE)

# arm
bin/linux/arm/5/$(EXECUTABLE):
	GOARM=5 GOARCH=arm GOOS=linux go build -o "$@" $(BUILD_ARGS)
bin/linux/arm/7/$(EXECUTABLE):
	GOARM=7 GOARCH=arm GOOS=linux go build -o "$@" $(BUILD_ARGS)

# 386
bin/darwin/386/$(EXECUTABLE):
	GOARCH=386 GOOS=darwin go build -o "$@" $(BUILD_ARGS)
bin/linux/386/$(EXECUTABLE):
	GOARCH=386 GOOS=linux go build -o "$@" $(BUILD_ARGS)
bin/windows/386/$(EXECUTABLE):
	GOARCH=386 GOOS=windows go build -o "$@" $(BUILD_ARGS)

# amd64
bin/freebsd/amd64/$(EXECUTABLE):
	GOARCH=amd64 GOOS=freebsd go build -o "$@" $(BUILD_ARGS)
bin/darwin/amd64/$(EXECUTABLE):
	GOARCH=amd64 GOOS=darwin go build -o "$@" $(BUILD_ARGS)
bin/linux/amd64/$(EXECUTABLE):
	GOARCH=amd64 GOOS=linux go build -o "$@" $(BUILD_ARGS)
bin/windows/amd64/$(EXECUTABLE).exe:
	GOARCH=amd64 GOOS=windows go build -o "$@" $(BUILD_ARGS)

# compressed artifacts, makes a huge difference (Go executable is ~9MB,
# after compressing ~2MB)
%.bz2: %
	bzip2 -c < "$<" > "$@"
%.zip: %.exe
	zip "$@" "$<"

# git tag -a v$(RELEASE) -m 'release $(RELEASE)'
release: clean
	$(MAKE) $(COMPRESSED_EXECUTABLE_TARGETS)
	git push && git push --tags
	git log --format=%B $(LAST_TAG) -1 | \
		github-release release -u $(USER) -r $(EXECUTABLE) \
		-t $(LAST_TAG) -n $(LAST_TAG) -d - || true
	$(foreach FILE,$(COMPRESSED_EXECUTABLES),$(UPLOAD_CMD);)

# install and/or update all dependencies, run this from the project directory
# go get -u ./...
# go test -i ./
dep:
	go mod vendor

$(EXECUTABLE): dep
	go build -o "$@" $(BUILD_ARGS)

install:
	go install

clean:
	rm go-app || true
	rm $(EXECUTABLE) || true
	rm -rf bin/

test: clean dep
	go test ./...

test_e2e:
	test/e2e.sh

lint:
	golangci-lint run

.PHONY: clean release dep install test test_e2e lint
