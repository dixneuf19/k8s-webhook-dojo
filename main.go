package main

import (
	"encoding/json"
	"log"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"

	"github.com/busser/k8s-webhook-dojo/handlers"
)

func main() {

	http.HandleFunc("/inject-tolerations", httpHandlerFrom(handlers.AddTolerations))

	log.Println("Listening on port 8443...")
	// err := http.ListenAndServe(":8443", nil)
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
		log.Println("Handling webhook request ...")

		// Verify the content type is accurate.
		if req.Header.Get("Content-Type") != "application/json" {
			http.Error(resp, "Only application/json is supported", http.StatusUnsupportedMediaType)
			return
		}

		var adRev admissionv1.AdmissionReview
		// Decode the admission review from the HTTP request.
		err := json.NewDecoder(req.Body).Decode(&adRev)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusBadRequest)
			return
		}

		// Call the admission handler.
		adRev.Response = handler(*adRev.Request)

		// Encode the admission review into the HTTP response.
		resp.Header().Add("Content-Type", "application/json")
		err = json.NewEncoder(resp).Encode(adRev)
		if err != nil {
			log.Printf("Error handling webhook request: %v", err)
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Print("Webhook request handled successfully")
		// io.WriteString(resp, '{"Hello, world!"}")

	}
}
