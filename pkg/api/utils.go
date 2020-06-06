package api

import (
	"path/filepath"

	"github.com/alecjacobs5401/kube-client-wrapper/pkg/types"
	corev1 "k8s.io/api/core/v1"
)

// PodIsReady returns true if all Pod Conditions are marked as True and if all Containers are Ready
func PodIsReady(pod corev1.Pod) bool {
	return len(FailedPodConditions(pod)) == 0 && len(NotReadyPodContainerStatus(pod)) == 0
}

// FailedPodConditions returns an array of PodConditions that have failed
func FailedPodConditions(pod corev1.Pod) []types.PodCondition {
	conditions := []types.PodCondition{}
	for _, condition := range pod.Status.Conditions {
		if condition.Status != corev1.ConditionTrue {
			conditions = append(
				conditions,
				types.PodCondition{
					Type:       string(condition.Type),
					Successful: false,
					Reason:     condition.Reason,
					Message:    condition.Message,
				},
			)

		}
	}

	return conditions
}

// NotReadyPodContainerStatus returns an array of ContainerStatus for Pod containers that are not ready
func NotReadyPodContainerStatus(pod corev1.Pod) []types.ContainerStatus {
	statuses := []types.ContainerStatus{}
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if !containerStatus.Ready {
			statuses = append(
				statuses,
				types.ContainerStatus{
					Name:  containerStatus.Name,
					Ready: false,
				},
			)

		}
	}

	return statuses
}

// ConfigPathFromDirectory determines the kube config location from the HOME environment variable. If HOME is not defined, return empty.
func ConfigPathFromDirectory(d string) string {
	if d != "" {
		return filepath.Join(d, ".kube", "config")
	}
	return ""
}
