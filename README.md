# disco 

Utility for containerize workload discovery.

Features:

* Discover currently deployed container images across multiple projects and regions
* Resolve deployed images to their actual digests
* Report on operating system and package-level vulnerabilities or licenses in these images

> Note: this is a personal project not an official Google product.

* [Why disco](#why-disco)
* [Install](#install)
* [Supported Runtimes](#supported-runtimes)
  * [Cloud Run](#cloud-run)
  * [GKE](#gke)


## Why disco

It's easy to end up with a large number of services across many GCP projects and regions. Google Container Analysis service can scan your Artifact Registry images for vulnerabilities, but currently it only covers base OS, and it's not always easy to know which of these images are actually running in Cloud Run. Cloud Run also supports multiple revisions, each potentially using different version of an image, or even different image all together.

`disco` provides an easy way of `disco`vering which of these container images are currently deployed and are being used in Cloud Run. It extracts the digests (even if the revision is using only a tag (e.g. `v1.2.3`), or that misunderstood `latest`.

## Install 

If you have Go 1.17+, you can install `disco` directly using this command:

```shell
go install github.com/mchmarny/disco/cmd/disco@latest
```

You can also download the [latest release](https://github.com/mchmarny/disco/releases/latest) version of `disco` for your operating system/architecture from [here](https://github.com/mchmarny/disco/releases/latest). Put the binary somewhere in your $PATH, and make sure it has that executable bit.

> The official `disco` releases include SBOMs

### Prerequisites 

Since you are interested in `disco`, you probably already have GCP account and project. Here are some of the other prerequisites:

#### gcloud

To invoke GCP APIs, `disco` uses `gcloud`. You can find instructions on how to install it [here](https://cloud.google.com/sdk/docs/install). Once installed, you will need to provisioned Application Default Credentials (ADC):
  
```shell
gcloud auth application-default login
```

#### Service APIs

`disco` also depends on a few GCP service APIs to be enabled on each project you want to access:

```shell
gcloud services enable \
  artifactregistry.googleapis.com \
  cloudresourcemanager.googleapis.com \
  containerregistry.googleapis.com \
  containerscanning.googleapis.com \
  run.googleapis.com 
```

#### Roles

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



## Supported Runtimes

The general usage looks like this:

```shell
disco <runtime> <command>
```

Options:

* `--project` - runs only on specific project (project ID)
* `--output`  - saves report to file at this path (stdout by default) 
* `--format`  - specifies report format: `json`, `yaml`, `raw` (`json` by default)
* `--help`    - shows help 

> Currently, `disco` implements only Cloud Run as target runtime.

### Cloud Run 

#### Discover container images currently deployed in Cloud Run

```shell
disco run images
```

All the generic options listed above, plus: 

* `--digest`  - outputs only image digests (default: false)

> The `--digest` flag is helpful when you want to pipe the resulting list to another program.

The resulting JSON formatted report looks something like this (abbreviated):

```json
[
  {
    "region": "us-central1",
    "project": "cloudy-demos",
    "service": "hello",
    "image": "https://us-docker.pkg.dev/cloudrun/container/hello@sha256:2e70803dbc92a7bffcee3af54b5d264b23a6096f304f00d63b7d1e177e40986c"
  },
  ...
]
```

#### Discover licenses used in container images currently deployed in Cloud Run

```shell
disco run licenses
```

The resulting JSON formatted report looks something like this (abbreviated):

```json
[
  {
    "image": "us-docker.pkg.dev/cloudrun/container/hello@sha256:2e70803dbc92a7bffcee3af54b5d264b23a6096f304f00d63b7d1e177e40986c",
    "licenses": [
      {
        "name": "GPL-2.0",
        "source": "alpine-baselayout"
      },
      {
        "name": "MPL-2.0",
        "source": "ca-certificates"
      },
      {
        "name": "MIT",
        "source": "ca-certificates"
      },
      ...
    ]
  },
  ...
]
```


#### Discover potential vulnerabilities in container images currently deployed in Cloud Run

```shell
disco run licenses
```

All the generic options listed above, plus: 

* `--cve` - filters report only to a specific CVE
* `--ca`  - invokes Container Analysis API in stead of the local scanner (default: false)

> The `--cve` is a quick way to finding out if anything currently running is exposed to new CVE.                       

The resulting JSON formatted report looks something like this (abbreviated):

```json
[
  {
    "image": "gcr.io/cloudy-demos/hello-broken@sha256:0900c08e7d40f9485c8497c035de07391ba3c274a1035f504f8602531b2314e6",
    "vulnerabilities": [
      {
        "source": "CVE-2022-3715",
        "severity": "LOW",
        "package": "bash",
        "version": "5.1-6ubuntu1",
        "title": "bash: a heap-buffer-overflow in valid_parameter_transform",
        "description": "A flaw was found in the bash package, where a heap-buffer overflow can occur in valid_parameter_transform. This issue may lead to memory problems.",
        "url": "https://avd.aquasec.com/nvd/cve-2022-3715",
        "updated": "2022-12-23T16:52:00Z"
      },
      ...
    ]
  },
  ...
]
```

### GKE

> Not yet implemented.


## Disclaimer

This is my personal project and it does not represent my employer. While I do my best to ensure that everything works, I take no responsibility for issues caused by this code.
