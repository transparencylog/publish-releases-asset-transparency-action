# Start by building the application.
FROM golang:1.14 as build

WORKDIR /go/src/go.transparencylog.net/transparencylog-github-releases-action
COPY . .

ENV GO111MODULE=on
ENV GOFLAGS=-mod=vendor

# TODO: re-enable tests on all subpackages
# RUN go test -v ./...
RUN go install -v ./...

# Now copy it into our base image.
FROM gcr.io/distroless/base
COPY --from=build /go/bin/transparencylog-github-releases-action /

ENTRYPOINT ["/transparencylog-github-releases-action"]
