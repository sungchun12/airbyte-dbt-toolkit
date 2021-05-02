# packer airbyte_gc_image

## Set up your Google Cloud Project and Service Account

```bash
# setup env vars
GOOGLE_CLOUD_PROJECT_ID="dbt-demos-sung"
INSTANCE_NAME="packer-demo"
YOUR_GCP_ZONE="us-central1-a"
YOUR_GCP_NETWORK="default"

# enable all apis in scope
gcloud services enable \
    cloudresourcemanager.googleapis.com \
    iam.googleapis.com \
    compute.googleapis.com

# create a packer service account
gcloud iam service-accounts create packer \
  --project $GOOGLE_CLOUD_PROJECT_ID \
  --description="Packer Service Account" \
  --display-name="Packer Service Account"

# assign compute admin permissions
gcloud projects add-iam-policy-binding $GOOGLE_CLOUD_PROJECT_ID \
    --member=serviceAccount:packer@$GOOGLE_CLOUD_PROJECT_ID.iam.gserviceaccount.com \
    --role=roles/compute.instanceAdmin.v1

# assign the service account user permission
gcloud projects add-iam-policy-binding $GOOGLE_CLOUD_PROJECT_ID \
    --member=serviceAccount:packer@$GOOGLE_CLOUD_PROJECT_ID.iam.gserviceaccount.com \
    --role=roles/iam.serviceAccountUser

# test if you can create a VM with the packer service acccount
gcloud compute instances create $INSTANCE_NAME \
  --project $GOOGLE_CLOUD_PROJECT_ID \
  --image-family ubuntu-2004-lts \
  --image-project ubuntu-os-cloud \
  --network $YOUR_GCP_NETWORK \
  --zone $YOUR_GCP_ZONE \
  --service-account=packer@$GOOGLE_CLOUD_PROJECT_ID.iam.gserviceaccount.com \
  --scopes="https://www.googleapis.com/auth/cloud-platform"

# delete after your test is complete
gcloud compute instances delete $INSTANCE_NAME \
  --project $GOOGLE_CLOUD_PROJECT_ID \
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
    --project $GOOGLE_CLOUD_PROJECT_ID \
    --image="packer-1619640073" \
    --zone $YOUR_GCP_ZONE \
    --service-account=packer@$GOOGLE_CLOUD_PROJECT_ID.iam.gserviceaccount.com

# delete after your test is complete
gcloud compute instances delete $INSTANCE_NAME \
  --project $GOOGLE_CLOUD_PROJECT_ID \
  --zone $YOUR_GCP_ZONE

```

## Manually deploy airybte demo instance

```bash
# setup env vars
GOOGLE_CLOUD_PROJECT_ID="dbt-demos-sung"
INSTANCE_NAME="airbyte-demo"
YOUR_GCP_ZONE="us-central1-a"
YOUR_GCP_NETWORK="default"
MACHINE_TYPE="n1-standard-2" # or e2.medium to save money

# create airybte demo instance
gcloud compute instances create $INSTANCE_NAME \
  --project $GOOGLE_CLOUD_PROJECT_ID \
  --machine-type=$MACHINE_TYPE \
  --image-family debian-10  \
  --image-project debian-cloud \
  --network $YOUR_GCP_NETWORK \
  --zone $YOUR_GCP_ZONE \
  --service-account=packer@$GOOGLE_CLOUD_PROJECT_ID.iam.gserviceaccount.com \
  --scopes="https://www.googleapis.com/auth/cloud-platform"

# In your workstation terminal
# Verify you can see your instance
gcloud --project $GOOGLE_CLOUD_PROJECT_ID compute instances list
[...] # You should see the airbyte instance you just created

# Connect to your instance
# In your workstation terminal
gcloud --project=$GOOGLE_CLOUD_PROJECT_ID beta compute ssh $INSTANCE_NAME

# Install docker
# In your ssh session on the instance terminal
sudo apt-get update
sudo apt-get install -y apt-transport-https ca-certificates curl gnupg2 software-properties-common
curl -fsSL https://download.docker.com/linux/debian/gpg | sudo apt-key add --
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/debian buster stable"
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io
sudo usermod -a -G docker $USER

# Install docker-compose
# In your ssh session on the instance terminal
sudo apt-get -y install wget
sudo wget https://github.com/docker/compose/releases/download/1.26.2/docker-compose-$(uname -s)-$(uname -m) -O /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
docker-compose --version

# Close the ssh connection to ensure the group modification is taken into account
# In your ssh session on the instance terminal
logout

# In your workstation terminal
gcloud --project=$GOOGLE_CLOUD_PROJECT_ID beta compute ssh $INSTANCE_NAME

# In your ssh session on the instance terminal
mkdir airbyte && cd airbyte
wget https://raw.githubusercontent.com/airbytehq/airbyte/master/{.env,docker-compose.yaml}
sudo docker-compose up -d

# Run the airbyte webserver locally on your workstation through an ssh tunnel
gcloud --project=$GOOGLE_CLOUD_PROJECT_ID beta compute ssh $INSTANCE_NAME \
    --zone $YOUR_GCP_ZONE \
    -- -L 8000:localhost:8000 -L 8001:localhost:8001 -N -f

# In your local workstation browser
http://localhost:8000


# delete after your test is complete
gcloud compute instances delete $INSTANCE_NAME \
  --project $GOOGLE_CLOUD_PROJECT_ID \
  --zone $YOUR_GCP_ZONE

```

## Sources

[Tutorial](https://www.packer.io/docs/builders/googlecompute)
[airbyte Google Compute Engine Deployment](https://docs.airbyte.io/deploying-airbyte/on-gcp-compute-engine)
