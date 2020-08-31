# Minibroker Integration Tests for CAP

This repository holds the integration tests for running Minibroker as a service
broker with Cloud Foundry.

## Running the tests

1. Install KubeCF.
2. Install Minibroker.
3. Build the image and chart:
```
./build/build.sh
```
Optionally, build the image using Minikube's Docker daemon:
```
MINIKUBE=true ./build/build.sh
```
4. Run the tests with Helm:
```
kubectl create namespace mits
helm install mits \
  --namespace mits output/mits-<version>.tgz \
  --set "config.cf.admin.username=admin" \
  --set "config.cf.admin.password=<password for the admin user>" \
  --set "config.cf.api.endpoint=<URL for the KubeCF API>"
```

## Creating a new release

MITS uses GitHub Actions to create a new release.

1. Navigate to the [Actions](https://github.com/SUSE/minibroker-integration-tests/actions)
  page and select the Release pipeline.
2. Trigger the pipeline manually using the `workflow_dispatch` event trigger.
