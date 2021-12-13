PROJECT_NAME = discord-inspirobot
GOARCH ?= $(shell uname -m | sed 's/^x86_64$$/amd64/')
GOOS ?= $(shell uname -s | sed 's/^Darwin$$/darwin/')

# Target: bin
# Build the binary for the current arch.
.PHONY: bin
bin:
	go build -o "bin/$(PROJECT_NAME)" "./cmd/$(PROJECT_NAME)"

.PHONY: run
run: bin
	"./bin/$(PROJECT_NAME)" -v 4 run

# Target: bin-[os]-[arch]
# Build the binary for the target arch.
.PHONY: bin-%
bin-%:
	@GOARCH="`echo "$*" | cut -d- -f2`" \
	GOOS="`echo "$*" | cut -d- -f1`" \
	go build -o "bin/$(PROJECT_NAME)-$*" "./cmd/$(PROJECT_NAME)"
