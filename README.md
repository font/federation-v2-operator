# Cluster Operator for FederationV2

This repo contains the in-progress prototype of a cluster operator for
FederationV2 and ClusterRegistry.  Since Federation and the Cluster Registry
CRD are currently cluster-scoped, only a single instance of each should be
deployed to a given cluster.

FederationV2 expects to be run in the federation-system namespace and
have sufficient access to read and add finalizers to all federation
resources in the cluster.  Since by default an operator can only
create resources in the namespace in which it is deployed, the
federation operator also needs to be deployed to the federation-system
namespace.


## Prepare the federation-system namespace

- Create namespace `federation-system`

```bash
kubectl create namespace federation-system
```

- Grant the federation control plane the permissions it needs to run.

```bash
# TODO(marun) Limit cluster-scoped permissions to federation resources (template/placement/etc)
kubectl create clusterrolebinding federation-admin \
  --clusterrole=cluster-admin --serviceaccount=federation-system:default
```

## Deploying the operator

The federation operator can be deployed manually or via OLM.

 - Manually:

```bash
kubectl create -n federation-system -f deploy/rbac.yaml
kubectl create -f deploy/crd.yaml -f deploy/cluster-registry/crd.yaml
kubectl create -n federation-system -f deploy/operator.yaml
```

 - Via OLM (must be [installed](https://github.com/operator-framework/operator-lifecycle-manager/blob/master/Documentation/install/install.md)):

```bash
kubectl create -f deploy/olm-catalog/crd.yaml
kubectl create -n federation-system -f deploy/olm-catalog/csv.yaml
```

- Checking that the operator is running

```bash
# Look for a running pod with the name prefix of 'federation-v2-operator-'
kubectl get pods -n federation-system
```

## Deploying the Cluster Registry

- Create the operator CR for the cluster registry in the federation-system
  namespace:

```bash
kubectl create -n federation-system -f deploy/cluster-registry/cr.yaml
```

- Check that the cluster registry `clusters` resource is recognized and reports
  `No resources found`:

```bash
kubectl get clusters.clusterregistry.k8s.io
```

- Check that the `kube-multicluster-public` namespace has been created:

```bash
kubectl get ns kube-multicluster-public
```

## Deploying the federation control plane

- Create the operator CR in the federation-system namespace

```bash
kubectl create -n federation-system -f deploy/cr.yaml
```

- Check that the federation control manager is running

```bash
# Look for a running pod with the name prefix of 'federation-controller-manager-'
kubectl get pods -n federation-system
```
