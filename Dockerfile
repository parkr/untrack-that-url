FROM --platform=$BUILDPLATFORM cgr.dev/chainguard/wolfi-base as build
RUN apk update && apk search -x go && apk add build-base git openssh go=1.20.6
WORKDIR /srv/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -buildvcs=false -o untrack-that-url ./cmd/untrack-that-url

ENTRYPOINT [ "/srv/app/untrack-that-url" ]
