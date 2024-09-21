# Build
FROM golang:latest as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

COPY pkg ./pkg
COPY cmd ./cmd
COPY configs ./configs
COPY internal ./internal

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -a -o oss-codenames-api cmd/api/main.go

# Package
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/oss-codenames-api .

CMD ["./oss-codenames-api"] 