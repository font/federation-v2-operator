package federation

import (
	"fmt"

	v1alpha1 "github.com/marun/federation-v2-operator/pkg/apis/operator/v1alpha1"
	"github.com/marun/federation-v2-operator/pkg/common"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// Name is hard-coded for now.  Allowing the name to vary would
	// allow more than one instance of the controller manager and this
	// is not supported.
	fedName = "federation-controller-manager"

	// Allowing the namespace to vary would require a canonical way to
	// discover the namespace so that client tools supporting cluster
	// join would still work.
	fedNamespace = "federation-system"
)

// Handle processes the FederationV2 Operator event object for Federation-V2.
func Handle(event sdk.Event) error {
	// Ignore the delete event since the garbage collector will clean up all secondary resources for the CR
	// All secondary resources must have the CR set as their OwnerReference for this to be the case
	if event.Deleted {
		return nil
	}

	fedv2 := event.Object.(*v1alpha1.FederationV2)

	// TODO(marun) Only output error once per instance
	if fedv2.Name != fedName {
		return fmt.Errorf("Only one FederationV2 resource can exist, and it must be named %q.", fedName)
	}

	// Create the deployment if it doesn't exist
	dep := deploymentForFederationV2(fedv2)
	err := sdk.Create(dep)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return fmt.Errorf("failed to create deployment: %v", err)
	}

	// TODO(marun) What kind of status would be useful to record?
	return nil
}

// deploymentForFederationV2 returns a federation-v2 Deployment object
func deploymentForFederationV2(f *v1alpha1.FederationV2) *appsv1.Deployment {
	ls := map[string]string{
		"api":           "federation",
		"control-plane": "controller-manager",
	}

	// TODO(marun) Switch to a stateful set to ensure only a single instance
	replicas := int32(1)

	dep := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      f.Name,
			Namespace: fedNamespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: v1.PodSpec{
					// TODO(marun) Apply the supplied configuration as feature gate paramters
					Containers: []v1.Container{{
						Image:   "quay.io/kubernetes-multicluster/federation-v2:latest",
						Name:    "controller-manager",
						Command: []string{"/root/controller-manager"},
					}},
				},
			},
		},
	}
	common.AddOwnerRefToObject(dep, common.AsOwner(f, true))
	return dep

}
