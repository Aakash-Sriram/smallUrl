FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/url-shortener ./

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /root/
COPY --from=builder /usr/local/bin/url-shortener .
EXPOSE 9808
ENV PORT=9808
CMD ["./url-shortener"]
