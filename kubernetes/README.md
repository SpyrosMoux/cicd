### Install Keda

```shell
kustomize build kustomize/kubernetes/overlays/keda --enable-helm > k8s-manifest.yaml
kubectl create namespace keda
kubectl apply --server-side -f k8s-manifest.yaml
```

### Deploy app

**Prerequisites**

- Secret for rabbitmq:

  ```yaml
  apiVersion: v1
  kind: Secret
  metadata:
    name: keda-rabbitmq-secret
  data:
    host: <base64 secret>
  ```

**Run the app**

```shell
kubectl apply -k kubernetes/kustomize/overlays/flowforge
```
