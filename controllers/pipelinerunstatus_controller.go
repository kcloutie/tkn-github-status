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

package controllers

import (
	"context"
	"fmt"

	statusv1 "github.com/kcloutie/tkn-github-status/api/v1"
	pipelineBeta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	controllerLabelPrefix = "kcloutie.com"
	kubeConfigEnvVar      = "KUBECONFIG"
)

var controllerConfigLabel = fmt.Sprintf("%s/tknGithubStatus", controllerLabelPrefix)

// PipelineRunStatusReconciler reconciles a PipelineRunStatus object
type PipelineRunStatusReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=status.kcloutie.com,resources=pipelinerunstatuses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=status.kcloutie.com,resources=pipelinerunstatuses/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=status.kcloutie.com,resources=pipelinerunstatuses/finalizers,verbs=update
//+kubebuilder:rbac:groups="tekton.dev",resources=pipelineruns,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PipelineRunStatus object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *PipelineRunStatusReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	log.V(3).Info("Reconciling TektonHelperConfig")

	var status statusv1.PipelineRunStatus
	if err := r.Get(ctx, req.NamespacedName, &status); err != nil {
		log.Error(err, "unable to get PipelineRunStatus", "namespace", req.Namespace, "name", req.Name)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if status.Spec.Disabled {
		return ctrl.Result{}, nil
	}

	foundPipelineRuns, err := r.listPipelineRuns(ctx, req.Namespace, status)
	if err != nil {
		return ctrl.Result{}, nil
	}

	if status.Status.PipelineRuns == nil {
		status.Status.PipelineRuns = map[string]statusv1.PipelineRunsStatus{}
	}

	foundPipelineRunNames := map[string]pipelineBeta1.PipelineRun{}
	// executeStatusUpdate := false

	for _, item := range foundPipelineRuns {
		shouldProcess := true

		foundPipelineRunNames[string(item.UID)] = item

		_, ok := status.Status.PipelineRuns[string(item.UID)]

		if ok {
			// Need to determine the status of the pipelineRun and only update the status if it is different than what we had set previously
		} else {
			// We have not yest set the status in github, lets set it
			shouldProcess = true
		}

		if shouldProcess {

		}
	}

	// toRemove := []string{}
	// for key := range status.Status.Processed {
	// 	_, ok := foundPipelineRunNames[key]
	// 	if !ok {
	// 		toRemove = append(toRemove, key)
	// 	}
	// }

	// for _, item := range toRemove {
	// 	delete(status.Status.Processed, item)
	// 	executeStatusUpdate = true
	// }
	// if executeStatusUpdate {
	// 	log.V(1).Info("Executing Status Update")
	// 	if err := r.Status().Update(ctx, &config); err != nil {
	// 		return ctrl.Result{}, fmt.Errorf("unable to update TektonHelperConfig status: %w", err)
	// 	}
	// }

	return ctrl.Result{}, nil
}

func (r *PipelineRunStatusReconciler) listPipelineRuns(ctx context.Context, namespace string, status statusv1.PipelineRunStatus) (map[string]pipelineBeta1.PipelineRun, error) {
	log := log.FromContext(ctx)

	configNameLabelReq, _ := labels.NewRequirement(controllerConfigLabel, selection.Equals, []string{status.Name})
	selector := labels.NewSelector()
	selector = selector.Add(*configNameLabelReq)

	listOps := &client.ListOptions{
		LabelSelector: selector,
		Namespace:     namespace,
	}

	foundPipelineRuns := map[string]pipelineBeta1.PipelineRun{}
	pipelineRuns := &pipelineBeta1.PipelineRunList{}
	err := r.List(ctx, pipelineRuns, listOps)
	if err != nil {
		log.Error(err, "failed to get the list of pipelineRuns", "namespace", namespace, "pipelineRunStatusName", status.Name)
		return foundPipelineRuns, err
	}

	for _, item := range pipelineRuns.Items {
		log.V(1).Info("Found a pipelineRun with the correct label value", "namespace", namespace, "pipelineRunStatusName", status.Name, "PipelineRun", item.Name, "label", controllerConfigLabel)
		foundPipelineRuns[string(item.UID)] = item
	}
	return foundPipelineRuns, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PipelineRunStatusReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&statusv1.PipelineRunStatus{}).
		Complete(r)
}
