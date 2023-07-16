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

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	common "github.com/gigiozzz/depiy/common-libs/commons"
	v1alpha1 "github.com/gigiozzz/depiy/operators/gateway-operator/api/v1alpha1"
	"github.com/gigiozzz/depiy/operators/gateway-operator/controllers/reconcilers"
	"github.com/go-logr/logr"
)

const (
	entandoGatewayFinalizer  = "gateway.entando.org/finalizer"
	controllerIngressLogName = "EntandoGatewayV2 Controller"
)

// EntandoGatewayV2Reconciler reconciles a EntandoGatewayV2 object
type EntandoGatewayV2Reconciler struct {
	Base     common.BaseK8sStructure
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=gateway.entando.org,resources=entandogatewayv2s,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=gateway.entando.org,resources=entandogatewayv2s/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=gateway.entando.org,resources=entandogatewayv2s/finalizers,verbs=update
// Annotation for generating RBAC role for writing Events
//+kubebuilder:rbac:groups="*",resources=events,verbs=create;patch
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete

func NewEntandoGatewayV2Reconciler(client client.Client, log logr.Logger, scheme *runtime.Scheme, recorder record.EventRecorder) *EntandoGatewayV2Reconciler {
	return &EntandoGatewayV2Reconciler{
		Base:     common.BaseK8sStructure{Client: client, Log: log},
		Scheme:   scheme,
		Recorder: recorder,
	}
}

func (r *EntandoGatewayV2Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	//_ = log.FromContext(ctx)
	log := r.Base.Log.WithName(controllerIngressLogName)
	log.Info("Start reconciling EntandoGatewayV2 custom resources")

	cr := &v1alpha1.EntandoGatewayV2{}
	err := r.Base.Get(ctx, req.NamespacedName, cr)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Check if the EntandoApp instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	isEntandoIngressV2MarkedToBeDeleted := cr.GetDeletionTimestamp() != nil
	if isEntandoIngressV2MarkedToBeDeleted {
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

	log.Info("Reconciled EntandoGatewayV2 custom resources")
	return res, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *EntandoGatewayV2Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.EntandoGatewayV2{}).
		WithEventFilter(predicate.GenerationChangedPredicate{}). //solo modifiche a spec
		Complete(r)
}

// =====================================================================
// Add the cleanup steps that the operator
// needs to do before the CR can be deleted. Examples
// of finalizers include performing backups and deleting
// resources that are not owned by this CR, like a PVC.
// =====================================================================
func (r *EntandoGatewayV2Reconciler) finalizeEntandoApp(log logr.Logger, m *v1alpha1.EntandoGatewayV2) error {
	log.Info("Successfully finalized entandoApp")
	return nil
}

func (r *EntandoGatewayV2Reconciler) addFinalizer(ctx context.Context, cr *v1alpha1.EntandoGatewayV2) error {
	if !controllerutil.ContainsFinalizer(cr, entandoGatewayFinalizer) {
		controllerutil.AddFinalizer(cr, entandoGatewayFinalizer)
		return r.Base.Update(ctx, cr)
	}
	return nil
}

func (r *EntandoGatewayV2Reconciler) removeFinalizer(ctx context.Context, cr *v1alpha1.EntandoGatewayV2, log logr.Logger) error {
	if controllerutil.ContainsFinalizer(cr, entandoGatewayFinalizer) {
		// Run finalization logic for entandoAppFinalizer. If the
		// finalization logic fails, don't remove the finalizer so
		// that we can retry during the next reconciliation.
		if err := r.finalizeEntandoApp(log, cr); err != nil {
			return err
		}

		// Remove entandoAppFinalizer. Once all finalizers have been
		// removed, the object will be deleted.
		controllerutil.RemoveFinalizer(cr, entandoGatewayFinalizer)
		err := r.Base.Update(ctx, cr)
		if err != nil {
			return err
		}
	}
	return nil
}
