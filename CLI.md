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
  containerscanning.googleapis.com \
  run.googleapis.com 
```

### Roles

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

### Supported Vulnerability Scanners 

`disco` shells out the vulnerability scans to one of the supported OSS scanners: 

* [trivy (disco default)](https://aquasecurity.github.io/trivy/v0.35/getting-started/installation/)
* [grype](https://github.com/anchore/grype#installation)


## CLI Usage

```shell
disco [runtime] [command] [arguments...]
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

* [Images](#images)
* [Licenses](#licenses)
* [Vulnerabilities](#vulnerabilities)

#### Images

To discover container images currently deployed in Cloud Run:

```shell
disco run img
```

The `images` or `img` command supports all of the generic options listed above, plus: 

* `--uri` - outputs only image uri (default: false). This is helpful when you want to pipe the resulting images to another program.

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

#### Licenses

To discover licenses used in container images currently deployed in Cloud Run.

```shell
disco run lic
```

The `licenses` or `lic` command supports all of the generic options listed above, plus: 

* `--source` - path to image list file to use as source. This allows you to use the previously generated list of images (`disco run img --uri -o images.txt`), instead of running through potentially lengthy discovery. 
* `--image` - specific image URI to scan. Note: `source` and `image` are mutually exclusive.

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

#### Vulnerabilities

To discover potential vulnerabilities in container images currently deployed in Cloud Run.

```shell
disco run vul
```

The `vul` or `vulnerabilities` command supports all of the generic options listed above, plus: 

* `--source` - path to image list file to use as source. This allows you to use the previously generated list of images (e.g. `disco run img --uri -o images.txt`). If not provided, `disco` will discover images first. 
* `--image` - specific image URI to scan. Note: `source` and `image` are mutually exclusive.
* `--min-severity` - minimum severity of vulnerability to include in report (e.g. low, medium, high, critical, default: all).
* `--cve` - filters report on a specific CVE. This enables quick search if anything currently running is exposed to a new CVE.

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

## Disclaimer

This is my personal project and it does not represent my employer. While I do my best to ensure that everything works, I take no responsibility for issues caused by this code.
