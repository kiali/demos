package resources

import (
	"encoding/json"
	"fmt"

	"github.com/leandroberetta/mimik/pkg/api"
	"github.com/leandroberetta/mimik/pkg/generator"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func CreateTopology(numNamespaces, numServices, numConnections, numRandomConnections int) corev1.List {
	topology := generator.Generate(numServices, numConnections, numNamespaces, numRandomConnections)
	list := corev1.List{
		TypeMeta: metav1.TypeMeta{
			Kind:       "List",
			APIVersion: "v1",
		},
		ListMeta: metav1.ListMeta{},
	}

	for namespace, services := range topology {
		tg := true

		ns := NamespaceForMimik(namespace)
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

func NamespaceForMimik(namespace string) *corev1.Namespace {
	return &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   namespace,
			Labels: map[string]string{"istio-injection": "enabled", "generated-by": "mimik"},
		},
	}
}

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

func DeploymentForMimik(s api.Service, namespace string, tg bool) *appsv1.Deployment {
	labels := labelsForMimik(s.Name, s.Version)
	replicas := int32(1)
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
					Annotations: map[string]string{"sidecar.istio.io/inject": "true"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: "quay.io/leandroberetta/mimik:v0.0.2",
						Name:  "mimik",
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
