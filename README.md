# Minibroker Integration Tests for CAP

This repository holds the integration tests for running Minibroker as a service
broker with Cloud Foundry.

## Running the tests

1. Install KubeCF.
2. Install Minibroker.
3. (optional) Build the image using the Minibroker Docker daemon:
```
MINIKUBE=true ./build/image.sh
```
4. Run the tests with Helm:
```
helm install mits \
  --namespace mits chart/mits/ \
  --set "cf.admin.username=admin" \
  --set "cf.admin.password=<password for the admin user>" \
  --set "cf.api=<URL for the KubeCF API>"
```
