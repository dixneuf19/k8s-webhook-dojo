package handlers

import (
	"encoding/json"
	"fmt"
	"log"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	podGroupVersionKind = metav1.GroupVersionKind{Version: "v1", Kind: "Pod"}
	patchType           = admissionv1.PatchTypeJSONPatch
)

// TolerationKey is the injected toleration's key.
const TolerationKey = "padok.fr/namespace"

// AddTolerations responds to an AdmissionRequest for a Pod with a patch that
// adds a toleration based on the pod's namespace.
func AddTolerations(request admissionv1.AdmissionRequest) *admissionv1.AdmissionResponse {
	// Make sure that the request's Kind is for a Pod resource.
	if request.Kind != podGroupVersionKind {
		return admissionResponseError(fmt.Errorf("expect resource to be %s", podGroupVersionKind))
	}

	// Decode the Pod from the request.
	var pod corev1.Pod
	err := json.Unmarshal(request.Object.Raw, &pod)
	if err != nil {
		return admissionResponseError(err)
	}

	// Prepare a JSON patch.
	var patch JSONPatch

	// If the Pod has no tolerations, set an initial value: an empty slice of
	// tolerations.
	if len(pod.Spec.Tolerations) == 0 {
		patch.Append(JSONPatchOperation{
			Op:    "add",
			Path:  "/spec/tolerations",
			Value: make([]corev1.Toleration, 0),
		})
	}

	// Define the toleration to add.
	toleration := corev1.Toleration{
		Key:    TolerationKey,
		Value:  request.Namespace,
		Effect: corev1.TaintEffectNoSchedule,
	}

	// Check if the pod already has the toleration.
	podHasToleration := false
	for _, t := range pod.Spec.Tolerations {
		if toleration.MatchToleration(&t) {
			log.Printf("pod has already toleration %s", t.String())
			podHasToleration = true
			break
		}
	}

	// Add the toleration if it is missing.
	if !podHasToleration {
		log.Printf("add toleration %s to pod", toleration.String())
		patch.Append(JSONPatchOperation{
			Op:    "add",
			Path:  fmt.Sprintf("/spec/tolerations/%d", len(pod.Spec.Tolerations)),
			Value: toleration,
		})
	}

	// Encode the patch as JSON.
	patchBytes, err := json.Marshal(patch)
	if err != nil {
		return admissionResponseError(fmt.Errorf("could not marshal JSON patch: %v", err))
	}

	// Prepare a response.
	resp := admissionv1.AdmissionResponse{
		UID:     request.UID,
		Allowed: true,
	}

	// Include the patch in the response.
	resp.Patch = patchBytes
	resp.PatchType = &patchType

	return &resp

}

func admissionResponseError(err error) *admissionv1.AdmissionResponse {
	return &admissionv1.AdmissionResponse{
		Result: &metav1.Status{
			Message: err.Error(),
		},
	}
}
