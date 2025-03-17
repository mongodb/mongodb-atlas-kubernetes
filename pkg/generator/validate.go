package generator

import (
	"context"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/validation"
)

func ValidateCRD(ctx context.Context, crd *apiextensions.CustomResourceDefinition) error {
	errorList := validation.ValidateCustomResourceDefinition(ctx, crd)
	if len(errorList) > 0 {
		return errorList.ToAggregate()
	}
	return nil
}
