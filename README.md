# vctl

WIP

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