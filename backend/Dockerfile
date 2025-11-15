FROM golang:1.24 AS builder
ENV GOTOOLCHAIN=auto
ARG GOPROXY=https://goproxy.cn,direct
ENV GOPROXY=${GOPROXY}
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/user-api .

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /out/user-api /app/user-api
EXPOSE 8888
VOLUME ["/data/storage"]
ENTRYPOINT ["/app/user-api","-f","/app/etc/user-api.yaml"]