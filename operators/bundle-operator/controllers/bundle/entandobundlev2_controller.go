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

package bundle

import (
	"sigs.k8s.io/controller-runtime/pkg/client"

	"context"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	common "github.com/gigiozzz/depiy/common-libs/commons"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	bundlev1alpha1 "github.com/gigiozzz/depiy/operators/bundle-operator/api/v1alpha1"
)

const (
	entandoBundleFinalizer = "bundle.entando.org/finalizer"
	controllerLogName      = "EntandoBundleV2 Controller"
)

// EntandoBundleV2Reconciler reconciles a EntandoBundleV2 object
type EntandoBundleV2Reconciler struct {
	Base     common.BaseK8sStructure
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=bundle.entando.org,resources=entandobundlev2s,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=bundle.entando.org,resources=entandobundlev2s/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=bundle.entando.org,resources=entandobundlev2s/finalizers,verbs=update

func NewEntandoBundleV2Reconciler(client client.Client, log logr.Logger, scheme *runtime.Scheme, recorder record.EventRecorder) *EntandoBundleV2Reconciler {
	return &EntandoBundleV2Reconciler{
		Base:     common.BaseK8sStructure{Client: client, Log: log},
		Scheme:   scheme,
		Recorder: recorder,
	}
}

func (r *EntandoBundleV2Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Base.Log.WithName(controllerLogName)
	log.Info("Start reconciling EntandoBundleV2 custom resources")

	cr := &bundlev1alpha1.EntandoBundleV2{}
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

	recoBundleManager := NewReconcileBundleManager(r.Base.Client, r.Base.Log, r.Scheme, r.Recorder)
	res, err := recoBundleManager.MainReconcile(ctx, req, cr)

	log.Info("Reconciled EntandoBundleV2 custom resources")
	return res, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *EntandoBundleV2Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&bundlev1alpha1.EntandoBundleV2{}).
		WithEventFilter(predicate.GenerationChangedPredicate{}). //solo modifiche a spec
		Complete(r)
}

// =====================================================================
// Add the cleanup steps that the operator
// needs to do before the CR can be deleted. Examples
// of finalizers include performing backups and deleting
// resources that are not owned by this CR, like a PVC.
// =====================================================================
func (r *EntandoBundleV2Reconciler) finalizeEntandoApp(log logr.Logger, m *bundlev1alpha1.EntandoBundleV2) error {
	log.Info("Successfully finalized entandoApp")
	return nil
}

func (r *EntandoBundleV2Reconciler) addFinalizer(ctx context.Context, cr *bundlev1alpha1.EntandoBundleV2) error {
	if !controllerutil.ContainsFinalizer(cr, entandoBundleFinalizer) {
		controllerutil.AddFinalizer(cr, entandoBundleFinalizer)
		return r.Base.Update(ctx, cr)
	}
	return nil
}

func (r *EntandoBundleV2Reconciler) removeFinalizer(ctx context.Context, cr *bundlev1alpha1.EntandoBundleV2, log logr.Logger) error {
	if controllerutil.ContainsFinalizer(cr, entandoBundleFinalizer) {
		// Run finalization logic for entandoAppFinalizer. If the
		// finalization logic fails, don't remove the finalizer so
		// that we can retry during the next reconciliation.
		if err := r.finalizeEntandoApp(log, cr); err != nil {
			return err
		}

		// Remove entandoAppFinalizer. Once all finalizers have been
		// removed, the object will be deleted.
		controllerutil.RemoveFinalizer(cr, entandoBundleFinalizer)
		err := r.Base.Update(ctx, cr)
		if err != nil {
			return err
		}
	}
	return nil
}
