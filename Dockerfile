FROM golang:1.25 as builder
WORKDIR /build

COPY . .
RUN CGO_ENABLED=0 go build -a -ldflags "-s -w" -o manager cmd/manager/main.go

FROM alpine:3
WORKDIR /app

RUN apk --no-cache add ca-certificates && update-ca-certificates
RUN addgroup --gid 1000 app
RUN adduser --disabled-password --gecos "" --ingroup app --no-create-home --uid 1000 app

COPY --from=builder /build/manager /app/manager

USER 1000
ENTRYPOINT ["/app/manager"]
