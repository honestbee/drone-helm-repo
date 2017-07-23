# drone-helm-repo

[![Coverage Status](https://coveralls.io/repos/github/honestbee/drone-helm-repo/badge.svg)](https://coveralls.io/github/honestbee/drone-helm-repo)
[![Docker Repository on Quay](https://quay.io/repository/honestbee/drone-helm-repo/status "Docker Repository on Quay")](https://quay.io/repository/honestbee/drone-helm-repo)

Drone plugin to package and upload Helm charts to selected storage services

Supported static file hosts:

- AWS S3

Used github.com/silverstripeltd/s3sync as a reference

## Testing:

```

export pkg=go list ./...
go get -v -t $pkg
go test $pkg -v -cover
```

## Building:

```
go build -a -tags netgo -o bin/drone-helm-repo
```

## Known issues:

- Updated index does not support merging with existing index, this means previous Chart versions are not kept in index.yaml
  should add a feature to pull an existing index.yaml first and merge with it