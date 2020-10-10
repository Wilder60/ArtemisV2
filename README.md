# KeyRing #

## Introduction ##

KeyRing is the service that deals with the storing the data related to the password manager section of theArtemis application. This service acts as a a CRUD for the CloudSql database.  This api allows for data pagination
This project is mostly being used for me to learn different DevOps and cloud technologies like terraform, docker, and gcp.

## EndPoints ##

* /keyring
  * Get
  * Post
  * Patch
  * Delete

* /health
  * Get

Note: these endpoints are subject to change

## Technologies Used ##

### Currently ###

* GCP CloudSQL
* Docker
* JWT

## Planned ###
* Auth0 Validation
* Terraform
* GCP Compute Engine
