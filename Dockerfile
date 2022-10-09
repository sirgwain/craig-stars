# =============================================================================
#  Multi-stage Dockerfile Example
# =============================================================================
#  This is a simple Dockerfile that will build an image of scratch-base image.
#  Usage:
#    docker build -t simple:local . && docker run --rm simple:local
# =============================================================================

FROM node:alpine AS build-node

WORKDIR /workspace

COPY ./frontend /workspace/

RUN yarn install && yarn run build

# -----------------------------------------------------------------------------
#  Build Stage
# -----------------------------------------------------------------------------
FROM golang:alpine AS build

# Important:
#   Because this is a CGO enabled package, you are required to set it as 1.
ENV CGO_ENABLED=1

RUN apk add --no-cache \
    # Important: required for go-sqlite3
    gcc \
    # Required for Alpine
    musl-dev

WORKDIR /workspace

COPY . /workspace/

COPY --from=build-node /workspace/build /workspace/frontend/build

RUN \
    go mod tidy && \
    go install -ldflags='-s -w -extldflags "-static"' ./main.go


# -----------------------------------------------------------------------------
#  Main Stage
# -----------------------------------------------------------------------------
FROM busybox

WORKDIR /workspace

RUN mkdir -p /workspace/artifacts/craig-stars
RUN mkdir -p /workspace/artifacts/frontend
RUN mkdir -p /workspace/dist

# copy craig-stars binary for executing with ENTRYPOINT
COPY --from=build /go/bin/main /usr/local/bin/craig-stars 

# copy craig-stars binary into artifacts folder and tar it up
COPY --from=build /go/bin/main /workspace/artifacts/craig-stars/craig-stars
RUN tar -cvf /workspace/dist/craig-stars.tgz -C /workspace/artifacts/craig-stars/ .

# copy front end from node builder into artifacts and tar it up
COPY --from=build-node /workspace/build /workspace/artifacts/frontend/build
RUN tar -cvf /workspace/dist/frontend.tgz -C /workspace/artifacts/frontend/build .

CMD ["cp", "/workspace/dist/craig-stars.tgz", "/workspace/dist/frontend.tgz", "/dist"]

# ENTRYPOINT [ "/usr/local/bin/craig-stars" ]