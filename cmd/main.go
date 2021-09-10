package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	admsv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/klog"
)

var (
	scheme       = runtime.NewScheme()
	codecs       = serializer.NewCodecFactory(scheme)
	deserializer = codecs.UniversalDeserializer()
	configFile   string
	patchedFlag  = Patch{Op: "add", Path: "/metadata/labels", Value: map[string]string{"injected": "true"}}
)

func errorResponse(err error) *admsv1.AdmissionResponse {
	return &admsv1.AdmissionResponse{
		Result: &metav1.Status{
			Message: err.Error(),
		},
	}
}

func getPodPatchs(pod corev1.Pod) []byte {
	config := loadConfig(configFile)
	podInjection := config.PodInjection

	patchs := []Patch{patchedFlag}

	for _, initContainer := range podInjection.InitContainers {
		patch := Patch{}
		patch.Op = "add"
		if pod.Spec.InitContainers != nil {
			patch.Path = "/spec/initContainers/-"
			patch.Value = initContainer
		} else {
			patch.Path = "/spec/initContainers"
			patch.Value = []corev1.Container{initContainer}
		}
		patchs = append(patchs, patch)
	}

	for _, container := range podInjection.Containers {
		patch := Patch{}
		patch.Op = "add"
		if pod.Spec.Containers != nil {
			patch.Path = "/spec/containers/-"
			patch.Value = container
		} else {
			patch.Path = "/spec/containers"
			patch.Value = []corev1.Container{container}
		}
		patchs = append(patchs, patch)
	}

	for _, volume := range podInjection.Volumes {
		patch := Patch{}
		patch.Op = "add"
		if pod.Spec.Volumes != nil {
			patch.Path = "/spec/volumes/-"
			patch.Value = volume
		} else {
			patch.Path = "/spec/volumes"
			patch.Value = []corev1.Volume{volume}
		}
		patchs = append(patchs, patch)
	}

	for index, podContainer := range pod.Spec.Containers {
		volumeMounts := podInjection.VolumeMountPatchs[podContainer.Name]
		for _, volumeMount := range volumeMounts {
			patch := Patch{}
			patch.Op = "add"
			if podContainer.VolumeMounts != nil {
				patch.Path = fmt.Sprintf("/spec/containers/%d/volumeMounts/-", index)
				patch.Value = volumeMount
			} else {
				patch.Path = fmt.Sprintf("/spec/containers/%d/volumeMounts", index)
				patch.Value = []corev1.VolumeMount{volumeMount}
			}
			patchs = append(patchs, patch)
		}
	}

	podPatchs, _ := json.Marshal(patchs)
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
	var configFileFlag = flag.String("config-file", "/config/injection.yaml", "Injection configuration")
	klog.InitFlags(nil)
	flag.Parse()
	configFile = *configFileFlag

	http.HandleFunc("/inject", serveInject)
	server := &http.Server{
		Addr: ":8443",
	}
	server.ListenAndServeTLS(*certFileFlag, *keyFileFlag)
}
