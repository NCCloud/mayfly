FROM golang:1.19 as builder
ARG GITHUB_TOKEN

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN git config --global url."https://${GITHUB_TOKEN}:@github.com/".insteadOf "https://github.com/"
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags "-s -w" -o mayfly cmd/main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /app

COPY --from=builder /workspace .

ENTRYPOINT ["/app/mayfly"]
