package main

import (
	"os"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/klog"
)

type (
	InjectionConfig struct {
		PodInjection PodInjectionConfig `json:"podInjection,omitempty"`
	}

	PodInjectionConfig struct {
		InitContainers  []corev1.Container        `json:"initContainers,omitempty"`
		Containers      []corev1.Container        `json:"containers,omitempty"`
		Volumes         []corev1.Volume           `json:"volumes,omitempty"`
		ContainerPatchs map[string]ContainerPatch `json:"containerPatchs,omitempty"`
	}

	ContainerPatch struct {
		VolumeMounts []corev1.VolumeMount `json:"volumeMounts,omitempty"`
	}

	Patch struct {
		Op    string      `json:"op,omitempty"`
		Path  string      `json:"path,omitempty"`
		Value interface{} `json:"value,omitempty"`
	}
)

func loadConfig(configFile string) InjectionConfig {
	config := &InjectionConfig{}
	configYaml, err := os.ReadFile(configFile)
	if err != nil {
		klog.Error(err)
		return *config
	}

	err = yaml.Unmarshal(configYaml, config)
	if err != nil {
		klog.Error(err)
	}

	return *config
}
