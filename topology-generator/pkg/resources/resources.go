package resources

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/kiali/demos/topology-generator/pkg/api"
	generators "github.com/kiali/demos/topology-generator/pkg/generator"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// GenerateTopology generates a topology and the Kubernetes resources to create it
func GenerateTopology(generator api.Generator, config api.Configurations) corev1.List {
	topology := generators.GenerateTopology(generator, config)
	list := corev1.List{
		TypeMeta: metav1.TypeMeta{
			Kind:       "List",
			APIVersion: "v1",
		},
		ListMeta: metav1.ListMeta{},
	}

	for namespace, services := range topology {
		tg := true

		ns := NamespaceForMimik(namespace, config)
		ns.SetGroupVersionKind(schema.GroupVersionKind{Kind: "Namespace", Group: "", Version: "v1"})
		list.Items = append(list.Items, runtime.RawExtension{Object: ns})

		for _, service := range services {
			cm := ConfigMapForMimik(service, namespace)
			cm.SetGroupVersionKind(schema.GroupVersionKind{Kind: "ConfigMap", Group: "", Version: "v1"})
			list.Items = append(list.Items, runtime.RawExtension{Object: cm})

			svc := ServiceForMimik(service, namespace)
			svc.SetGroupVersionKind(schema.GroupVersionKind{Kind: "Service", Group: "", Version: "v1"})
			list.Items = append(list.Items, runtime.RawExtension{Object: svc})

			deploy := DeploymentForMimik(service, namespace, tg)
			deploy.SetGroupVersionKind(schema.GroupVersionKind{Kind: "Deployment", Group: "apps", Version: "v1"})
			list.Items = append(list.Items, runtime.RawExtension{Object: deploy})

			tg = false
		}
	}

	return list
}

// NamespaceForMimik creates a namespace
func NamespaceForMimik(namespace string, C api.Configurations) *corev1.Namespace {
	return &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   namespace,
			Labels: InjectionLabel(C),
		},
	}
}

// InjectionLabel get injection Labels
func InjectionLabel(C api.Configurations) map[string]string {
	injectionLabel := C.InjectionLabel
	kv := strings.Split(injectionLabel, ":")

	if len(kv) != 2 {
		log.Fatalln("Get Injection Label error")
	}

	labels := make(map[string]string)
	labels[kv[0]] = kv[1]
	labels["generated-by"] = "mimik"

	return labels
}

// ConfigMapForMimik creates a ConfigMap for an instance
func ConfigMapForMimik(s api.Service, namespace string) *corev1.ConfigMap {
	jsonData, _ := json.Marshal(s.Endpoints)
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getNameWithVersion(s.Name, s.Version),
			Namespace: namespace,
			Labels:    labelsForMimik(s.Name, s.Version),
		},
		Data: map[string]string{fmt.Sprintf("%s.json", getNameWithVersion(s.Name, s.Version)): string(jsonData)},
	}
	return cm
}

// ServiceForMimik creates a Service for an instance
func ServiceForMimik(s api.Service, namespace string) *corev1.Service {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.Name,
			Namespace: namespace,
			Labels:    labelsForMimik(s.Name, s.Version),
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{
				Name: "http",
				Port: 8080,
			},
			},
			Selector: map[string]string{
				"app": s.Name,
			},
		},
	}
	return svc
}

// DeploymentForMimik creates a Deployment for an instance
func DeploymentForMimik(s api.Service, namespace string, tg bool) *appsv1.Deployment {
	labels := labelsForMimik(s.Name, s.Version)
	annotations := annotationsForMimik(s.C)
	resources := deploymentResource(s.C)
	replicas := int32(s.C.Replicas)
	command := []string{"instance"}
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getNameWithVersion(s.Name, s.Version),
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: annotations,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:   fmt.Sprintf("%s:%s", s.C.ImageTag, s.C.ImageVersion),
						Name:    s.C.Name,
						Command: command,
						Env: []corev1.EnvVar{
							{
								Name:  "MIMIK_SERVICE_NAME",
								Value: s.Name,
							},
							{
								Name:  "MIMIK_SERVICE_PORT",
								Value: "8080",
							},
							{
								Name:  "MIMIK_ENDPOINTS_FILE",
								Value: fmt.Sprintf("/data/%s.json", getNameWithVersion(s.Name, s.Version)),
							},
							{
								Name:  "MIMIK_LABELS_FILE",
								Value: "/tmp/etc/pod_labels",
							},
						},
						Resources:       resources,
						ImagePullPolicy: corev1.PullAlways,
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      "pod-info",
								MountPath: "/tmp/etc",
							},
							{
								Name:      "endpoints",
								MountPath: "/data",
							},
						},
						Ports: []corev1.ContainerPort{{
							ContainerPort: 8080,
							Name:          "http",
						}},
					}},
					Volumes: []corev1.Volume{
						{
							Name: "pod-info",
							VolumeSource: corev1.VolumeSource{
								DownwardAPI: &corev1.DownwardAPIVolumeSource{
									Items: []corev1.DownwardAPIVolumeFile{
										{
											FieldRef: &corev1.ObjectFieldSelector{
												FieldPath: "metadata.labels",
											},
											Path: "pod_labels",
										},
									},
								},
							},
						},
						{
							Name: "endpoints",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									Items: []corev1.KeyToPath{
										{
											Key:  fmt.Sprintf("%s.json", getNameWithVersion(s.Name, s.Version)),
											Path: fmt.Sprintf("%s.json", getNameWithVersion(s.Name, s.Version)),
										},
									},
									LocalObjectReference: corev1.LocalObjectReference{
										Name: getNameWithVersion(s.Name, s.Version),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if tg {
		tge := corev1.EnvVar{
			Name:  "MIMIK_TRAFFIC_GENERATOR",
			Value: "true",
		}

		deploy.Spec.Template.Spec.Containers[0].Env = append(deploy.Spec.Template.Spec.Containers[0].Env, tge)
	}

	return deploy
}

func getNameWithVersion(name, version string) string {
	return fmt.Sprintf("%s-%s", name, version)
}

func labelsForMimik(name, version string) map[string]string {
	return map[string]string{"app": name, "version": version}
}

func annotationsForMimik(C api.Configurations) map[string]string {
	config := map[string]string{
		"sidecar.istio.io/inject":      C.EnableInjection,
		"sidecar.istio.io/proxyCPU":    C.IstioProxyRequestCPU,
		"sidecar.istio.io/proxyMemory": C.IstioProxyRequestMemory,
	}
	return config
}

func deploymentResource(C api.Configurations) corev1.ResourceRequirements {

	resourceCPU, err := resource.ParseQuantity(C.MimikRequestCPU)
	if err != nil {
		log.Fatalf("ParseQuantity error: %v", err)
	}

	resourceMemory, err := resource.ParseQuantity(C.MimikRequestMemory)
	if err != nil {
		log.Fatalf("ParseQuantity error: %v", err)
	}

	resourceLimitsCPU, err := resource.ParseQuantity(C.MimikLimitCPU)
	if err != nil {
		log.Fatalf("ParseQuantity error: %v", err)
	}

	resourceLimitsMemory, err := resource.ParseQuantity(C.MimikLimitMemory)
	if err != nil {
		log.Fatalf("ParseQuantity error: %v", err)
	}

	return corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resourceCPU,
			corev1.ResourceMemory: resourceMemory,
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resourceLimitsCPU,
			corev1.ResourceMemory: resourceLimitsMemory,
		},
	}
}
