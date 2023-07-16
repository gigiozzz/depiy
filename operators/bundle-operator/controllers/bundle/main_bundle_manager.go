package bundle

import (
	"context"

	common "github.com/gigiozzz/depiy/common-libs/commons"
	"github.com/gigiozzz/depiy/operators/bundle-operator/api/v1alpha1"
	"github.com/gigiozzz/depiy/operators/bundle-operator/controllers/services"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ReconcileBundleManager struct {
	Base      *common.BaseK8sStructure
	Scheme    *runtime.Scheme
	Recorder  record.EventRecorder
	Condition *services.ConditionService
}

func NewReconcileBundleManager(client client.Client, log logr.Logger, scheme *runtime.Scheme, recorder record.EventRecorder) *ReconcileBundleManager {
	base := &common.BaseK8sStructure{Client: client, Log: log}
	return &ReconcileBundleManager{
		Base:      base,
		Scheme:    scheme,
		Recorder:  recorder,
		Condition: services.NewConditionService(base),
	}
}

func (r *ReconcileBundleManager) MainReconcile(ctx context.Context, req ctrl.Request, cr *v1alpha1.EntandoBundleV2) (ctrl.Result, error) {

	log := r.Base.Log
	bundleService := services.NewBundleService()

	if err := r.Condition.SetConditionBundleReadyUnknow(ctx, cr); err != nil {
		log.Info("error on set instance ready unknow")
		return ctrl.Result{}, err
	}

	// generate bundleCode
	if err := r.generateAndSaveBundleCode(ctx, cr, bundleService); err != nil {
		log.Info("error generateAndSaveBundleCode reschedule reconcile", "error", err)
		r.Condition.SetConditionBundleReadyFalse(ctx, cr)
		return ctrl.Result{}, err
	}

	// verify signature
	if err := r.verifyBundleSignatures(ctx, cr, bundleService); err != nil {
		log.Info("error verifyBundleSignatures reschedule reconcile", "error", err)
		r.Condition.SetConditionBundleReadyFalse(ctx, cr)
		return ctrl.Result{}, err
	}

	r.Condition.SetConditionBundleReadyTrue(ctx, cr)
	return ctrl.Result{}, nil
}

func (r *ReconcileBundleManager) generateAndSaveBundleCode(ctx context.Context,
	cr *v1alpha1.EntandoBundleV2,
	bundleService *services.BundleService) error {
	r.Base.Log.Info("generate code for bundle from repo: " + cr.Spec.Repository)
	bundleCode := bundleService.GenerateBundleCode(cr)
	r.Base.Log.Info("generated bundle code: " + bundleCode)

	annotations := cr.GetAnnotations()
	annotations["bundleCode"] = bundleCode
	cr.SetAnnotations(annotations)

	err := r.Base.Client.Update(ctx, cr)

	return err
}

func (r *ReconcileBundleManager) verifyBundleSignatures(ctx context.Context,
	cr *v1alpha1.EntandoBundleV2,
	bundleService *services.BundleService) error {
	verifiedList, err := bundleService.CheckBundleSignature(ctx, cr, r.Base.Log)
	if err == nil {
		annotations := cr.GetAnnotations()
		for k, v := range verifiedList {
			annotations[k] = v
		}
		cr.SetAnnotations(annotations)
		errSave := r.Base.Client.Update(ctx, cr)
		if errSave != nil {
			r.Base.Log.Error(errSave, "error saving verified sign")
		}
	}
	return nil
}
