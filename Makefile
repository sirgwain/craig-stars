VERSION := 0.0.0-develop
COMMIT := `git rev-parse HEAD`

# detect os type and swap instructions accordingly
ifeq ($(OS),Windows_NT) 
# Forcibly swap to powershell if on windows to prevent make from using
# git's sh.eve instead (which is extremely limited in capabilities) 
SHELL := powershell.exe
.SHELLFLAGS := -NoProfile -Command
BUILDTIME := $(shell Get-Date -Format "yy.MM.dd HH:mm:ss")
BINARY_NAME := craig-stars.exe
# conditionals used to mimic behavior on unix-like systems  
mkdir = if ( -not ( Test-Path $(1) ) ) { mkdir "$(1)" }
rm = if ( Test-Path $(1) ) { rm -Recurse -Force "$(1)" }
cp = Copy-Item -Path "$(1)" -Destination "$(2)" -Force
download = Invoke-WebRequest -Uri "$(1)" -OutFile "$(2)"
unzip = Expand-Archive -LiteralPath "$(1)" -DestinationPath "./"
else
# Unix commands
BUILDTIME := $$(date +'%y.%m.%d %H:%M:%S')
BINARY_NAME := craig-stars
mkdir = mkdir -p $(1)
rm = rm -rf $(1)
cp = cp $(1) $(2)
download = wget -Q -P $(2) -B $(1)
unzip = unzip $(1)
endif

# replaces backslashes with unix-style frontslashes 
# and strips ending whitespace to allow tacking on backslashes later
goroot := $(subst \,/,$(shell go env GOROOT))

# defined separately to avoid backslash separators affecting recipes 
build_thing := go build \
	-o dist/${BINARY_NAME} \
	-ldflags \
	"-X 'github.com/sirgwain/craig-stars/cmd.semver=${VERSION}' \
	-X 'github.com/sirgwain/craig-stars/cmd.commit=${COMMIT}' \
	-X 'github.com/sirgwain/craig-stars/cmd.buildTime=${BUILDTIME}'" \
	main.go

# always redo these
.PHONY: run images build test clean dev dev_backend dev_frontend

# clean, deploy and launch all in 1 command
run: clean build dev 

images:
	cd frontend/static; $(call download,https://craig-stars.net/images/images.zip,images.zip);$(call unzip,images.zip);$(call rm,images.zip)
	# @sirgwain: pls remove the _MACOSX zip from the images zip file thx 
	$(call rm, frontend/static/_MACOSX.zip)

build: build_frontend tidy vendor generate build_wasm build_server

build_frontend:
	cd frontend; npm install; npm run build

build_server:
	$(call mkdir,dist)
	$(build_thing)

build_wasm:
	$(call mkdir,frontend/src/lib/wasm)
	go env -w GOOS=js GOARCH=wasm; go build -o frontend/src/lib/wasm/cs.wasm wasm/main.go
	$(call cp,$(goroot)/misc/wasm/wasm_exec.js,./frontend/src/lib/wasm/wasm_exec.js)
	go env -u GOOS GOARCH
	
# use docker to build an amd64 image for linux deployment
build_docker:
	docker build -f builder.Dockerfile --platform linux/amd64 . -t craig-stars-builder
	docker run -f builder.Dockerfile --platform linux/amd64 -v ${CURDIR}/dist:/dist craig-stars-builder

generate:
	go generate ./...

test:
	go test ./...
	cd frontend; npm run test

clean:
	go clean
	$(call rm,dist)
	$(call rm,vendor)
	$(call rm,frontend/build)

# uninstall unused modules
tidy:
	go mod tidy -v

# get those deps local!
vendor:
	go mod vendor

dev_frontend:
	cd frontend; npm run dev

dev_backend:
	air

dev:
	make -j 2 dev_backend dev_frontend
