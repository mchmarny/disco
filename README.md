# vctl 

Helper utility for containerize workload security. 



## prerequisites 



`gcloud auth application-default login`

project must have enabled:

```shell
gcloud services enable \
  artifactregistry.googleapis.com \
  cloudresourcemanager.googleapis.com \
  containerregistry.googleapis.com \
  containerscanning.googleapis.com \
  run.googleapis.com 
```

roles:

```shell
roles/artifactregistry.admin
roles/containeranalysis.occurrences.viewer
```

## why 

https://cloud.google.com/sdk/gcloud/reference/artifacts/docker/images/describe

## trivy 

```shell
trivy image us-docker.pkg.dev/cloudrun/container/hello@sha256:2e70803dbc92a7bffcee3af54b5d264b23a6096f304f00d63b7d1e177e40986c \
    --security-checks license \
    --timeout 10m \
    --format json \
    --output test.json
```