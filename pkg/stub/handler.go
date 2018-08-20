package stub

import (
	"context"

	v1alpha1 "github.com/marun/federation-v2-operator/pkg/apis/operator/v1alpha1"
	"github.com/marun/federation-v2-operator/pkg/clusterregistry"
	"github.com/marun/federation-v2-operator/pkg/federation"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
)

func NewHandler() sdk.Handler {
	return &Handler{}
}

type Handler struct {
}

// TODO(font): operator-sdk does not currently support multiple handlers for
// different resource types.
func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	var err error
	switch event.Object.(type) {
	case *v1alpha1.ClusterRegistry:
		err = clusterregistry.Handle(event)
	case *v1alpha1.FederationV2:
		err = federation.Handle(event)
	}
	return err
}
