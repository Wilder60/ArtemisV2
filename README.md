# KeyRing #

## Introduction ##

KeyRing is the service that deals with the storing the data related to the password manager section of the Artemis application.  While the idea of the service is
simple, this is mostly being done to help me focus my skills on various cloud technologies like Docker, Terraform, and various GCP functionality.

The design of the application is a CRUD REST api for storing and accessing passwords and metadata in a secure fashion.  For retrieving passwords, pagination can be used
to be make easier for recyclerviews to make use of this api.

## EndPoints ##

* api/v1/keyring
  * Get
  * Post
  * Patch
  * Delete

* /health
  * Get

Note: these endpoints are subject to change

## Technologies Used ##

### Currently ###

* CloudSQL
* Cloud Logging
* Docker
* JWT

## Planned ###

* gRPC for interservice communication
* Edge Security with Auth0 Validation
* Terraform
* Running in GCP Compute Engine
