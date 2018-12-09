package main

import (
	"fmt"
	"strings"

	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var podResource = metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}

// only nginx is allowed to run as Root
func admitPod(ar v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	if ar.Request.Resource != podResource {
		err := fmt.Errorf("expect resource to be %s", podResource)
		fmt.Println(err)
		return toAdmissionResponse(err)
	}

	raw := ar.Request.Object.Raw
	pod := corev1.Pod{}
	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(raw, nil, &pod); err != nil {
		fmt.Println(err)
		return toAdmissionResponse(err)
	}

	reviewResponse := v1beta1.AdmissionResponse{}
	reviewResponse.Allowed = true

	for _, container := range pod.Spec.Containers {
		fmt.Printf("admission request for image %s\n", container.Image)
		if !strings.Contains(container.Image, "nginx") {
			if container.SecurityContext == nil || !*container.SecurityContext.RunAsNonRoot {
				reviewResponse.Allowed = false
				reviewResponse.Result = &metav1.Status{Message: "must specify RunAsNonRoot for all containers except nginx"}
				fmt.Printf("pod not permitted: %v\n", reviewResponse.Result.Message)
				return &reviewResponse
			}
		}
	}

	fmt.Println("pod permitted")
	return &reviewResponse
}
