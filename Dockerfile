# Build Server
FROM golang:1.18 as builder
WORKDIR /workspace
COPY ./ ./
ARG VERSION
ARG GITVERSION
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on GOPROXY="https://goproxy.cn,direct"\
    go build -a -ldflags "-s -w -X github.com/mmfc-labs/driving-assistant/version.Version=${VERSION:-undefined} -X github.com/mmfc-labs/driving-assistant/version.GitRevision=${GITVERSION:-undefined}" \
    -o driving-assistant ./main.go


# driving-assistant image
FROM alpine:latest
WORKDIR /driving-assistant

# server
COPY --from=builder /workspace/driving-assistant .
COPY config.yaml ./config.yaml
COPY ./templates ./templates

EXPOSE 80 80

ENTRYPOINT ["./driving-assistant"]
