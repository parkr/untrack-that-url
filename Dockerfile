FROM cgr.dev/chainguard/go:1.20.1 as build
WORKDIR /srv/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -buildvcs=false -o untrack-that-url ./cmd/untrack-that-url

ENTRYPOINT [ "/srv/app/untrack-that-url" ]
