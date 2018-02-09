FROM golang:1.9

WORKDIR /go/src/gitlab.com/swarmfund/api
COPY . .
ENTRYPOINT ["go", "test"]