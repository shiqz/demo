# syntax=docker/dockerfile:1
FROM golang:1.21-alpine as builder
WORKDIR /src
COPY .. .
RUN go env -w GOPROXY=https://proxy.golang.com.cn,direct && go build -o app main.go

FROM scratch AS runner
COPY --from=builder /src/app .
EXPOSE 8080