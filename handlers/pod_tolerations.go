package handlers

import (
	"errors"

	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TolerationKey is the injected toleration's key.
const TolerationKey = "padok.fr/namespace"

// AddTolerations responds to an AdmissionRequest for a Pod with a patch that
// adds a toleration based on the pod's namespace.
func AddTolerations(request admissionv1.AdmissionRequest) *admissionv1.AdmissionResponse {
	// Make sure that the request's Kind is for a Pod resource.

	// Decode the Pod from the request.

	// Prepare a JSON patch.

	// If the Pod has no tolerations, set an initial value: an empty slice of
	// tolerations.

	// Define the toleration to add.

	// Check if the pod already has the toleration.

	// Add the toleration if it is missing.

	// Encode the patch as JSON.

	// Prepare a response.

	// Include the patch in the response.

	return admissionResponseError(errors.New("Not implemented"))
}

func admissionResponseError(err error) *admissionv1.AdmissionResponse {
	return &admissionv1.AdmissionResponse{
		Result: &metav1.Status{
			Message: err.Error(),
		},
	}
}
