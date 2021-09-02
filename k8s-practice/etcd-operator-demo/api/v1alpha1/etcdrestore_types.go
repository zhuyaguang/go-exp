/*
Copyright 2021 cnych.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type EtcdRestorePhase string

var (
	EtcdRestorePhasePending   EtcdRestorePhase = "Pending"
	EtcdRestorePhaseFailed    EtcdRestorePhase = "Failed"
	EtcdRestorePhaseCompleted EtcdRestorePhase = "Completed"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// EtcdRestoreSpec defines the desired state of EtcdRestore
type EtcdRestoreSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	EtcdCluster EtcdClusterRef `json:"etcdCluster"`

	BackupStorageType BackupStorageType `json:"backupStorageType"`

	// restore data source
	RestoreSource `json:",inline"`
}

type EtcdClusterRef struct {
	// Name is the EtcdCluster resource name.
	Name string `json:"name"`
}

type RestoreSource struct {
	S3  *S3BackupSource  `json:"s3,omitempty"`
	OSS *OSSBackupSource `json:"oss,omitempty"`
}

// EtcdRestoreStatus defines the observed state of EtcdRestore
type EtcdRestoreStatus struct {
	Phase EtcdRestorePhase `json:"phase,omitempty"`
	// Reason indicates the reason for any restore related failures.
	Reason string `json:"reason,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// EtcdRestore is the Schema for the etcdrestores API
type EtcdRestore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EtcdRestoreSpec   `json:"spec,omitempty"`
	Status EtcdRestoreStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// EtcdRestoreList contains a list of EtcdRestore
type EtcdRestoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EtcdRestore `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EtcdRestore{}, &EtcdRestoreList{})
}
