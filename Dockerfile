FROM --platform=$BUILDPLATFORM cgr.dev/chainguard/wolfi-base as build
RUN apk update --no-cache \
  && apk search -x go \
  && apk add \
    build-base \
    git \
    go=1.24.6-r1 \
    openssh

WORKDIR /srv/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -buildvcs=false -o untrack-that-url ./cmd/untrack-that-url

ENTRYPOINT [ "/srv/app/untrack-that-url" ]
