# Helm Repo plugin for Drone

[![Coverage Status](https://coveralls.io/repos/github/honestbee/drone-helm-repo/badge.svg)](https://coveralls.io/github/honestbee/drone-helm-repo)
[![Docker Repository on Quay](https://quay.io/repository/honestbee/drone-helm-repo/status "Docker Repository on Quay")](https://quay.io/repository/honestbee/drone-helm-repo)

Drone plugin to package and upload Helm charts to selected storage services

**Note**: This plugin does not merge index.yaml and only keeps the latest chart version on the Repository index. To fix this, we wrote a plugin that pushes to [kubernetes-helm/ChartMuseum](https://github.com/kubernetes-helm/chartmuseum) instead (which supports S3/GCS/... )

See [honestbee/drone-chartmuseum](https://github.com/honestbee/drone-chartmuseum)

## Description

When managing Charts for your organisation, you may either choose to put Chart definitions within each project or centralised in a `helm-charts` repository.

Keeping Charts in a central repository allows a more flexible definition of how each project is deloyed and a full decoupling of the projects (similar to the central public-charts repo)

This plugin provides a Drone build step to package and update the Helm Repository for centralised charts repository.

Read our [blog post](http://tech.honestbee.com/articles/devops/2017-07/drone-helm-repository) for more details.

## Usage

Secrets:

```bash
drone secret add --image=quay.io/honestbee/drone-helm-repo \
  your-user/your-repo AWS_ACCESS_KEY_ID AKIA...

drone secret add --image=quay.io/honestbee/drone-helm-repo \
  your-user/your-repo AWS_SECRET_ACCESS_KEY ...

drone secret add --image=quay.io/honestbee/drone-helm-repo \
  your-user/your-repo AWS_REGION ap-southeast-1
```

Drone Usage:

```YAML
pipeline:
  update_helm_repo:
    image: quay.io/honestbee/drone-helm-repo
    exclude: .git
    repo_url: http://helm-charts.example.com
    storage_url: s3://helm-charts.example.com
    aws_region: ap-southeast-1
    when:
      branch: [master]

```

CLI Options:

```bash
   --source-dir PATH                PATH to recursively search for Charts (default: ".") [$PLUGIN_SOURCE_DIR, $SOURCE_DIR]
   --exclude LIST                   LIST of excluded directories in source-dir to exclude from Chart search [$PLUGIN_EXCLUDE, $EXCLUDE]
   --repo-url BASE_URL              BASE_URL for the helm repository [$PLUGIN_REPO_URL, $REPO_URL]
   --storage-url URL                URL of the container to store charts to (i.e s3://my-bucket/prefix) [$PLUGIN_STORAGE_URL, $STORAGE_URL]
   --aws-access-key AWS_ACCESS_KEY  AWS Access Key AWS_ACCESS_KEY [$AWS_ACCESS_KEY_ID, $AWS_ACCESS_KEY]
   --aws-secret-key AWS_SECRET_KEY  AWS Secret Key AWS_SECRET_KEY [$AWS_SECRET_ACCESS_KEY, $AWS_SECRET_KEY]
   --aws-region AWS_REGION          AWS Region AWS_REGION [$AWS_REGION]
   --debug                          show debug logs [$PLUGIN_DEBUG, $DEBUG]
   --help, -h                       show help
   --version, -v                    print the version
```

**Note**: Debug option will print all env vars, be careful when enabling this for a public build log as it may expose your secrets

```bash
docker run --rm \
  - EXCLUDE=".git" \
  - REPO_URL="http://helm-charts.example.com" \
  - STORAGE_URL="s3://helm-charts.example.com" \
  - AWS_REGION="ap-southeast-1" \
  - AWS_ACCESS_KEY="AKI..." \
  - AWS_SECRET_KEY="qlr5B..." \
  quay.io/honestbee/drone-helm-repo

```

References:

- [ipedrazas/drone-helm](github.com/ipedrazas/drone-helm)
- [silverstripeltd/s3sync](github.com/silverstripeltd/s3sync)

## Testing:

To test all packages (ensure to install test dependencies)

```
export pkg=go list ./...
go get -v -t $pkg
go test $pkg -v -cover
```

`helm` binary will be called by tests to generate files and should be available as well

## Building:

Generate Binary

```
go build -o bin/drone-helm-repo
```

## Known issues:

- Updated index does not support merging with existing index, this means previous Chart versions are not kept in index.yaml
  
  should add a feature to pull an existing index.yaml first and merge with it
- Add support for GCS (see [lovoo/drone-gcloud-helm](https://github.com/lovoo/drone-gcloud-helm))
