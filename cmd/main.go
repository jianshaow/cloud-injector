package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	admsv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/klog"
)

var scheme = runtime.NewScheme()
var codecs = serializer.NewCodecFactory(scheme)

func errorResponse(err error) *admsv1.AdmissionResponse {
	return &admsv1.AdmissionResponse{
		Result: &metav1.Status{
			Message: err.Error(),
		},
	}
}

func admit(admissionRequest admsv1.AdmissionRequest) *admsv1.AdmissionResponse {
	klog.V(2).Info("admitting pods")
	podResource := metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
	if admissionRequest.Resource != podResource {
		err := fmt.Errorf("expect resource to be %s", podResource)
		klog.Error(err)
		return errorResponse(err)
	}

	raw := admissionRequest.Object.Raw
	pod := corev1.Pod{}
	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(raw, nil, &pod); err != nil {
		klog.Error(err)
		return errorResponse(err)
	}
	admissionResponse := admsv1.AdmissionResponse{}
	admissionResponse.Allowed = true

	var msg string
	if v, ok := pod.Labels["app"]; ok {
		if v == "no-inject-app" {
			admissionResponse.Allowed = false
			msg = msg + "the pod contains unwanted label; "
		}
	}
	if !admissionResponse.Allowed {
		admissionResponse.Result = &metav1.Status{Message: strings.TrimSpace(msg)}
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

	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(body, nil, &requestedAR); err != nil {
		klog.Error(err)
		responseAR.Response = errorResponse(err)
	} else {
		responseAR.Response = admit(*requestedAR.Request)
	}

	responseAR.Response.UID = requestedAR.Request.UID

	klog.V(2).Info(fmt.Sprintf("sending response: %v", responseAR.Response))

	respBytes, err := json.Marshal(responseAR)
	if err != nil {
		klog.Error(err)
	}
	if _, err := w.Write(respBytes); err != nil {
		klog.Error(err)
	}
}

func configTLS(certFile string, keyFile string) *tls.Config {
	sCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		klog.Fatal(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{sCert},
	}
}

func main() {
	var certFile = flag.String("cert_file", "/certs/server.cer", "TLS certificate")
	var keyFile = flag.String("key_file", "/certs/server.key", "TLS private key")
	klog.InitFlags(nil)
	flag.Parse()

	http.HandleFunc("/inject", serveInject)
	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: configTLS(*certFile, *keyFile),
	}
	server.ListenAndServeTLS("", "")
}
