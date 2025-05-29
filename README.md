# CRUD application in GO

A simple API for managing Devices

## Installation

Add PROD.env and TEST.env file at the root of the project, follow the .env.example file pattern

COMMANDS for running the server
 
 ```sh
go get -d ./...
```
at the root

 ```sh
go run app.go
```

## Endpoints List

This project is deployed on a VM, test it accessing the link below:

https://apidog.com/apidoc/shared/970c9244-c847-42d1-95d3-b7d7c5a17687/status-17486321e0

## Docker

COMMANDS for building and push the docker image

 ```sh
sudo docker rmi --force yourdockerhubuser/test-devices-api

sudo docker build . --no-cache -t yourdockerhubuser/test-devices-api:latest

sudo docker push yourdockerhubuser/test-devices-api:latest
```

for running the docker container

```ssh
sudo docker run -d -t i- p 3001:3001 -e EXTERNAL_AUTH="your auth" -e MONGO_URL="your mongo url" -e PORT="0000" -e DB_NAME="your db name" yourdockerhubuser/test-devices-api:latest
```

## Infra

COMMANDS for deploying infra on MagaluCloud TF

 ```sh
export TF_VAR_mgc_api_key="your api key"
export TF_VAR_region="your region"
export TF_VAR_mgc_vpc_i_id="your vpc id if fixed"
export TF_VAR_mgc_sg_id="your security group if fixed"
export TF_VAR_mgc_ssh_key="your ssh key id if fixed"

terraform init

terraform apply

```

If you want to repeat

```sh
```

## Overall documentation

*EXTERNAL_AUTH on env, the key the client should sent on headers, placed at "x-api-key" header;

Pending:

x Rate limiters
x Non-static external auth
x Handlers fully tested
x ID json tag as "id" on structs
x actions



