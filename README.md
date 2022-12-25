# disco 

Helper utility for containerize workload discovery.

Features:

* Discover currently deployed container images across multiple projects and regions
* Resolve deployed images to their digests
* Use discovered digests with open source SBOM generation or vulnerability scanner tools of your choice to discover:
  * OS and packages
  * Vulnerabilities 
  * Licenses 

> Note: this is a personal project not an official Google product.

* [Supported Runtimes](#supported-runtimes)
  * [Cloud Run](#cloud-run)
* [Install](#install)
* [Usage](#usage)

## Supported Runtimes

### Cloud Run 

Google Cloud Run is a great runtime for many use-cases. It's easy to end up with a large number of services across many GCP projects and regions. Google Container Analysis service can scan your Artifact Registry images for vulnerabilities, but currently it only covers base OS, and it's not always easy to know which of these images are actually currently deployed. Cloud Run also supports multiple revisions, each potentially using different version of an image, or even different image all together.

`disco` provides an easy way of `disco`vere which of these container images are currently deployed and are being used in Cloud Run. It extracts the digests (even if the revision is using only a tag (e.g. `v1.2.3`), or that misunderstood `latest`.

#### Prerequisites 

Since you are interested in `disco`, you probably already have GCP account and project. Here are some of the other prerequisites:

###### gcloud

To invoke GCP APIs, `disco` uses `gcloud`. You can find instructions on how to install it [here](https://cloud.google.com/sdk/docs/install). Once installed, you will need to provisioned Application Default Credentials (ADC):
  
```shell
gcloud auth application-default login
```

##### Service APIs

`disco` also depends on a few GCP service APIs to be enabled on each project you want to access:

```shell
gcloud services enable \
  artifactregistry.googleapis.com \
  cloudresourcemanager.googleapis.com \
  containerregistry.googleapis.com \
  containerscanning.googleapis.com \
  run.googleapis.com 
```

##### Roles

Finally, make sure you have the required Identity and Access Management (IAM) roles: 

```shell
roles/artifactregistry.reader
roles/containeranalysis.occurrences.viewer
roles/containeranalysis.notes.viewer
roles/run.viewer
```

You can check if you already have these roles for a given project like this:

```shell
gcloud projects get-iam-policy $PROJECT_ID --format=json > policy.json
```

> Learn how to grant multiple IAM roles to a user [here](https://cloud.google.com/iam/docs/granting-changing-revoking-access#multiple-roles)

## Install 

If you have Go 1.17+, you can install `disco` directly using this command:

```shell
go install github.com/mchmarny/disco/cmd/disco@latest
```

You can also download the [latest release](https://github.com/mchmarny/disco/releases/latest) version of `disco` for your operating system/architecture from [here](https://github.com/mchmarny/disco/releases/latest). Put the binary somewhere in your $PATH, and make sure it has that executable bit.

## Usage



### trivy 

```shell
trivy image us-docker.pkg.dev/cloudrun/container/hello@sha256:2e70803dbc92a7bffcee3af54b5d264b23a6096f304f00d63b7d1e177e40986c \
    --security-checks license \
    --timeout 10m \
    --format json \
    --output test.json
```

## Disclaimer

This is my personal project and it does not represent my employer. While I do my best to ensure that everything works, I take no responsibility for issues caused by this code.
