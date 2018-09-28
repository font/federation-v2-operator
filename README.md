# Operator for FederationV2

The `deploy` directory contains the in-progress prototype of a cluster or
namespaced operator for FederationV2.  When Federation is cluster-scoped, only
a single instance of the control plane should be deployed to a given cluster.

FederationV2 can be run in any namespace but must have sufficient access to
read and add finalizers to all federation resources in the cluster (if using
cluster-scoped deployment) or namespace (if using namespace-scoped deployment).

## Cluster or Namespace Scoped Deployment

Set the appropriate variable depending on whether you'd like to perform a
cluster-scoped deployment:

```bash
CLUSTER_OR_NAMESPACE=cluster
```

or namespace-scoped deployment:

```bash
CLUSTER_OR_NAMESPACE=namespace
```

## Prepare the namespaces

- Create namespace `FEDERATION_NAMESPACE`:

```bash
FEDERATION_NAMESPACE=federation-system
kubectl create namespace ${FEDERATION_NAMESPACE}
```

- Grant the federation control plane the permissions it needs to run.

```bash
# TODO(marun) Limit cluster-scoped permissions to federation resources (template/placement/etc)
# TODO(font) use OLM cluster permissions.
if [[ ${CLUSTER_OR_NAMESPACE} == 'cluster' ]]; then
    kubectl create clusterrolebinding federation-admin \
      --clusterrole=cluster-admin --serviceaccount=${FEDERATION_NAMESPACE}:federation-controller-manager
fi
```

## Deploying the federation control plane operator

The federation operator can be deployed manually or via OLM.

 - Manually:

```bash
kubectl create -n ${FEDERATION_NAMESPACE} -f deploy/${CLUSTER_OR_NAMESPACE}/rbac.yaml
kubectl create --validate=false -f deploy/${CLUSTER_OR_NAMESPACE}/crd.yaml
kubectl create -n ${FEDERATION_NAMESPACE} -f deploy/${CLUSTER_OR_NAMESPACE}/operator.yaml
```

 - Via OLM (must be [installed](https://github.com/operator-framework/operator-lifecycle-manager/blob/master/Documentation/install/install.md)):

```bash
kubectl create --validate=false -f deploy/${CLUSTER_OR_NAMESPACE}/olm-catalog/crd.yaml
kubectl create -n ${FEDERATION_NAMESPACE} -f deploy/${CLUSTER_OR_NAMESPACE}/olm-catalog/csv.yaml
```

- Checking that the federation-v2 controller manager is running

```bash
# Look for a running pod with the name prefix of 'federation-controller-manager-'
kubectl get pods -n ${FEDERATION_NAMESPACE}
```

## Enabling the federation controllers

There are numerous custom resources that you can create to enable various
controllers. The following lists the controllers and how to enable or disable
them.

### FederatedType Controllers

TBD

### MultiClusterIngressDNS Controller

TBD

### MultiClusterServiceDNS Controller

TBD

### SchedulingPreferences Controller

TBD

## Cleanup

- Delete namespaces:

```bash

kubectl delete namespace ${FEDERATION_NAMESPACE}
```

- Remove the federation control plane permissions it needs to run:

```bash
if [[ ${CLUSTER_OR_NAMESPACE} == 'cluster' ]]; then
    kubectl delete clusterrolebinding federation-admin
fi
```

- Remove required CRDs:

```bash
kubectl delete -f deploy/${CLUSTER_OR_NAMESPACE}/crd.yaml
```
