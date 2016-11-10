# Refresher
## Auto-restart your k8s pods when dependent configuration maps are updated

To use:

- Deploy configmap-updater in your kubernetes cluster somewhere
- Label your pods with the `astuart.co/configMapBehavior` label, with the following acceptable values:
  - `Delete`: Delete this pod when any configmap volumes are updated

Upon updating a configmap used by the labeled pods, all instances will be
deleted, allowing your deployment or replicationcontroller to scale them back
up.

## TODO

- Use a rollout
- Allow label value to control behavior
