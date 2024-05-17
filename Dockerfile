FROM golang:1.22 as builder
WORKDIR /build

COPY . .
RUN CGO_ENABLED=0 go build -a -ldflags "-s -w" -o manager cmd/manager/main.go

FROM alpine:3
WORKDIR /workspace

RUN apk --no-cache add ca-certificates && update-ca-certificates

COPY --from=builder /build/manager manager
ENTRYPOINT ["./manager"]
