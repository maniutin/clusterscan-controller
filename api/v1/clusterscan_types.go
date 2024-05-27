package v1

import (
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClusterScanSpec defines the desired state of ClusterScan
type ClusterScanSpec struct {
	// JobTemplate is a template for creating Jobs
	JobTemplate corev1.PodTemplateSpec `json:"jobTemplate,omitempty"`

	// CronJobTemplate is a template for creating CronJobs
	CronJobTemplate batchv1beta1.CronJobSpec `json:"cronJobTemplate,omitempty"`
}

// ClusterScanStatus defines the observed state of ClusterScan
type ClusterScanStatus struct {
	// JobStatus is the current status of the created Job
	JobStatus batchv1.JobStatus `json:"jobStatus,omitempty"`

	// CronJobStatus is the current status of the created CronJob
	CronJobStatus batchv1beta1.CronJobStatus `json:"cronJobStatus,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ClusterScan is the Schema for the clusterscans API
type ClusterScan struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterScanSpec   `json:"spec,omitempty"`
	Status ClusterScanStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ClusterScanList contains a list of ClusterScan
type ClusterScanList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterScan `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterScan{}, &ClusterScanList{})
}
