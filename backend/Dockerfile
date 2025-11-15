FROM golang:1.24-alpine AS builder
ENV GOTOOLCHAIN=auto
ARG ALPINE_MIRROR=https://dl-cdn.alpinelinux.org/alpine
WORKDIR /src
RUN sed -i "s|https://dl-cdn.alpinelinux.org/alpine|${ALPINE_MIRROR}|g" /etc/apk/repositories && apk update && apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/user-api .

FROM alpine:3.20
ARG ALPINE_MIRROR=https://dl-cdn.alpinelinux.org/alpine
WORKDIR /app
RUN sed -i "s|https://dl-cdn.alpinelinux.org/alpine|${ALPINE_MIRROR}|g" /etc/apk/repositories && apk update && apk add --no-cache ca-certificates tzdata && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=builder /out/user-api /app/user-api
EXPOSE 8888
VOLUME ["/data/storage"]
ENTRYPOINT ["/app/user-api","-f","/app/etc/user-api.yaml"]