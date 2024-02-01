# ec-kube

As part of my weekend AI streams, I am building a Kubernetes controller that will manage Embedchain AI apps.

Streams:
- Part 1: [Watch here](https://www.youtube.com/watch?v=X-irXixeo1Y)
- Part 2: [Join here](https://www.youtube.com/watch?v=q2sG9cRJh-w)

## What are we building?

1. Create the Kubernetes controller (I'll use the kubebuilder)

```yaml
spec:
  secretRef:
    name: this-is-my-secret

  config: |-
    {
      // Embedchain config
    }
```

2. Embedchain API Docker image

## Running the image in Kubernetes

Create a ConfigMap that holds the Embedchain config file:

```sh
kubectl create configmap ec-config --from-file=ec-image/config/config.yaml
```

Create a Secret that holds the `OPENAI_API_KEY` (for now):

```sh
kubectl create secret generic ec-secret --from-literal='OPENAI_API_KEY=${OPENAI_API_KEY}'
```

## Kubernetes controller

```shell
kubebuilder init --domain learncloudnative.com --repo learncloudnative.com/aiapps
kubebuilder create api --group aiapps --version v1 --kind EmbedchainApp
```
