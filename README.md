# Cluster Operator for FederationV2

This repo contains the in-progress prototype of a cluster operator for
FederationV2.  Since Federation is currently cluster-scoped, only a
single instance of the control plane should be deployed to a given
cluster.

FederationV2 expects to be run in the federation-system namespace and
have sufficient access to read and add finalizers to all federation
resources in the cluster.  Since by default an operator can only
create resources in the namespace in which it is deployed, the
federation operator also needs to be deployed to the federation-system
namespace.


## Prepare the namespaces

- Create namespace `federation-system` and `kube-multicluster-public`:

```bash
kubectl create namespace federation-system
kubectl create namespace kube-multicluster-public
```

- Grant the federation control plane the permissions it needs to run.

```bash
# TODO(marun) Limit cluster-scoped permissions to federation resources (template/placement/etc)
kubectl create clusterrolebinding federation-admin \
  --clusterrole=cluster-admin --serviceaccount=federation-system:federation-controller-manager
```

## Deploying the federation control plane operator

The federation operator can be deployed manually or via OLM.

 - Manually:

```bash
kubectl create -n federation-system -f deploy/rbac.yaml
kubectl create --validate=false -f deploy/crd.yaml
kubectl create -n federation-system -f deploy/operator.yaml
```

 - Via OLM (must be [installed](https://github.com/operator-framework/operator-lifecycle-manager/blob/master/Documentation/install/install.md)):

```bash
kubectl create --validate=false -f deploy/olm-catalog/crd.yaml
kubectl create -n federation-system -f deploy/olm-catalog/csv.yaml
```

- Checking that the federation-v2 operator is running

```bash
# Look for a running pod with the name prefix of 'federation-controller-manager-'
kubectl get pods -n federation-system
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

kubectl delete namespace federation-system
kubectl delete namespace kube-multicluster-public
```

- Remove the federation control plane permissions it needs to run:

```bash
kubectl delete clusterrolebinding federation-admin
```

- Remove required CRDs:

```bash
kubectl delete -f deploy/crd.yaml
```
