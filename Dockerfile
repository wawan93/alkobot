FROM alpine as alpine
RUN apk add -U --no-cache ca-certificates

FROM golang as builder
RUN mkdir -p /app/alkobot/
WORKDIR /app/alkobot/
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"'

FROM scratch
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/alkobot/alkobot /app/
WORKDIR /app
EXPOSE 8445
ENTRYPOINT ["./alkobot"]
