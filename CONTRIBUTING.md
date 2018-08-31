# Contributing to API

## Checklist

[ ] start a `feature` or `hotfix` branch
[ ] implement your stuff
[ ] don't forget `go generate ./...` if new migrations added or interfaces were updated
[ ] ensure `go test ./...` passes
[ ] review your changes one more time
[ ] send merge request to master and assign someone
[ ] poke assignee time to time to get feedback and merge
TODO Explain design options
TODO Ways to create user
TODO Explain TFA
## Flow

API uses go-bindata package for generating migrations
It can be found here: [https://github.com/jteeuwen/go-bindata
