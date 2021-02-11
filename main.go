package main

import (
	"log"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
)

func main() {
	log.Println("Listening on port 8443...")
	err := http.ListenAndServeTLS(
		":8443",
		"/tmp/toleration-injector/serving-certs/tls.crt",
		"/tmp/toleration-injector/serving-certs/tls.key",
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
}

// AdmissionHandlerFunc responds to a Kubernetes admission request.
type AdmissionHandlerFunc func(admissionv1.AdmissionRequest) *admissionv1.AdmissionResponse

func httpHandlerFrom(handler AdmissionHandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		// Verify the content type is accurate.

		// Decode the admission review from the HTTP request.

		// Call the admission handler.

		// Encode the admission review into the HTTP response.
	}
}
