OUT = gote
VERSION = `git describe --tags --long --dirty`
LDFLAGS = -ldflags "-X github.com/hawry/gote/cmd.buildVersion=$(VERSION)"
PRODLDFLAGS = -ldflags "-X github.com/hawry/gote/cmd.buildVersion=$(VERSION) -X main.logLevel=production"
DOCSLDFLAGS = -ldflags "-X github.com/hawry/gote/cmd.buildVersion=$(VERSION) -X main.createDocumentation=true"
PRODTAG = `git describe --tags --abbrev=0`
UNAME := $(shell uname)

ifeq ($(UNAME), MSYS_NT-10.0)
	OUT = gote.exe
endif

.PHONY: all
.SILENT:

all: default

default:
	go build $(LDFLAGS) -o $(OUT)

run: default
	./$(OUT) note -d

test:
	go test ./... -v

clean:
	rm -rf ./$(OUT); \
	rm -rf $(GOPATH)/bin/$(OUT)

prod:
	go build $(PRODLDFLAGS) -o $(OUT)

documentation:
	@echo "building with docs flag set! $(DOCSLDFLAGS)" \
	go build $(DOCSLDFLAGS) -o $(OUT)

debug:
	@echo 'build version will be $(VERSION)\n' \
	@echo "build prod tag will be $(PRODTAG)\n" \
	@echo "binary executable will be $(OUT)" \
	@echo "docs ld $(DOCSLDFLAGS)"

RELEASE_OUT = ./archives
U_ARCHS = amd64 arm64 386 arm
W_ARCHS = amd64 386

install:
	cp `pwd`/$(OUT) $(GOPATH)/bin/$(OUT)
	# unlink /usr/local/bin/$(OUT); \
	# ln -s `pwd`/$(OUT) /usr/local/bin/$(OUT)

# 386 arm64 arm

markdown:
	cd ./docs; \
	go run gendoc.go; \
	cd ../

release: linux windows

linux:
	@echo "**** Creating release archive for LINUX ***** "
	for arch in $(U_ARCHS); do \
		# echo "Building for $$arch"; \
		TARNAME="$(OUT)-$(PRODTAG)-linux-$$arch.tar.gz"; \
		echo "Building '$$TARNAME'"; \
		GOOS=linux GOARCH=$$arch go build $(PRODLDFLAGS) -o $(OUT); \
		tar -czvf ./$(RELEASE_OUT)/$$TARNAME ./$(OUT); \
		rm ./$(OUT); \
	done

windows:
	@echo "**** Creating release archive for WINDOWS ***** "
	for arch in $(W_ARCHS); do \
		TARNAME="$(OUT)-$(PRODTAG)-windows-$$arch.tar.gz"; \
		echo "Building '$$TARNAME'"; \
		GOOS=windows GOARCH=$$arch go build $(PRODLDFLAGS) -o $(OUT).exe; \
		tar -czvf ./$(RELEASE_OUT)/$$TARNAME ./$(OUT).exe; \
		rm ./$(OUT).exe; \
	done
