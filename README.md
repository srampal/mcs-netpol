# Network Policy for Kubernetes Multi-Cluster Services API

PoC of CRD and controller for an enhanced network policy to address multi-cluster service networking implementation of the Multi-Cluster Services API. 

Note: This is a limited proof of concept  implementation currently for basic demo purposes. The longer term goal is to have a more complete implementation based on the eventual upstream kubernetes community's  Multi-Cluster Network Policy API definition (a current [draft proposal](https://docs.google.com/document/d/1sLq_1CwfdvhmqxslQRdt7-RjLFsZRJkd/edit#heading=h.gjdgxs))

This repo defines an api (api/va1lpha1/multiclusterpolicy_types.go) and includes a golang implementation of a controller for this api. This implementation should work with multiple CNI implementations and multiple implementations of the K8s MultiCluster Services (MCS) API,  although has so far only been tested using the upstream Kubernetes, weave CNI and the [Submariner implementation](https://submariner.io/) of the MCS api. 


## Quick demo

1. Setup 2 k8s clusters with an implementation of the Multi-Cluster Services API (example
using the procedure outlined at [this link](https://submariner.io/getting-started/quickstart/kind/) 

2. Export an nginx deployment/ service from cluster 2 (example shown [here](https://submariner.io/getting-started/quickstart/kind/#verify-manually) but any standard method works. Confirm the service can be accessed by a client pod in cluster1 as described in the [same link](https://submariner.io/getting-started/quickstart/kind/#verify-manually). 

3. Clone this repo

4. Run 'make install' to deploy the multicluster policy CRD

5. Run 'make deploy IMG=srampal/mcs-netpol:0.4' to deploy the controller (if a newer version of the container is available on dockerhub, feel free to try the latest version)

6. Deploy test pods in cluster1 (kubectl apply -f config/samples/test-pods.yaml)

7. Verify that both test pods can access the remote nginx service (use 'curl nginx.default.svc.clusterset.local'  from within the bash shell of each test pod.

8. Now apply the sample multicluster policy provided (kubectl apply -f config/samples/policy-1.yaml)

9. Repeat the test from step 7, confirm that the pod labeled color:blue are unable to access the remote service whereas the pod labeled color:red can continue to access. This confirms multicluster netwpro policy filtering operation.

## Community, discussion, contribution, and support


### Code of conduct

Participation in the Kubernetes community is governed by the [Kubernetes Code of Conduct](code-of-conduct.md).

[owners]: https://git.k8s.io/community/contributors/guide/owners.md
[Creative Commons 4.0]: https://git.k8s.io/website/LICENSE
