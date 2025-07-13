# k8s-operator-poc

## Description

## Getting Started

### Prerequisites
- go version v1.24.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

### To Deploy on the cluster
**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/guestbook:tag
```

**NOTE:** This image ought to be published in the personal registry you specified.
And it is required to have access to pull the image from the working environment.
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/guestbook:tag
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

>**NOTE**: Ensure that the samples has default values to test it out.

### To Uninstall
**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Project Distribution

Following the options to release and provide this solution to the users.

### By providing a bundle with all YAML files

1. Build the installer for the image built and published in the registry:

```sh
make build-installer IMG=<some-registry>/guestbook:tag
```

**NOTE:** The makefile target mentioned above generates an 'install.yaml'
file in the dist directory. This file contains all the resources built
with Kustomize, which are necessary to install this project without its
dependencies.

2. Using the installer

Users can just run 'kubectl apply -f <URL for YAML BUNDLE>' to install
the project, i.e.:

```sh
kubectl apply -f https://raw.githubusercontent.com/<org>/guestbook/<tag or branch>/dist/install.yaml
```

### By providing a Helm Chart

1. Build the chart using the optional helm plugin

```sh
kubebuilder edit --plugins=helm/v1-alpha
```

2. See that a chart was generated under 'dist/chart', and users
can obtain this solution from there.

**NOTE:** If you change the project, you need to update the Helm Chart
using the same command above to sync the latest changes. Furthermore,
if you create webhooks, you need to use the above command with
the '--force' flag and manually ensure that any custom configuration
previously added to 'dist/chart/values.yaml' or 'dist/chart/manager/manager.yaml'
is manually re-applied afterwards.

## Contributing
// TODO(user): Add detailed information on how you would like others to contribute to this project

**NOTE:** Run `make help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

# Command used to generate Controller/CRD
kubebuilder create api --group webapp --version v1 --kind Guestbook

# Command used to generate webhook 
kubebuilder create webhook --group webapp --version v1 --kind Guestbook --defaulting false --programmatic-validation true

# How to test webhook?
For fast iteration its best to run the webhook locally of course, as it allows you to quickly compile and run avoiding CICD build times.

Step 1:
Go into your webhook configuration, for this repo it is located under /config/webhook/manifests.yaml (its contains both mutatation and validation configurations)

Set the url property under the client config, make sure the url contains localhost or your local ip (mine is 10.1.1.6)

Step 2:
The kubernetes API expects TLS so if you try to do anything without it, you may find an error like the following

Error from server (InternalError): error when creating "webapp_v1_guestbook.yaml": Internal error occurred: failed calling webhook "mguestbook-v1.kb.io": failed to call webhook: Post "https://10.1.1.6:9443/mutate-webapp-my-domain-v1-guestbook?timeout=10s": tls: failed to verify certificate: x509: cannot validate certificate for 10.1.1.6 because it doesn't contain any IP SANs

So lets fix that ahead of time, in this example I am using minikube so lets create a certificate for my webhook server
and lets sign it by the minikube certificate authority that is generated when we created our minikube cluster.


create a file called **csr.cnf** and insert below text
```
[req]
default_bits       = 2048
prompt             = no
default_md         = sha256
distinguished_name = dn
req_extensions     = req_ext

[dn]
CN = <your-common-name>

[req_ext]
subjectAltName = @alt_names

[alt_names]
IP.1 = <your-local-ip-here>
```

These commands create your server.csr as well as the tls.crt and tls.key file that you will need. 
It also signs it using the minikube certificate authority. Note: These commmands asssume everything is in one directory.
```
openssl req -new -nodes -out server.csr -newkey rsa:2048 -keyout tls.key -config csr.conf
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out tls.crt -days 365 -extensions req_ext -extfile csr.conf
```