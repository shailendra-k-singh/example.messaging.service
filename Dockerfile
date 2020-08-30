# build server
FROM golang:1.15 AS builder
LABEL stage=server-intermediate
WORKDIR /go/src/github.com/shailendra-k-singh/example.messaging.service
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -i -v -o bin/server ./cmd/server

# copy the built binary
FROM alpine:latest AS runner
LABEL app=server

COPY --from=builder /go/src/github.com/shailendra-k-singh/example.messaging.service/bin/server /bin/
WORKDIR /bin

ENTRYPOINT ["./server"]
