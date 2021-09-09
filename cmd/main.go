package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	admsv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/klog"
)

const (
	defaultPatch string = `[{"op":"add","path":"/metadata/annotations","value":{"injected":"true"}}]`
)

var (
	scheme       = runtime.NewScheme()
	codecs       = serializer.NewCodecFactory(scheme)
	deserializer = codecs.UniversalDeserializer()
	patchFile    string
)

func errorResponse(err error) *admsv1.AdmissionResponse {
	return &admsv1.AdmissionResponse{
		Result: &metav1.Status{
			Message: err.Error(),
		},
	}
}

func getPodPatchs(pod corev1.Pod) []byte {
	podPatchs, err := os.ReadFile(patchFile)
	if err != nil {
		return []byte(defaultPatch)
	}
	return podPatchs
}

func mutatePods(admissionRequest admsv1.AdmissionRequest) *admsv1.AdmissionResponse {
	klog.V(2).Info("mutating pods")
	podResource := metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
	if admissionRequest.Resource != podResource {
		err := fmt.Errorf("expect resource to be %s", podResource)
		klog.Error(err)
		return errorResponse(err)
	}

	raw := admissionRequest.Object.Raw
	pod := corev1.Pod{}
	if _, _, err := deserializer.Decode(raw, nil, &pod); err != nil {
		klog.Error(err)
		return errorResponse(err)
	}

	patchType := admsv1.PatchTypeJSONPatch
	admissionResponse := admsv1.AdmissionResponse{
		UID:       admissionRequest.UID,
		Allowed:   true,
		Patch:     getPodPatchs(pod),
		PatchType: &patchType,
	}

	return &admissionResponse
}

func serveInject(w http.ResponseWriter, r *http.Request) {
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		klog.Errorf("contentType=%s, expect application/json", contentType)
		return
	}

	klog.V(2).Info(fmt.Sprintf("handling request: %s", body))

	requestedAR := admsv1.AdmissionReview{}
	responseAR := admsv1.AdmissionReview{}

	if _, _, err := deserializer.Decode(body, nil, &requestedAR); err != nil {
		klog.Error(err)
		responseAR.Response = errorResponse(err)
	} else {
		responseAR.Response = mutatePods(*requestedAR.Request)
	}

	responseAR.TypeMeta = requestedAR.TypeMeta

	klog.V(2).Info(fmt.Sprintf("sending response: %v", responseAR.Response))

	respBytes, err := json.Marshal(responseAR)
	if err != nil {
		klog.Error(err)
	}
	if _, err := w.Write(respBytes); err != nil {
		klog.Error(err)
	}
}

func main() {
	var certFileFlag = flag.String("cert-file", "/certs/server.cer", "TLS certificate")
	var keyFileFlag = flag.String("key-file", "/certs/server.key", "TLS private key")
	var patchFileFlag = flag.String("patch-file", "/patchs/patch.json", "TLS private key")
	klog.InitFlags(nil)
	flag.Parse()
	patchFile = *patchFileFlag

	http.HandleFunc("/inject", serveInject)
	server := &http.Server{
		Addr: ":8443",
	}
	server.ListenAndServeTLS(*certFileFlag, *keyFileFlag)
}
