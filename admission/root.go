package main

import (
	"fmt"

	"k8s.io/api/admission/v1beta1"
)

func admitRoot(ar v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	fmt.Printf("resource: %v\n", ar.Request.Resource)
	reviewResponse := v1beta1.AdmissionResponse{}
	reviewResponse.Allowed = true
	return &reviewResponse
}
