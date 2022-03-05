/*
Copyright 2021.

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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MultiClusterPolicySpec defines the desired state of MultiClusterPolicy
type MultiClusterPolicySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Selects the pods to which this MultiClusterPolicy object applies.
	// This field is NOT optional and follows standard
	// label selector semantics. An empty podSelector matches all pods in this
	// namespace.
	PodSelector metav1.LabelSelector `json:"podSelector" protobuf:"bytes,1,opt,name=podSelector"`

	// List of egress rules to be applied to the selected pods.
	// +optional
	Egress []MultiClusterPolicyEgressRule `json:"egress,omitempty" protobuf:"bytes,3,rep,name=egress"`
}

// MultiClusterPolicyEgressRule  describes a particular set of traffic that is allowed out of pods
// matched by a MultiClusterPolicySpec's podSelector
type MultiClusterPolicyEgressRule struct {

	// List of destination MultiCluster Service Imports for outgoing traffic
	// of pods selected for this rule.
	// Items in this list are combined using a logical OR operation. If this field is
	// empty or missing, this rule matches all destinations (traffic not restricted by
	// destination). If this field is present and contains at least one item, this rule
	// allows traffic only if the traffic matches at least one item in the to list.
	// +optional
	To []MultiClusterPolicyPeer `json:"to,omitempty" protobuf:"bytes,2,rep,name=to"`
}

// MultiClusterPolicyPeer describes a peer to allow traffic to.
type MultiClusterPolicyPeer struct {

	// Selects Namespaces using cluster-scoped labels. This field
	// must be present
	NamespaceSelector *metav1.LabelSelector `json:"namespaceSelector,omitempty" protobuf:"bytes,2,opt,name=namespaceSelector"`

	//List of ServiceImports
	ServiceImportRefs []string `json:"serviceImportRefs,omitempty" protobuf:"bytes,2,opt,name=serviceImportRefs"`
}

// MultiClusterPolicyStatus defines the observed state of MultiClusterPolicy
type MultiClusterPolicyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Boolbean flag indicates if policy is valid and applied to the data plane
	Valid bool `json:"valid,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MultiClusterPolicy is the Schema for the multiclusterpolicies API
type MultiClusterPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MultiClusterPolicySpec   `json:"spec,omitempty"`
	Status MultiClusterPolicyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MultiClusterPolicyList contains a list of MultiClusterPolicy
type MultiClusterPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MultiClusterPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MultiClusterPolicy{}, &MultiClusterPolicyList{})
}
