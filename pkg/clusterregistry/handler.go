package clusterregistry

import (
	"fmt"
	"io"
	"net/http"

	v1alpha1 "github.com/marun/federation-v2-operator/pkg/apis/operator/v1alpha1"
	"github.com/marun/federation-v2-operator/pkg/common"
	"github.com/sirupsen/logrus"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"k8s.io/api/core/v1"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

const (
	// Cluster Registry name is hard-coded for now since CRDs are cluster-scoped.
	clusterRegistryName = "cluster-registry"

	// This is the canonical namespace for multi-cluster components. This
	// allows the Cluster Registry to be namespaced while designating a global
	// namespace convention for storing clusters for cluster-admin use cases.
	// Also, fedv2 expects to find clusters in this namespace. However, there
	// needs to be a mechanism for discovering namespaces in order to allow
	// them to vary.
	multiClusterPublicNamespace = "kube-multicluster-public"

	// Latest Cluster Registry Version.
	clusterRegistryLatestVersion = "v0.0.6"
)

// Handle processes the Cluster Registry Operator event object for the Cluster
// Registry.
func Handle(event sdk.Event) error {
	crObj := event.Object.(*v1alpha1.ClusterRegistry)

	// Ignore the delete event since the garbage collector will clean up all
	// secondary resources for the CR. All secondary resources must have the CR
	// set as their OwnerReference for this to be the case.
	// BUG(font): Garbage collector does not clean up CRDs with OwnerReference
	// set.
	if event.Deleted {
		// Prevent inadvertent deletion of Cluster Registry CRD when deleting
		// an invalid Cluster Registry Operator CR.
		if crObj.Name != clusterRegistryName {
			return nil
		}

		err := deleteClusterRegistryCRD()
		if err != nil {
			return err
		}
		return nil
	}

	err := validateClusterRegistryOperatorCR(crObj)
	if err != nil {
		return err
	}

	err = createClusterRegistry(crObj)
	if err != nil {
		return err
	}

	// TODO(font): create necessary RBAC resources.
	// TODO(font): any status info to capture?
	return nil
}

// deleteClusterRegistryCRD handles the deletion of the Cluster Registry CRD.
func deleteClusterRegistryCRD() error {
	crCRD := &apiextv1beta1.CustomResourceDefinition{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apiextensions.k8s.io/v1beta1",
			Kind:       "CustomResourceDefinition",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "clusters.clusterregistry.k8s.io",
		},
	}

	err := sdk.Delete(crCRD)
	if err != nil {
		return fmt.Errorf("Failed to delete Cluster Registry CRD: %v", err)
	}

	logrus.Infof("Deleted Cluster Registry CRD")
	return nil
}

// validateClusterRegistryOperatorCR verifies the name of the Cluster Registry
// Operator CR and that the version of the Cluster Registry, if specified, is
// valid.
func validateClusterRegistryOperatorCR(c *v1alpha1.ClusterRegistry) error {
	if c.Name != clusterRegistryName {
		return fmt.Errorf("Only one ClusterRegistry resource can exist, and it must be named %q", clusterRegistryName)
	} else if c.Spec.Version == "" {
		c.Spec.Version = clusterRegistryLatestVersion
	}

	crVersionRegex := `^v([0-9]+)\.([0-9]+)\.([0-9]+)`
	valid := common.ValidateVersion(crVersionRegex, c.Spec.Version)
	if !valid {
		return fmt.Errorf("Failed to validate Cluster Registry version %s using regex `%s`",
			c.Spec.Version, crVersionRegex)
	}

	logrus.Infof("Successfully validated Cluster Registry version %s", c.Spec.Version)
	return nil
}

// createClusterRegistry handles the creation of the Cluster Registry namespace
// and CRD.
func createClusterRegistry(c *v1alpha1.ClusterRegistry) error {
	crNamespace := namespaceForClusterRegistry(c)
	err := sdk.Create(crNamespace)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return fmt.Errorf("failed to create namespace %s for Cluster Registry : %v",
			crNamespace.Name, err)
	} else if err == nil {
		logrus.Infof("Created namespace %s for Cluster Registry", crNamespace.Name)
	}

	crCRD, err := clusterRegistryCRD(c)
	if err != nil {
		defer sdk.Delete(crNamespace)
		return fmt.Errorf("failed to retrieve Cluster Registry CRD : %v", err)
	}

	err = sdk.Create(crCRD)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		defer sdk.Delete(crNamespace)
		return fmt.Errorf("failed to create Cluster Registry CRD : %v", err)
	} else if err == nil {
		logrus.Infof("Created Cluster Registry CRD version %s", c.Spec.Version)
	}

	return nil
}

// namespaceForClusterRegistry returns the namespace object for the Cluster
// Registry and adds the Cluster Registry Operator CR object as its owner.
func namespaceForClusterRegistry(c *v1alpha1.ClusterRegistry) *v1.Namespace {
	crNamespace := &v1.Namespace{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Namespace",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: multiClusterPublicNamespace,
		},
	}

	common.AddOwnerRefToObject(crNamespace, common.AsOwner(c, false))
	return crNamespace
}

// clusterRegistryCRD returns the Cluster Registry CRD as an unstructured
// object after downloading, converting, and adding the Cluster Registry
// Operator CR object as its owner.
func clusterRegistryCRD(c *v1alpha1.ClusterRegistry) (*unstructured.Unstructured, error) {
	// TODO(font): what's the best way to get Cluster Registry CRD?
	clusterRegistryURL := "https://raw.githubusercontent.com/kubernetes/cluster-registry/" + c.Spec.Version + "/cluster-registry-crd.yaml"

	resp, err := http.Get(clusterRegistryURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode YAML
	crd, err := readerToObj(resp.Body)
	if err != nil {
		return nil, err
	}

	// NOTE: Adding OwnerReference to CRD does not trigger garbage collector to
	// remove it.
	common.AddOwnerRefToObject(crd, common.AsOwner(c, false))
	return crd, nil
}

// readerToObj converts the YAML reader object to JSON in order to decode it as
// an unstructured object.
func readerToObj(r io.Reader) (*unstructured.Unstructured, error) {
	decoder := yaml.NewYAMLToJSONDecoder(r)
	obj := &unstructured.Unstructured{}
	err := decoder.Decode(obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}
