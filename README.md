# helm-repo-sync

Drone plugin to sync Repository to static file host

Supported static file hosts:

- AWS S3

used github.com/silverstripeltd/s3sync as a base

## Testing:

```
export pkg=go list ./...
go test $pkg -v -cover
```

## Building:

```
go build -a -tags netgo -o bin/drone-helm-repo
```

## Known issues:

- Index is not merged, this means previous versions are not kept in index.yaml
  should add a feature to pull an existing index first and merge with it