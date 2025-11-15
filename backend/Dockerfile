FROM golang:1.24 AS builder
ENV GOTOOLCHAIN=auto
ARG GOPROXY=https://goproxy.cn,direct
ENV GOPROXY=${GOPROXY}
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/user-api .

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=builder /out/user-api /app/user-api
EXPOSE 8888
VOLUME ["/data/storage"]
USER nonroot
ENTRYPOINT ["/app/user-api","-f","/app/etc/user-api.yaml"]