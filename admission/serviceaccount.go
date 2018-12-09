package main

import (
	"fmt"

	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var serviceAccountResource = metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "serviceaccounts"}

// Don't allow new serviceAccounts to be created
func admitServiceAccount(ar v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	reviewResponse := v1beta1.AdmissionResponse{}
	reviewResponse.Allowed = false
	reviewResponse.Result = &metav1.Status{
		Reason: "not letting you create service accounts",
	}

	fmt.Printf("service account not permitted: %v\n", reviewResponse.Result.Message)
	return &reviewResponse
}
