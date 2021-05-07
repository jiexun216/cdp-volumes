package hook

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/golang/glog"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// cdp-volumes admission is to add volumeMounts
func createDeploymentAddVolumePatch(deployment appsv1.Deployment, availableAnnotations map[string]string, annotations map[string]string) ([]byte, error) {
	var patch []patchOperation
	// update Annotation to set admissionWebhookAnnotationStatusKey: "mutated"
	patch = append(patch, updateAnnotation(availableAnnotations, annotations)...)

	// read configMap to decide modify the sts
	volumeMap := getConfigMap()
	if volumeMap != nil {
		if value, ok := volumeMap["deployment."+deployment.Name]; ok {
			// modify
			var deployTemplate appsv1.Deployment
			if err := json.Unmarshal([]byte(value), &deployTemplate); err != nil {
				glog.Errorf("Can't json.Unmarshal stsTemplate: %v", err)
			}
			// pod level
			if deployTemplate.Spec.Template.Spec.Volumes != nil {
				// modify /spec/template/spec/volumes
				replaceVolumes := patchOperation{
					Op:    "replace",
					Path:  "/spec/template/spec/volumes",
					Value: deployTemplate.Spec.Template.Spec.Volumes,
				}
				glog.Infof("add Deployment Volumes  /spec/template/spec/Volumes for value: %v", replaceVolumes)
				patch = append(patch, replaceVolumes)
			}

			// initContainers level
			if len(deployTemplate.Spec.Template.Spec.InitContainers) > 0 {
				initContainers := deployment.Spec.Template.Spec.InitContainers
				var initSize = len(initContainers)
				for i := 0; i < initSize; i++ {
					for _, val := range deployTemplate.Spec.Template.Spec.InitContainers {
						if val.Name == initContainers[i].Name {
							// only add Volumes initcontainers
							replaceVolumeMounts := patchOperation{
								Op:    "replace",
								Path:  "/spec/template/spec/initContainers/" + strconv.Itoa(i) + "/volumeMounts",
								Value: val.VolumeMounts,
							}
							patch = append(patch, replaceVolumeMounts)
						}
					}
				}
			}

			// containers level
			if len(deployTemplate.Spec.Template.Spec.Containers) > 0 {
				containers := deployment.Spec.Template.Spec.Containers
				var containerSize = len(containers)
				for i := 0; i < containerSize; i++ {
					for _, val := range deployTemplate.Spec.Template.Spec.Containers {
						if val.Name == containers[i].Name {
							// only add volumeMounts containers
							replaceVolumeMounts := patchOperation{
								Op:    "replace",
								Path:  "/spec/template/spec/containers/" + strconv.Itoa(i) + "/volumeMounts",
								Value: val.VolumeMounts,
							}
							patch = append(patch, replaceVolumeMounts)
						}
					}
				}
			}
		}
	}

	return json.Marshal(patch)
}

func createStatefulsetAddVolumePatch(statefulset appsv1.StatefulSet, availableAnnotations map[string]string, annotations map[string]string) ([]byte, error) {
	var patch []patchOperation
	// update Annotation to set admissionWebhookAnnotationStatusKey: "mutated"
	patch = append(patch, updateAnnotation(availableAnnotations, annotations)...)

	// read configMap to decide modify the sts
	volumeMap := getConfigMap()
	if volumeMap != nil {
		for k, value := range volumeMap {
			if strings.Contains("statefulset."+statefulset.Name, k) {
				// modify
				var stsTemplate appsv1.StatefulSet
				if err := json.Unmarshal([]byte(value), &stsTemplate); err != nil {
					glog.Errorf("Can't json.Unmarshal stsTemplate: %v", err)
				}
				// pod level
				if stsTemplate.Spec.Template.Spec.Volumes != nil {
					// modify /spec/template/spec/volumes
					replaceVolumes := patchOperation{
						Op:    "replace",
						Path:  "/spec/template/spec/volumes",
						Value: stsTemplate.Spec.Template.Spec.Volumes,
					}
					glog.Infof("add StatefulSet Volumes  /spec/template/spec/volumes for value: %v", replaceVolumes)
					patch = append(patch, replaceVolumes)
				}

				// initContainers level
				if len(stsTemplate.Spec.Template.Spec.InitContainers) > 0 {
					initContainers := statefulset.Spec.Template.Spec.InitContainers
					var initSize = len(initContainers)
					for i := 0; i < initSize; i++ {
						for _, val := range stsTemplate.Spec.Template.Spec.InitContainers {
							if val.Name == initContainers[i].Name {
								// only add volumeMounts initcontainers
								replaceVolumeMounts := patchOperation{
									Op:    "replace",
									Path:  "/spec/template/spec/initContainers/" + strconv.Itoa(i) + "/volumeMounts",
									Value: val.VolumeMounts,
								}
								patch = append(patch, replaceVolumeMounts)
							}
						}
					}
				}

				// containers level
				if len(stsTemplate.Spec.Template.Spec.Containers) > 0 {
					containers := statefulset.Spec.Template.Spec.Containers
					var containerSize = len(containers)
					for i := 0; i < containerSize; i++ {
						for _, val := range stsTemplate.Spec.Template.Spec.Containers {
							if val.Name == containers[i].Name {
								// only add volumeMounts containers
								replaceVolumeMounts := patchOperation{
									Op:    "replace",
									Path:  "/spec/template/spec/containers/" + strconv.Itoa(i) + "/volumeMounts",
									Value: val.VolumeMounts,
								}
								patch = append(patch, replaceVolumeMounts)
							}
						}
					}
				}
			}

		}

	}

	return json.Marshal(patch)
}

func createJobAddVolumePatch(job batchv1.Job, availableAnnotations map[string]string, annotations map[string]string) ([]byte, error) {
	var patch []patchOperation
	// update Annotation to set admissionWebhookAnnotationStatusKey: "mutated"
	patch = append(patch, updateAnnotation(availableAnnotations, annotations)...)

	// read configMap to decide modify the sts
	volumeMap := getConfigMap()
	if volumeMap != nil {
		for k, value := range volumeMap {
			if strings.Contains("job."+job.Name, k) {
				// modify
				var jobTemplate batchv1.Job
				if err := json.Unmarshal([]byte(value), &jobTemplate); err != nil {
					glog.Errorf("Can't json.Unmarshal stsTemplate: %v", err)
				}
				// pod level
				if jobTemplate.Spec.Template.Spec.Volumes != nil {
					// modify /spec/template/spec/volumes
					replaceVolumes := patchOperation{
						Op:    "replace",
						Path:  "/spec/template/spec/volumes",
						Value: jobTemplate.Spec.Template.Spec.Volumes,
					}
					glog.Infof("add StatefulSet Volumes  /spec/template/spec/volumes for value: %v", replaceVolumes)
					patch = append(patch, replaceVolumes)
				}

				// initContainers level
				if len(jobTemplate.Spec.Template.Spec.InitContainers) > 0 {
					initContainers := job.Spec.Template.Spec.InitContainers
					var initSize = len(initContainers)
					for i := 0; i < initSize; i++ {
						for _, val := range jobTemplate.Spec.Template.Spec.InitContainers {
							if val.Name == initContainers[i].Name {
								// only add VolumeMounts initcontainers
								replaceVolumeMounts := patchOperation{
									Op:    "replace",
									Path:  "/spec/template/spec/initContainers/" + strconv.Itoa(i) + "/volumeMounts",
									Value: val.VolumeMounts,
								}
								patch = append(patch, replaceVolumeMounts)
							}
						}
					}
				}

				// containers level
				if len(jobTemplate.Spec.Template.Spec.Containers) > 0 {
					containers := job.Spec.Template.Spec.Containers
					var containerSize = len(containers)
					for i := 0; i < containerSize; i++ {
						for _, val := range jobTemplate.Spec.Template.Spec.Containers {
							if val.Name == containers[i].Name {
								// only add VolumeMounts containers
								replaceVolumeMounts := patchOperation{
									Op:    "replace",
									Path:  "/spec/template/spec/containers/" + strconv.Itoa(i) + "/volumeMounts",
									Value: val.VolumeMounts,
								}
								patch = append(patch, replaceVolumeMounts)
							}
						}
					}
				}
			}
		}

	}

	return json.Marshal(patch)
}

func updateAnnotation(target map[string]string, added map[string]string) (patch []patchOperation) {
	for key, value := range added {
		if target == nil || target[key] == "" {
			target = map[string]string{}
			patch = append(patch, patchOperation{
				Op:   "add",
				Path: "/metadata/annotations",
				Value: map[string]string{
					key: value,
				},
			})
		} else {
			patch = append(patch, patchOperation{
				Op:    "replace",
				Path:  "/metadata/annotations/" + key,
				Value: value,
			})
		}
	}
	return patch
}

func getConfigMap() map[string]string {
	config, err := rest.InClusterConfig()
	if err != nil {
		glog.Errorf("Can't get ClusterConfig: %v", err)
		return nil
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Errorf("Can't connect to kubernetes: %v", err)
		return nil
	}

	configMapClient := clientset.CoreV1().ConfigMaps("cdp-customizer")
	// labelSelector cdp.cloudera.io/security-context
	labelSelector := metav1.LabelSelector{
		//MatchLabels: map[string]string{"cdp.cloudera.io/security-context":"true"},
		MatchExpressions: []metav1.LabelSelectorRequirement{
			{
				Key:      "cdp.cloudera.io/security-context",
				Operator: "Exists",
			},
		},
	}
	listOptions := metav1.ListOptions{
		LabelSelector: labels.Set(labelSelector.MatchLabels).String(),
	}
	configMapList, err := configMapClient.List(context.Background(), listOptions)
	if err != nil {
		glog.Errorf("Can't get the specific configMap: %v", err)
		return nil
	}
	mergeResult := make(map[string]string)
	for _, configMap := range configMapList.Items {
		for k, v := range configMap.Data {
			mergeResult[k] = v
		}
	}
	if len(mergeResult) == 0 {
		glog.Errorf("get the specific configMap,but is empty")
		return nil
	}
	return mergeResult

	/*configMap, err := configMapClient.Get(context.Background(), "volumes-patches", metav1.GetOptions{})
	if err != nil  {
		glog.Errorf("Can't get the specific configMap: %v", err)
		return nil
	}
	glog.Infof("get the specific configMap data value: %v", configMap.Data)

	return configMap.Data*/

}
