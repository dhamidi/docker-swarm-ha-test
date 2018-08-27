# Description

This repository explores how to achieve high-availability deploys using docker swarm when services take a variable amount of time to start up before accepting HTTP requests.

## Trying it out

**Prerequisites**:

- a working installation of Docker Swarm
- GNU make

Run `make healthcheck` to deploy a stack named `ha`.

Ping the HTTP services in the stack using `make watch`.

Follow the logs using `make logs`.

See the (Makefile)[./Makefile] for more information.

## [slow-starting-service]

The directory [slow-starting-service] contains a simple golang HTTP server that takes n to m seconds before accepting HTTP requests.  By default the number of seconds waited is a random number between 10 and 60 seconds.  Additionally the service exposes an HTTP endpoint under `/health`, which always returns a `200 OK` response.

A `Dockerfile` is included to build a docker image for this service.  Run `make -C slow-starting-service` from the project root directory to build the docker image.

[slow-starting-service]: ./slow-starting-service

## [naive.yml]

A naïve docker swarm configuration can be found in [naive.yml].

A docker service is created with two instances of the slow starting service.  The service is exposed on host port `:12345`

[naive.yml]: ./naive.yml


## [healthcheck]

The same as the naïve stack, except that Docker healthchecks are also configured.  The healthchecks work as follows:

1. Each container is given 10 seconds to perform any preparations for its startup.  No healthchecks are performed by Docker Swarm during this time.
2. After the initial period has ended, a health check is performed every 10 seconds.  If 3 checks in a row fail, the container is considered unhealthy and removed from the Docker Swarm load balancer.
3. If a single healthcheck takes more than 10s to respond, the check is considered failed.

The stackfile is also configured to deploy the test service one container at a time, starting new containers first before killing old containers.

Effect: new containers can be deployed with zero downtime.

[healthcheck]: ./healthcheck.yml

Suggested settings for an application that takes several minutes to start:


    healthcheck:
      interval: 10s
      timeout: 10s
      retries: 3
      start_period: 15m
