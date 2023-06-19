BINARY_NAME=craig-stars

build_all: build_frontend build_server_all

build_frontend:
	yarn --cwd frontend install
	yarn --cwd frontend run build

build_server_all:
	mkdir -p dist
	GOARCH=amd64 GOOS=darwin go build -o dist/${BINARY_NAME}-darwin-amd64 main.go
	GOARCH=arm64 GOOS=darwin go build -o dist/${BINARY_NAME}-darwin-arm64 main.go
	GOARCH=amd64 GOOS=linux go build -o dist/${BINARY_NAME}-linux main.go
	GOARCH=amd64 GOOS=windows go build -o dist/${BINARY_NAME}-windows main.go
	lipo -create -output dist/${BINARY_NAME}-darwin-universal dist/${BINARY_NAME}-darwin-amd64 dist/${BINARY_NAME}-darwin-arm64

build:
	mkdir -p dist
	GOARCH=arm64 GOOS=darwin go build -o dist/${BINARY_NAME}-darwin-arm64 main.go

test:
	go test ./...

reset_db:
	rm data/data.db && go run main.go migrate && go run main.go create user -u admin -p admin

serve:
	./dist/${BINARY_NAME}-darwin-universal serve

build_and_run: build run

clean:
	go clean
	rm dist/${BINARY_NAME}-darwin-amd64
	rm dist/${BINARY_NAME}-darwin-arm64
	rm dist/${BINARY_NAME}-darwin-universal
	rm dist/${BINARY_NAME}-linux
	rm dist/${BINARY_NAME}-windows

# uninstall unused modules
tidy:
	go mod tidy -v

# get those deps local!
vendor:
	go mod vendor