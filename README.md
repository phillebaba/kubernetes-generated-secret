# Kubernetes Generated Secret
[![GitHub](https://img.shields.io/github/license/phillebaba/kubernetes-generated-secret)](https://github.com/phillebaba/kubernetes-generated-secret)
[![Travis (.org)](https://img.shields.io/travis/phillebaba/kubernetes-generated-secret)](https://travis-ci.org/phillebaba/kubernetes-generated-secret)
[![Go Report Card](https://goreportcard.com/badge/github.com/phillebaba/kubernetes-generated-secret)](https://goreportcard.com/report/github.com/phillebaba/kubernetes-generated-secret)
[![Docker Pulls](https://img.shields.io/docker/pulls/phillebaba/kubernetes-generated-secret)](https://hub.docker.com/r/phillebaba/kubernetes-generated-secret)

Controller to easily generate random secret values.

## Install
Add the CRD to the cluster.
```bash
kustomzie build config/crd | kubectl apply -f -
```

Deploy the controller in your cluster
```bash
kustomize build config/default | kubectl apply -f -
```

## How to use
### Create a new random secret
```yaml
```

### Constrain secret generation
```yaml
```

### Multiple secrets
```yaml
```

## Development
The project is setup with [Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) so it is good to install it as the integration tests depend on it, follow the [installation instructions](https://book.kubebuilder.io/quick-start.html#installation).

To simplify development it helps to use a local cluster, [Kind](https://github.com/kubernetes-sigs/kind) is a good example of such a tool. Given that a cluster is configured in a kubeconfig file run the following command to install the CRD.
```bash
make install
```

Then run the controller, the following command will run the controller binary. Make sure to disable any webhooks as they will not work when running outside of the cluster.
```bash
export ENABLE_WEBHOOKS=false
make run
```

Or you can run the controller inside of the cluster, like you would when actually deploying it.
```bash
make deploy
```

Run the test rule to run the integration tests.
```bash
make test
```

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
