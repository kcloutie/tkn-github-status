/*
Copyright 2023 Ken Cloutier.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PipelineRunVariableSpec struct {

	// The name of the pipelineRun variable that contains the github organization/owner name
	OwnerVariable string `json:"ownerVariable,omitempty"`

	// The name of the pipelineRun variable that contains the github repository name
	RepoVariable string `json:"repoVariable,omitempty"`

	// The name of the pipelineRun variable that contains the commit sha
	ShaVariable string `json:"shaVariable,omitempty"`

	// The name of the pipelineRun variable that contains the commit sha
	// TargetUrlVariable string `json:"targetUrlVariable,omitempty"`

}

// SecretRef contains the information required to reference a single secret string
// This is needed because the other secretRef types are not cross-namespace and do not
// actually contain the "SecretName" field, which allows us to access a single secret value.
type SecretRef struct {
	SecretKey  string `json:"secretKey,omitempty"`
	SecretName string `json:"secretName,omitempty"`
}

// PipelineRunStatusSpec defines the desired state of PipelineRunStatus
// Link: https://docs.github.com/en/rest/commits/statuses?apiVersion=2022-11-28#create-a-commit-status
type PipelineRunStatusSpec struct {
	// True or false whether the operator is disabled
	// +optional
	Disabled bool `json:"disabled,omitempty"`

	// A string label to differentiate this status from the status of other systems. This field is case-insensitive.
	// +optional
	Context string `json:"context,omitempty"`

	// The base target URL to associate with this status. This URL will be linked from the GitHub UI to allow users to easily see the source of the status.
	// +optional
	TargetUrlBaseUri string `json:"targetUrlBaseUri,omitempty"`

	// Pipeline variables extracted from the pipelineRun param section used to determine the owner/org name, repo name and commit sha of the pipelineRun
	PipelineRunVariables PipelineRunVariableSpec `json:"pipelineRunVariables,omitempty"`

	// The secret that contains the github personal access token used to make api calls to github
	SecretRef *SecretRef `json:"secretRef,omitempty"`

	GithubEnterprise string `json:"githubEnterprise,omitempty"`
}

// PipelineRunStatusStatus defines the observed state of PipelineRunStatus
type PipelineRunStatusStatus struct {
	PipelineRuns map[string]PipelineRunsStatus `json:"pipelineRuns,omitempty"`
}

type PipelineRunsStatus struct {
	// The name of the pipeline run
	PipelineRunName string `json:"pipelineRunName,omitempty"`

	// Wether an error occurred processing the PipelineRun
	HadError string `json:"hadError,omitempty"`

	// The result of processing the PipelineRun
	Result string `json:"result,omitempty"`

	// The last record status of the pipelineRun
	// +optional
	LastStatus CommitStatusesStatus `json:"lastStatus,omitempty"`
}
type CommitStatusesStatus struct {
	// The account owner of the repository
	Owner string `json:"owner,omitempty"`

	// The name of the repository
	Repo string `json:"repo,omitempty"`

	// The commit sha
	Sha string `json:"sha,omitempty"`

	// The state of the status
	State string `json:"state,omitempty"`

	// The target URL to associate with this status. This URL will be linked from the GitHub UI to allow users to easily see the source of the status
	TargetUrl string `json:"targetUrl,omitempty"`

	// A short description of the status
	Description string `json:"description,omitempty"`

	// A string label to differentiate this status from the status of other systems
	Context string `json:"context,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// PipelineRunStatus is the Schema for the pipelinerunstatuses API
type PipelineRunStatus struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PipelineRunStatusSpec   `json:"spec,omitempty"`
	Status PipelineRunStatusStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PipelineRunStatusList contains a list of PipelineRunStatus
type PipelineRunStatusList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PipelineRunStatus `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PipelineRunStatus{}, &PipelineRunStatusList{})
}
