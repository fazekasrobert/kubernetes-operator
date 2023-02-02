# operator
Basic kubernetes operator built with Kubebuilder to deploy httpd.

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**NOTE:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Prerequisites
Deploy `ingress-nginx` on your cluster:

```sh
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.5.1/deploy/static/provider/cloud/deploy.yaml
```

### Test It Out
1. Install the CRDs into the cluster:

```sh
make install
```

**NOTE:** `make install` also creates a namespace, called `sandbox-system`. Each CRD should be placed in this namespace.

2. Create a sample CRD on the cluster:

```sh
kubectl apply -f config/samples/operator_sample.yaml

```

3. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

4. Delete operator and CRDs:

```sh
make uninstall
```

**NOTE:** Also, `make uninstall` deletes `sandbox-system` namespace and all of its content.

**NOTE:** Run `make --help` for more information on all potential `make` targets.

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

