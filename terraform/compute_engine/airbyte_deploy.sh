echo "***********************"
echo "Start the airbyte server through docker-compose"
echo "***********************"
mkdir airbyte && cd airbyte
wget https://raw.githubusercontent.com/airbytehq/airbyte/master/{.env,docker-compose.yaml}
sudo docker-compose up -d