#!/bin/bash

echo "***********************"
echo "Install docker"
echo "***********************"
sudo apt-get update
sudo apt-get install -y apt-transport-https ca-certificates curl gnupg2 software-properties-common
curl -fsSL https://download.docker.com/linux/debian/gpg | sudo apt-key add --
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/debian buster stable"
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io
sudo usermod -a -G docker $USER

echo "***********************"
echo "Install docker-compose"
echo "***********************"
sudo apt-get -y install wget
sudo wget https://github.com/docker/compose/releases/download/1.26.2/docker-compose-$(uname -s)-$(uname -m) -O /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
docker-compose --version

echo "***********************"
echo "Start the airbyte server through docker-compose"
echo "This will NOT persist within the packer image"
echo "This will verify the airbyte server can start"
echo "***********************"
mkdir airbyte && cd airbyte
wget https://raw.githubusercontent.com/airbytehq/airbyte/master/{.env,docker-compose.yaml}
sudo docker-compose up -d