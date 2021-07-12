FROM golang:1.16 as builder

# All these steps will be cached
RUN mkdir /build
WORKDIR /build
# <- COPY go.mod and go.sum files to the workspace
COPY go.mod .
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download

COPY . .

RUN go mod vendor
RUN go build -o mygrpc-server ./cmd/server

FROM ubuntu:latest

COPY --from=builder /build/mygrpc-server /usr/local/bin/

ENTRYPOINT ["mygrpc-server"]
