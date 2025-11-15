FROM golang:1.22-alpine AS builder
WORKDIR /src
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/user-api .

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=builder /out/user-api /app/user-api
EXPOSE 8888
VOLUME ["/data/storage"]
ENTRYPOINT ["/app/user-api","-f","/app/etc/user-api.yaml"]