#!/bin/bash
# Bootstrap a container on ?ec2? w/ auto updates as new images are pushed to ecr

function read_sensitive() {
  local prompt=$1
  local var_name=$2
  read -rsp "$prompt" "$var_name"
  echo
  export "$var_name"
}

read_sensitive "Enter AWS_ACCESS_KEY_ID: " AWS_ACCESS_KEY_ID
read_sensitive "Enter AWS_SECRET_ACCESS_KEY: " AWS_SECRET_ACCESS_KEY
read_sensitive "Enter WATCHTOWER_KEY: " WATCHTOWER_KEY


export IMAGE_REPO_NAME="my_images"
export IMAGE_TAG="scuffed_metar_image"
export HELPER_IMAGE_TAG="aws-ecr-dock-cred-helper"
export AWS_DEFAULT_REGION=us-east-1

# Install aws cli
sudo apt-get install unzip
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip -qq awscliv2.zip
sudo ./aws/install

# grab account id (used create ecr address)
AWS_ACCOUNT_ID=$(aws sts get-caller-identity | awk -F '"' '/"Account":/ {print $4}')
export AWS_ACCOUNT_ID

if [[ -z $AWS_ACCOUNT_ID ]]; then
  echo "AWS_ACCOUNT_ID is empty exiting."
  exit 1
fi


export ECR_ADDRESS="$AWS_ACCOUNT_ID.dkr.ecr.$AWS_DEFAULT_REGION.amazonaws.com"


# The directory .docker (and the file .docker/config.json) will be created when you successfully authenticate
# to a Docker registry e.g Docker Hub running the command docker login.
mkdir -p ~/.docker

# Config for pulling images from ecr later
cat <<EOH > ~/.docker/config.json
{
     "credsStore" : "ecr-login",
     "HttpHeaders" : {
       "User-Agent" : "Docker-Client/19.03.1 (XXXXXX)"
     },
     "auths" : {
       "$AWS_ACCOUNT_ID.dkr.ecr.$AWS_DEFAULT_REGION.amazonaws.com" : {}
     },
     "credHelpers": {
       "$AWS_ACCOUNT_ID.dkr.ecr.$AWS_DEFAULT_REGION.amazonaws.com" : "ecr-login"
     }
  }
EOH


# Check that config was successfully made
if [[ ! -s ~/.docker/config.json ]]; then
  echo "failed to create ~/.docker/config.json, exiting"
  exit 1
fi


# Uninstall old dependencies
for pkg in docker.io docker-doc docker-compose docker-compose-v2 podman-docker containerd runc; do
  sudo apt-get remove $pkg
done

# Add Docker's official GPG key:
sudo apt-get update
sudo apt-get install ca-certificates curl gnupg
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg

# Add the repository to Apt sources:
echo \
  "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update

sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# add current user to users allowed to run docker without sudo
sudo groupadd docker
sudo usermod -aG docker "$USER"

sudo apt-get install amazon-ecr-credential-helper

newgrp docker <<'EOH'
# new group exits/starts a new shell so i gotta do this heredoc nonsense

function pull_aws_image() {
  echo "pulling: $1"

  docker pull "$1"

  if [[ $? == 1 ]]; then
    echo "Failed to pull $1 exiting."
    exit 1
  fi
}

# Check docker was installed properly and doesn't need sudo to run
if ! docker -v; then
  echo "Failed to install docker"
  exit 1
fi


# get helper image, setup volume for watchtower auth
docker volume create helper

HELPER_IMAGE_FULL_NAME="$ECR_ADDRESS/$IMAGE_REPO_NAME:$HELPER_IMAGE_TAG"
pull_aws_image "$HELPER_IMAGE_FULL_NAME" || exit 1

docker run -d --rm --name aws-cred-helper \
  --volume helper:/go/bin/ $HELPER_IMAGE_FULL_NAME


# get app image
APP_IMAGE_FULL_NAME="$ECR_ADDRESS/$IMAGE_REPO_NAME:$IMAGE_TAG"
pull_aws_image "$APP_IMAGE_FULL_NAME" || exit 1

# Start app container
docker run -d --name watched-container \
       --label "com.centurylinklabs.watchtower.enable=true" \
       --publish 80:80 \
       --restart always \
       "$APP_IMAGE_FULL_NAME"


# Start watchtower (handle image updates)
docker run -d --name watchtower \
       --volume /var/run/docker.sock:/var/run/docker.sock \
       --volume ~/.docker/config.json:/config.json \
       --volume helper:/go/bin \
       --env PATH="/go/bin" \
       --env WATCHTOWER_HTTP_API_TOKEN="$WATCHTOWER_KEY" \
       --env AWS_REGION="$AWS_DEFAULT_REGION" \
       --env AWS_ACCESS_KEY_ID="$AWS_ACCESS_KEY_ID" \
       --env AWS_SECRET_ACCESS_KEY="$AWS_SECRET_ACCESS_KEY" \
       --label "com.centurylinklabs.watchtower.enable=false" \
       --publish 8080:8080 \
       --restart always \
       containrrr/watchtower \
       --debug --http-api-update
EOH