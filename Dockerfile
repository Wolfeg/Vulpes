FROM golang:1.21-alpine AS builder

WORKDIR /usr/src/vulpes

COPY . .
RUN go build -v

FROM alpine:3

RUN apk --no-cache add ca-certificates

WORKDIR /usr/local/bin

COPY --from=builder /usr/src/vulpes/vulpes .
CMD ["./vulpes"]