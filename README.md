# Kubernetes Generated Secret
[![GitHub](https://img.shields.io/github/license/phillebaba/kubernetes-generated-secret)](https://github.com/phillebaba/kubernetes-generated-secret)
[![Travis (.org)](https://img.shields.io/travis/phillebaba/kubernetes-generated-secret)](https://travis-ci.org/phillebaba/kubernetes-generated-secret)
[![Go Report Card](https://goreportcard.com/badge/github.com/phillebaba/kubernetes-generated-secret)](https://goreportcard.com/report/github.com/phillebaba/kubernetes-generated-secret)
[![Docker Pulls](https://img.shields.io/docker/pulls/phillebaba/kubernetes-generated-secret)](https://hub.docker.com/r/phillebaba/kubernetes-generated-secret)

Kubernetes controller to easily generate random secrets inside your cluster. The project makes use of [`crypto/rand`](https://golang.org/pkg/crypto/rand/) to generate random values.

## Install
Easiest way is to add a git reference in your `kustomization.yaml` file.
```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
- github.com/phillebaba/kubernetes-generated-secret//config/default
```

Or you can add the CRD and Deploy the controller in your cluster manually.
```bash
kustomize build config/default | kubectl apply -f -
```

## How to use
A `Secret` is generated from a `GeneratedSecret` that configures the length, character content, and additional metadata of the secret. The `GeneratedSecret` is the parent of the `Secret` it creates, meaning that the `Secret` will be deleted when the `GeneratedSecret` is deleted.

### Simple random secret
Below is all you need to generate a `Secret` with a random value. It contains a data field in the spec that specifies the key length and value options for the generated value.
```yaml
apiVersion: core.phillebaba.io/v1alpha1
kind: GeneratedSecret
metadata:
  name: generatedsecret-sample
spec:
  data:
  - key: test
    length: 100
    options:
      - Uppercase
      - Lowercase
      - Numbers
      - Symbols
```

### Secret metadata
There is an optional `secretMetadata` that can be set. The metadata specified will propogate to the generated `Secret` with the exception of the name and namespace which is inherited by the parent `GeneratedSecret`.
```yaml
apiVersion: core.phillebaba.io/v1alpha1
kind: GeneratedSecret
metadata:
  name: generatedsecret-sample
spec:
  secretMetadata:
    labels:
      app: foobar
  data:
  - key: test
    length: 100
    options:
      - Uppercase
      - Lowercase
      - Numbers
      - Symbols
```

The resulting `Secret` will look like shown below.
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: generatedsecret-sample
  labels:
    app: foobar
spec:
  data:
    test: <RANDOM_VALUE>
```

### Multiple secrets
It is also possible to generate a `Secret` with multiple keys in it.
```yaml
apiVersion: core.phillebaba.io/v1alpha1
kind: GeneratedSecret
metadata:
  name: generatedsecret-sample
spec:
  data:
  - key: foo
    length: 100
    options:
      - Uppercase
      - Lowercase
  - key: bar
    length: 50
    options:
      - Numbers
      - Symbols
```

Each key will receive a different random value.
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: generatedsecret-sample
  labels:
    app: foobar
spec:
  data:
    foo: <RANDOM_VALUE_1>
    bar: <RANDOM_VALUE_2>
```

## Development
The project is setup with [Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) so it is good to install it as the integration tests depend on it, follow the [installation instructions](https://book.kubebuilder.io/quick-start.html#installation).

To simplify development it helps to use a local cluster, [Kind](https://github.com/kubernetes-sigs/kind) is a good example of such a tool. Given that a cluster is configured in a kubeconfig file run the following command to install the CRD.
```bash
make install
```

Then run the controller, the following command will run the controller binary.
```bash
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
