[![Tests](https://github.com/mchmarny/disco/actions/workflows/on-push.yaml/badge.svg?branch=main)](https://github.com/mchmarny/disco/actions/workflows/on-push.yaml)
[![Builds](https://github.com/mchmarny/disco/actions/workflows/on-tag.yaml/badge.svg?branch=main)](https://github.com/mchmarny/disco/actions/workflows/on-tag.yaml)
[![Release](https://img.shields.io/github/release/mchmarny/disco.svg)](https://github.com/mchmarny/disco/releases/latest)
[![Go: Version](https://img.shields.io/github/go-mod/go-version/mchmarny/disco.svg)](https://github.com/mchmarny/disco)
[![Go: Report](https://goreportcard.com/badge/github.com/mchmarny/disco)](https://goreportcard.com/report/github.com/mchmarny/disco)
[![License: Apache-2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/mchmarny/disco/blob/main/LICENSE)

# disco 

Utility for containerize workload discovery on GCP.

> Note: this is a personal project not an official Google product.

Features:

* Discover currently deployed container images
  * multiple project and region report with filters
  * deployed image to digest resolution
* Report on vulnerabilities or licenses in these images
  * supports operating system and package-level scans

## Why disco

It's easy to end up with a large number of services across many GCP projects and regions. Google Container Analysis service can scan your Artifact Registry images for vulnerabilities, but currently it only covers base OS, and it's not always easy to know which of these images are actually running in Cloud Run. Cloud Run also supports multiple revisions, each potentially using different version of an image, or even different image all together.

`disco` provides an easy way of `disco`vering which of these container images are currently deployed and are being used in Cloud Run. It extracts the digests (even if the revision is using only a tag (e.g. `v1.2.3`), or that misunderstood `latest`.

## Installation 

If you have Go 1.17 or newer, you can install latest `disco` using:

```shell
go install github.com/mchmarny/disco/cmd/disco@latest
```

You can also download the [latest release](https://github.com/mchmarny/disco/releases/latest) version of `disco` for your operating system/architecture from [here](https://github.com/mchmarny/disco/releases/latest). Put the binary somewhere in your $PATH, and make sure it has that executable bit.

> The official `disco` releases include SBOMs

### Prerequisites 

Since you are interested in `disco`, you probably already have GCP account and project. If not, you learn about creating and managing projects [here](https://cloud.google.com/resource-manager/docs/creating-managing-projects). The other prerequisites include:

#### gcloud

`disco` only uses `gcloud` to provision Application Default Credentials (ADC). You can find instructions on how to install `gcloud` [here](https://cloud.google.com/sdk/docs/install). To provision ADC run and follow the prompts:
  
```shell
gcloud auth application-default login
```

#### APIs

`disco` also depends on a few GCP service APIs. To enable these, run:

```shell
gcloud services enable \
  artifactregistry.googleapis.com \
  cloudresourcemanager.googleapis.com \
  containerregistry.googleapis.com \
  containerscanning.googleapis.com \
  run.googleapis.com 
```

#### Roles

Finally, `disco` is implicitly scoped to only the resources the authenticated user can see. To ensure you can discover resources across multiple projects, make sure you have the following Identity and Access Management (IAM) roles in each project: 

> Learn how to grant multiple IAM roles to a user [here](https://cloud.google.com/iam/docs/granting-changing-revoking-access#multiple-roles)

```shell
roles/artifactregistry.reader
roles/run.viewer
```

Additionally, if you plan to use the Container Analysis option, ensure you also have these roles: 

```shell
roles/containeranalysis.occurrences.viewer
roles/containeranalysis.notes.viewer
```

If you experience any issues, you can see the project level policy using following command:

```shell
gcloud projects get-iam-policy $PROJECT_ID --format=json > policy.json
```

## Supported Runtimes

The general `disco` usage follows this format:

```shell
disco <runtime> <command>
```

> You can use the `--help` flag on any level to get more information about the runtime, commands, of `disco` itself.

The command options available for all the runtimes include:

* `--project` - runs only on specific project (project ID)
* `--format`  - specifies report format: `json`, `yaml`, `raw` (`json` by default)
* `--output`  - saves report to file at this path (stdout by default) 

### Cloud Run 

To see all of the commands available for `run`:

```shell
disco run --help
```

#### Discover container images currently deployed in Cloud Run

```shell
disco run images
```

The `images` command supports all of the generic options listed above, plus: 

* `--digest` - outputs only image digests (default: false). This is helpful when you want to pipe the resulting image digests to another program.

The resulting report in JSON format will look something like this (abbreviated):

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

The resulting report in JSON format will look something like this (abbreviated):

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

The `licenses` command supports all of the generic options listed above, plus: 

* `--cve` - filters report on a specific CVE. This enables quick search if anything currently running is exposed to new CVE.
* `--ca`  - invokes Container Analysis API instead of the local scanner (default: false).     

The resulting report in JSON format will look something like this (abbreviated):

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
