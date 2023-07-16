/*
Copyright 2022.

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

	common "github.com/gigiozzz/depiy/common-libs/commons"
	"github.com/gigiozzz/depiy/operators/plugin-operator/api/v1alpha1"
	pluginv1alpha1 "github.com/gigiozzz/depiy/operators/plugin-operator/api/v1alpha1"
	"github.com/gigiozzz/depiy/operators/plugin-operator/controllers/reconcilers"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

const (
	entandoPluginFinalizer = "plugin.entando.org/finalizer"
	controllerLogName      = "EntandoPluginV2 Controller"
)

// EntandoPluginV2Reconciler reconciles a EntandoPluginV2 object
type EntandoPluginV2Reconciler struct {
	Base     common.BaseK8sStructure
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=plugin.entando.org,resources=entandopluginv2s,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=plugin.entando.org,resources=entandopluginv2s/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=plugin.entando.org,resources=entandopluginv2s/finalizers,verbs=update
//+kubebuilder:rbac:groups=gateway.entando.org,resources=entandogatewayv2s,verbs=get;list;watch;create;update;patch;delete
// Annotation for generating RBAC role for writing Events
//+kubebuilder:rbac:groups="*",resources=events,verbs=create;patch
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="*",resources=services,verbs=get;list;watch;create;update;patch;delete

func NewEntandoPluginV2Reconciler(client client.Client, log logr.Logger, scheme *runtime.Scheme, recorder record.EventRecorder) *EntandoPluginV2Reconciler {
	return &EntandoPluginV2Reconciler{
		Base:     common.BaseK8sStructure{Client: client, Log: log},
		Scheme:   scheme,
		Recorder: recorder,
	}
}

func (r *EntandoPluginV2Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	//	_ = log.FromContext(ctx)
	log := r.Base.Log.WithName(controllerLogName)
	log.Info("Start reconciling EntandoPluginV2 custom resources")

	cr := &v1alpha1.EntandoPluginV2{}
	err := r.Base.Get(ctx, req.NamespacedName, cr)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Check if the EntandoApp instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	isEntandoAppV2MarkedToBeDeleted := cr.GetDeletionTimestamp() != nil
	if isEntandoAppV2MarkedToBeDeleted {
		if err := r.removeFinalizer(ctx, cr, log); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// Add finalizer for this CR
	err = r.addFinalizer(ctx, cr)
	if err != nil {
		return ctrl.Result{}, err
	}

	recoManager := reconcilers.NewReconcileManager(r.Base.Client, r.Base.Log, r.Scheme, r.Recorder)
	res, err := recoManager.MainReconcile(ctx, req, cr)

	log.Info("Reconciled EntandoPluginV2 custom resources")
	return res, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *EntandoPluginV2Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&pluginv1alpha1.EntandoPluginV2{}).
		//Owns(&appsv1.Deployment{}).
		WithEventFilter(predicate.GenerationChangedPredicate{}). //solo modifiche a spec
		Complete(r)
}

// =====================================================================
// Add the cleanup steps that the operator
// needs to do before the CR can be deleted. Examples
// of finalizers include performing backups and deleting
// resources that are not owned by this CR, like a PVC.
// =====================================================================
func (r *EntandoPluginV2Reconciler) finalizeEntandoApp(log logr.Logger, m *v1alpha1.EntandoPluginV2) error {
	log.Info("Successfully finalized entandoApp")
	return nil
}

func (r *EntandoPluginV2Reconciler) addFinalizer(ctx context.Context, cr *v1alpha1.EntandoPluginV2) error {
	if !controllerutil.ContainsFinalizer(cr, entandoPluginFinalizer) {
		controllerutil.AddFinalizer(cr, entandoPluginFinalizer)
		return r.Base.Update(ctx, cr)
	}
	return nil
}

func (r *EntandoPluginV2Reconciler) removeFinalizer(ctx context.Context, cr *v1alpha1.EntandoPluginV2, log logr.Logger) error {
	if controllerutil.ContainsFinalizer(cr, entandoPluginFinalizer) {
		// Run finalization logic for entandoAppFinalizer. If the
		// finalization logic fails, don't remove the finalizer so
		// that we can retry during the next reconciliation.
		if err := r.finalizeEntandoApp(log, cr); err != nil {
			return err
		}

		// Remove entandoAppFinalizer. Once all finalizers have been
		// removed, the object will be deleted.
		controllerutil.RemoveFinalizer(cr, entandoPluginFinalizer)
		err := r.Base.Update(ctx, cr)
		if err != nil {
			return err
		}
	}
	return nil
}
