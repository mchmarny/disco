# Disco CLI

* [Installation](#cli-installation)
* [Usage](#cli-usage)

## CLI Installation 

You can install `disco` CLI using one of the following ways:

> If you are already using `gcloud` to manage GCP resources changes are you already have all the prerequisites. Otherwise, review the [prerequisites](#prerequisites) section.

* [Homebrew](#homebrew)
* [RHEL/CentOS](#rhelcentos)
* [Debian/Ubuntu](#debianubuntu)
* [Go](#go)
* [Binary](#binary)

See the [release section](https://github.com/mchmarny/disco/releases/latest) for `disco` checksums and SBOMs.

> See the list of supported open source [vulnerability scanner](#supported-vulnerability-scanners).

## Homebrew

On Mac or Linux, you can install `disco` with [Homebrew](https://brew.sh/):

```shell
brew tap mchmarny/disco
brew install disco
```

New release will be automatically picked up when you run `brew upgrade`

## RHEL/CentOS

```shell
rpm -ivh https://github.com/mchmarny/disco/releases/download/v$VERSION/disco-$VERSION_Linux-amd64.rpm
```

## Debian/Ubuntu

```shell
wget https://github.com/aquasecurity/disco/releases/download/v$VERSION/disco-$VERSION_Linux-amd64.deb
sudo dpkg -i disco-$VERSION_Linux-64bit.deb
```

## Go

If you have Go 1.17 or newer, you can install latest `disco` using:

```shell
go install github.com/mchmarny/disco/cmd/disco@latest
```

## Binary 

You can also download the [latest release](https://github.com/mchmarny/disco/releases/latest) version of `disco` for your operating system/architecture from [here](https://github.com/mchmarny/disco/releases/latest). Put the binary somewhere in your $PATH, and make sure it has that executable bit.

> The official `disco` releases include SBOMs

## Prerequisites 

Since you are interested in `disco`, you probably already have GCP account and project. If not, you learn about creating and managing projects [here](https://cloud.google.com/resource-manager/docs/creating-managing-projects). The other prerequisites include:

### gcloud

`disco` only uses `gcloud` to provision Application Default Credentials (ADC). You can find instructions on how to install `gcloud` [here](https://cloud.google.com/sdk/docs/install). To provision ADC run and follow the prompts:
  
```shell
gcloud auth application-default login
```

### APIs

`disco` also depends on a few GCP service APIs. To enable these, run:

```shell
gcloud services enable \
  artifactregistry.googleapis.com \
  cloudresourcemanager.googleapis.com \
  containerregistry.googleapis.com \
  run.googleapis.com 
```

### Roles

Finally, `disco` is implicitly scoped to only the resources the authenticated user can see. To ensure you can discover resources across multiple projects, make sure you have the following Identity and Access Management (IAM) roles in each project: 

> Learn how to grant multiple IAM roles to a user [here](https://cloud.google.com/iam/docs/granting-changing-revoking-access#multiple-roles)

```shell
roles/artifactregistry.reader
roles/run.viewer
```

If you experience any issues, you can see the project level policy using following command:

```shell
gcloud projects get-iam-policy $PROJECT_ID --format=json > policy.json
```

### Supported Vulnerability Scanners 

`disco` shells out the vulnerability scans to one of the supported OSS scanners: 

* [trivy (disco default)](https://aquasecurity.github.io/trivy/v0.35/getting-started/installation/)
* [grype](https://github.com/anchore/grype#installation)

## CLI Usage

```shell
disco command [command options] [arguments...]
```

> You can use the `--help` flag on any level to get more information about the runtime, commands, of `disco` itself.

### Images

Discover deployed images from specific runtime. To see all of the commands available for `img` (or `image`):

```shell
disco image --help
```

To discover container images currently deployed in all of the supported runtimes:

```shell
disco img
```

Options:

* `--output` - path where to save the output (stdout by default) 
* `--format` - output format (`json` or `yaml`, `json` is the default)
* `--uri` - outputs only image uri (default: false). This is helpful when you want to pipe the resulting images to another program.
* `--project` - scope discovery to a single project using project ID


The resulting report in JSON format will look something like this (abbreviated):

```json
{
  "meta": {
    "kind": "image",
    "version": "v0.3.19-next",
    "created": "2022-12-28T21:20:15Z",
    "count": 7
  },
  "items": [
    {
      "uri": "us-docker.pkg.dev/cloudrun/container/hello@sha256:2e70803dbc92a7bffcee3af54b5d264b23a6096f304f00d63b7d1e177e40986c",
      "context": {
        "container-name": "hello-1",
        "location-id": "us-central1",
        "location-name": "Iowa",
        "project-id": "cloudy-demos",
        "project-number": "799736955886",
        "service-id": "projects/cloudy-demos/locations/us-central1/services/hello",
        "service-name": "hello",
        "service-revision": "projects/cloudy-demos/locations/us-central1/services/hello/revisions/hello-00001-taj"
      }
    },
      ...
  ]
}
```

### Vulnerabilities

Discover potential vulnerabilities in container images. To see all of the commands available for `vul` (or `vulnerability`):

```shell
disco vulnerability --help
```

Options: 

* `--file` - image list input file path to serve as a source (instead of discovery) (e.g. `disco img --uri --output images.txt`).
* `--image` - specific image URI to scan. Note: `source` and `image` are mutually exclusive.
* `--output`  - saves report to file at this path (stdout by default) 
* `--format`  - report format: `json` or `yaml` (`json` is default)
* `--project` - during discovery, runs only on specific project (project ID)
* `--min-severity` - minimum severity of vulnerability to include in report (e.g. low, medium, high, critical, default: all)
* `--cve` - filter results on a specific CVE ID (e.g. `CVE-2020-22046`)
* `--target` - target data store to save the results to (e.g. `bq://my-project.some-dataset.table-name`)

> Using the `cve` filter you can quickly check if any of the currently deployed images have a vulnerability. 

The resulting report in JSON format will look something like this (abbreviated):

```json
{
  "meta": {
    "kind": "vulnerability",
    "version": "v0.3.19-next",
    "created": "2022-12-28T21:32:34Z",
    "count": 5
  },
  "items": [
    {
      "image": "gcr.io/cloudy-demos/hello-broken@sha256:0900c08e7d40f9485c8497c035de07391ba3c274a1035f504f8602531b2314e6",
      "vulnerabilities": [
        {
          "source": "CVE-2021-28165",
          "severity": "HIGH",
          "package": "org.eclipse.jetty:jetty-util",
          "version": "9.4.31.v20200723",
          "title": "jetty: Resource exhaustion when receiving an invalid large TLS frame",
          "description": "In Eclipse Jetty 7.2.2 to 9.4.38, 10.0.0.alpha0 to 10.0.1, and 11.0.0.alpha0 to 11.0.1, CPU usage can reach 100% upon receiving a large invalid TLS frame.",
          "url": "https://avd.aquasec.com/nvd/cve-2021-28165",
          "updated": "2022-07-29T17:05:00Z"
        },
        ...
      ]
    },
    ...
  ]
}
```

### Licenses

Discover licenses for OS and packages used in container images. To see all of the commands available for `lic` or `license`:

```shell
disco license --help
```

Options: 

* `--file` - image list input file path to serve as a source (instead of discovery) (e.g. `disco img --uri --output images.txt`).
* `--image` - specific image URI to scan. Note: `source` and `image` are mutually exclusive.
* `--output`  - saves report to file at this path (stdout by default)  
* `--format`  - report format: `json` or `yaml` (`json` is default)
* `--project` - during discovery, runs only on specific project (project ID)
* `--type` - license type filter (supports prefix: e.g. `apache`, `bsd`, `mit`, etc.).
* `--target` - target data store to save the results to (e.g. `bq://my-project.some-dataset.table-name`)

> Using the `type` you can quickly check if any of your currently deployed images are using specific license.

The resulting report in JSON format will look something like this (abbreviated):

```json
{
  "meta": {
    "kind": "license",
    "version": "v0.3.19-next",
    "created": "2022-12-28T21:23:20Z",
    "count": 7
  },
  "items": [
    {
      "image": "us-docker.pkg.dev/cloudrun/container/hello@sha256:2e70803dbc92a7bffcee3af54b5d264b23a6096f304f00d63b7d1e177e40986c",
      "licenses": [
        {
          "name": "GPL-2.0",
          "source": "alpine-baselayout-data"
        },
        {
          "name": "MIT",
          "source": "alpine-keys"
        },
        ...
      ]
    },
    ...
  ]
}
```

### Packages

Discover packages used in container images. To see all of the commands available for `pkg` (or `packages`):

```shell
disco packages --help
```

Options: 

* `--file` - image list input file path to serve as a source (instead of discovery) (e.g. `disco img --uri --output images.txt`).
* `--image` - specific image URI to scan. Note: `source` and `image` are mutually exclusive.
* `--output`  - saves report to file at this path (stdout by default)  
* `--format`  - report format: `json` or `yaml` (`json` is default)
* `--project` - during discovery, runs only on specific project (project ID)
* `--target` - target data store to save the results to (e.g. `bq://my-project.some-dataset.table-name`)

> Using the `type` you can quickly check if any of your currently deployed images are using specific license.

The resulting report in JSON format will look something like this (abbreviated):

```json
{
  "meta": {
    "kind": "package",
    "version": "v0.9.4",
    "created": "2023-01-08T00:37:26Z",
    "count": 8
  },
  "items": [
    {
      "image": "us-central1-docker.pkg.dev/cloudy-labz/gcf-artifacts/test--go119@sha256:80be8e3c174f20aebc555d932c8ba3a8452e1eb46af1923200ac39f4ccd30016",
      "packages": [
        {
          "package": "gobinary",
          "version": "",
          "format": "SPDX-2.2",
          "provider": "trivy",
          "source": "layers/google.go.build/bin/main"
        },
        ...
      ],
    ...
    }
  ]
}
```


## Disclaimer

This is my personal project and it does not represent my employer. While I do my best to ensure that everything works, I take no responsibility for issues caused by this code.
