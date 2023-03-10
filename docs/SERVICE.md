# Disco Service

![](../docs/img/diagram.png)

* [Deployment](#deployment)
* [Usage](#usage)

## Deployment
  
To deploy the prebuilt version of `disco`, first clone this repo:

```shell
git clone git@github.com:mchmarny/disco.git
```

Then navigate to the `deploy` directory inside of that cloned repo:

```shell
cd disco/deploy
```

Next, authenticate to GCP:

```shell
gcloud auth application-default login
```

Initialize Terraform: 

```shell
terraform init
```

> Note, this flow uses the default, local terraform state. Make sure you do not check the state files into your source control (see `.gitignore`), or consider using persistent state provider like GCS.

When done, apply the Terraform configuration:

```shell
terraform apply
```

When promoted, provide requested variables:

* `project_id` is the GCP project ID (not the name)
* `location` is GCP region to deploy to
* `git_repo` qualified name of the newly cloned repo (e.g. `username/disco`)

When completed, this will output the configured resource information. 

### Cross-project Discovery

`disco` service will automatically discover services in any project where it has been granted `viewer` role. The name of the service account used by `disco` to discover will be listed under `RUN_SERVICE_ACCOUNT` after service deployed.

For example, assuming the ID of the project where you deployed `disco` is `$DISCO_PROJECT_ID` and the project in which you want it to discover services is `$TARGET_PROJECT_ID`, the command granting the necessary role would look like this:


```shell
gcloud projects add-iam-policy-binding $TARGET_PROJECT_ID \
    --member "serviceAccount:disco-run-sa@${DISCO_PROJECT_ID}".iam.gserviceaccount.com" \
    --role "roles/viewer"
```

> Note, you can make `disco` discover across as many projects as you wish. 

### Test Deployment

To test the deployed `disco` service:

```shell
SERVICE_URL=$(gcloud run services describe disco \
    --region $REGION --format="value(status.url)")

curl -sS -H "Authorization: Bearer $(gcloud auth print-identity-token)" \
     -H "Content-Type: application/json" \
     -H "X-Goog-User-Project: ${PROJECT_ID}" \
     "${SERVICE_URL}/disco"
```

A correctly deployed service should return: 

```json
{ "status": "OK", "message": "Done" }
```

## Disclaimer

This is my personal project and it does not represent my employer. While I do my best to ensure that everything works, I take no responsibility for issues caused by this code.
