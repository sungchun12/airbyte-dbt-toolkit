# packer airbyte_gc_image

## Set up your Google Cloud Project and Service Account

```bash
# setup env vars
YOUR_GCP_PROJECT="dbt-demos-sung"
INSTANCE_NAME="packer-demo"
YOUR_GCP_ZONE="us-east1-b"
YOUR_GCP_NETWORK="default"

# enable all apis in scope
gcloud services enable \
    cloudresourcemanager.googleapis.com \
    iam.googleapis.com \
    compute.googleapis.com

# create a packer service account
gcloud iam service-accounts create packer \
  --project $YOUR_GCP_PROJECT \
  --description="Packer Service Account" \
  --display-name="Packer Service Account"

# assign compute admin permissions
gcloud projects add-iam-policy-binding $YOUR_GCP_PROJECT \
    --member=serviceAccount:packer@$YOUR_GCP_PROJECT.iam.gserviceaccount.com \
    --role=roles/compute.instanceAdmin.v1

# assign the service account user permission
gcloud projects add-iam-policy-binding $YOUR_GCP_PROJECT \
    --member=serviceAccount:packer@$YOUR_GCP_PROJECT.iam.gserviceaccount.com \
    --role=roles/iam.serviceAccountUser

# test if you can create a VM with the packer service acccount
gcloud compute instances create $INSTANCE_NAME \
  --project $YOUR_GCP_PROJECT \
  --image-family ubuntu-2004-lts \
  --image-project ubuntu-os-cloud \
  --network $YOUR_GCP_NETWORK \
  --zone $YOUR_GCP_ZONE \
  --service-account=packer@$YOUR_GCP_PROJECT.iam.gserviceaccount.com \
  --scopes="https://www.googleapis.com/auth/cloud-platform"

# delete after your test is complete
gcloud compute instances delete $INSTANCE_NAME \
  --project $YOUR_GCP_PROJECT \
  --zone $YOUR_GCP_ZONE

```

## Build Your Packer Image

```bash
# validate the packer template has no errors
packer validate airbyte_gce_image.pkr.hcl

# export Google Cloud credentials
export GOOGLE_APPLICATION_CREDENTIALS="../service_account.json"

# build the image and have Google Cloud store it under Compute Engine "Images"
packer build airbyte_gce_image.pkr.hcl

# test launching the image with the packer image id
gcloud compute instances create $INSTANCE_NAME \
    --project $YOUR_GCP_PROJECT \
    --image="packer-1619555494" \
    --zone $YOUR_GCP_ZONE \
    --service-account=packer@$YOUR_GCP_PROJECT.iam.gserviceaccount.com

# delete after your test is complete
gcloud compute instances delete $INSTANCE_NAME \
  --project $YOUR_GCP_PROJECT \
  --zone $YOUR_GCP_ZONE

```

## Sources

[Tutorial](https://www.packer.io/docs/builders/googlecompute)
