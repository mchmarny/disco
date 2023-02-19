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

> Make sure one of these is installed before running disco CLI

* [trivy (disco default)](https://aquasecurity.github.io/trivy/v0.35/getting-started/installation/)
* [grype](https://github.com/anchore/grype#installation)


## Disclaimer

This is my personal project and it does not represent my employer. While I do my best to ensure that everything works, I take no responsibility for issues caused by this code.
