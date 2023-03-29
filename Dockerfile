FROM --platform=$BUILDPLATFORM golang:alpine3.17 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=0 \
    GOOS=$TARGETOS \
    GOARCH=$TARGETARCH \
    go build -o main .

FROM alpine:3.17

WORKDIR /app

COPY --from=build /app/main .

ENTRYPOINT ["/app/main"]