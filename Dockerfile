FROM golang:1.17 AS build-env
WORKDIR /build
COPY . /build
RUN go build -a -tags 'osusergo netgo static_build' -ldflags '-w -extldflags "-static"' -o app

FROM scratch
COPY --from=build-env /build/app /
WORKDIR /data
ENTRYPOINT ["/app"]