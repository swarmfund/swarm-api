FROM golang:1.9

WORKDIR /go/src/gitlab.com/swarmfund/api
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /binary -v gitlab.com/swarmfund/api/cmd/api

FROM alpine:latest
COPY --from=0 /binary .
RUN apk --no-cache add ca-certificates
ENTRYPOINT ["./binary", "--config", "/config.yaml"]
