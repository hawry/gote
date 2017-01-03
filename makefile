OUT = gote
VERSION = `git describe --tags --long --dirty`
LDFLAGS = -ldflags "-X github.com/hawry/gote/cmd.buildVersion=$(VERSION)"

.PHONY: all
.SILENT:

all: default run

default:
	go build $(LDFLAGS) -o $(OUT)

run:
	./$(OUT)

clean:
	rm -rf ./$(OUT); \
	rm -rf $(GOPATH)/bin/$(OUT)

debug:
	@echo "build version will be $(VERSION)"

RELEASE_OUT = ./archives
U_ARCHS = amd64 arm64 386 arm
W_ARCHS = amd64 386

install:
	cp `pwd`/$(OUT) $(GOPATH)/bin/$(OUT)
	# unlink /usr/local/bin/$(OUT); \
	# ln -s `pwd`/$(OUT) /usr/local/bin/$(OUT)

# 386 arm64 arm
release:
	@echo "**** Creating release archive ***** "
	for arch in $(U_ARCHS); do \
		# echo "Building for $$arch"; \
		TARNAME="$(OUT)-linux-$$arch.tar.gz"; \
		echo "Building '$$TARNAME'"; \
		GOOS=linux GOARCH=$$arch go build $(LDFLAGS) -o $(OUT); \
		tar -czvf ./$(RELEASE_OUT)/$$TARNAME ./$(OUT); \
		rm ./$(OUT); \
	done
	for arch in $(W_ARCHS); do \
		TARNAME="$(OUT)-windows-$$arch.tar.gz"; \
		echo "Building '$$TARNAME'"; \
		GOOS=windows GOARCH=$$arch go build $(LDFLAGS) -o $(OUT).exe; \
		tar -czvf ./$(RELEASE_OUT)/$$TARNAME ./$(OUT).exe; \
		rm ./$(OUT).exe; \
	done
