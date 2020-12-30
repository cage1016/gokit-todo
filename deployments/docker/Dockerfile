FROM gcr.io/gcp-runtimes/go1-builder:1.14 AS builder

WORKDIR /go/src/github.com/cage1016/gokit-todo

RUN apt-get update && apt-get install -y make git

ENV GO111MODULE=on

COPY . .

RUN /usr/local/go/bin/go build -o app ./cmd/todo

# Application image.
FROM gcr.io/distroless/base:latest
COPY --from=builder /go/src/github.com/cage1016/gokit-todo/app /usr/local/bin/app

CMD ["/usr/local/bin/app"]