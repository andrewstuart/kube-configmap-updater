# Refresher
## Auto-refresh your k8s pods periodically

To use:

- Deploy refresher in your kubernetes cluster
- Label your pods with the `astuart.co/updateFrequency` label. Possible values:
  - Any valid Golang [duration](https://golang.org/pkg/time/#ParseDuration)
  - "daily"
  - "weekly"
