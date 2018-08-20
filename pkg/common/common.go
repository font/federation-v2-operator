package common

import (
	"regexp"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// AddOwnerRefToObject appends the desired OwnerReference to the object
func AddOwnerRefToObject(obj metav1.Object, ownerRef metav1.OwnerReference) {
	obj.SetOwnerReferences(append(obj.GetOwnerReferences(), ownerRef))
}

// AsOwner returns an OwnerReference set as the CR object passed in.
func AsOwner(obj metav1.Object, trueVar bool) metav1.OwnerReference {
	return metav1.OwnerReference{
		APIVersion: obj.(schema.ObjectKind).GroupVersionKind().Version,
		Kind:       obj.(schema.ObjectKind).GroupVersionKind().Kind,
		Name:       obj.GetName(),
		UID:        obj.GetUID(),
		Controller: &trueVar,
	}
}

// ValidateVersion returns true if a given string matches against the regular
// expression, otherwise returns false.
func ValidateVersion(regex, s string) bool {
	re := regexp.MustCompile(regex)
	return re.MatchString(s)
}
