BINARY_NAME=craig-stars

# always redo these
.PHONY: build test clean dev dev_backend dev_frontend

build: build_frontend tidy vendor build_server

build_frontend:
	cd frontend; npm install
	cd frontend; npm run build

build_server:
	mkdir -p dist
	go build -o dist/${BINARY_NAME} main.go

# use docker to build an amd64 image for linux deployment
build_docker:
	docker build -f builder.Dockerfile --platform linux/amd64 . -t craig-stars-builder
	docker run -f builder.Dockerfile --platform linux/amd64 -v ${CURDIR}/dist:/dist craig-stars-builder

test:
	go test ./...
	cd frontend; npm run test

clean:
	go clean
	rm -rf dist
	rm -rf vendor
	rm -rf frontend/build

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

