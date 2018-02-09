FROM golang:1.9

WORKDIR /go/src/gitlab.com/swarmfund/api
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /binary -v gitlab.com/swarmfund/api/cmd/api

FROM ubuntu:latest
COPY --from=0 /binary .
RUN apt-get update && apt-get install -y \
ca-certificates
ENTRYPOINT ["./binary", "--config", "/config.yaml"]
