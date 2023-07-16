package services

import (
	"context"
	"fmt"

	common "github.com/gigiozzz/depiy/common-libs/commons"
	utility "github.com/gigiozzz/depiy/common-libs/utilities"
	"github.com/gigiozzz/depiy/operators/bundle-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	// BundleInstance CR condition
	CONDITION_MANIFEST_APPLIED        = "ManifestApplied"
	CONDITION_MANIFEST_APPLIED_REASON = "ManifestIsApplied"
	CONDITION_MANIFEST_APPLIED_MSG    = "Your Manifest was applied"

	CONDITION_PLUGIN_CR_APPLIED        = "PluginCrApplied"
	CONDITION_PLUGIN_CR_APPLIED_REASON = "PluginCrIsApplied"
	CONDITION_PLUGIN_CR_APPLIED_MSG    = "Your Plugin cr was applied"

	CONDITION_PLUGIN_CR_READY        = "PluginCrReady"
	CONDITION_PLUGIN_CR_READY_REASON = "PluginCrIsReady"
	CONDITION_PLUGIN_CR_READY_MSG    = "Your Plugin cr is ready"

	CONDITION_INSTANCE_READY        = "InstanceReady"
	CONDITION_INSTANCE_READY_REASON = "InstanceIsReady"
	CONDITION_INSTANCE_READY_MSG    = "Your Instance is ready"

	// Bundle CR condition
	CONDITION_INSTANCE_CR_APPLIED        = "InstanceCrApplied"
	CONDITION_INSTANCE_CR_APPLIED_REASON = "InstanceCrIsApplied"
	CONDITION_INSTANCE_CR_APPLIED_MSG    = "Your instance cr was applied"

	CONDITION_INSTANCE_CR_READY        = "InstanceCrReady"
	CONDITION_INSTANCE_CR_READY_REASON = "InstanceCrIsReady"
	CONDITION_INSTANCE_CR_READY_MSG    = "Your instance cr is ready"

	CONDITION_BUNDLE_READY        = "BundleReady"
	CONDITION_BUNDLE_READY_REASON = "BundleIsReady"
	CONDITION_BUNDLE_READY_MSG    = "Your Bundle is ready"
)

type ConditionService struct {
	Base *common.BaseK8sStructure
}

func NewConditionService(base *common.BaseK8sStructure) *ConditionService {
	return &ConditionService{
		Base: base,
	}
}

func (cs *ConditionService) IsPluginCrReady(ctx context.Context, cr *v1alpha1.EntandoBundleInstanceV2) bool {

	condition, observedGeneration := cs.getConditionStatus(ctx, cr, CONDITION_PLUGIN_CR_READY)

	return metav1.ConditionTrue == condition && observedGeneration == cr.Generation
}

func (cs *ConditionService) SetConditionPluginCrReady(ctx context.Context, cr *v1alpha1.EntandoBundleInstanceV2) error {

	cs.deleteCondition(ctx, cr, CONDITION_PLUGIN_CR_READY)
	return utility.AppendCondition(ctx, cs.Base.Client, cr,
		CONDITION_PLUGIN_CR_READY,
		metav1.ConditionTrue,
		CONDITION_PLUGIN_CR_READY_REASON,
		CONDITION_PLUGIN_CR_READY_MSG,
		cr.Generation)
}

func (cs *ConditionService) IsManifestApplied(ctx context.Context, cr *v1alpha1.EntandoBundleInstanceV2, manifestId string) bool {

	condition, observedGeneration := cs.getConditionStatus(ctx, cr, CONDITION_MANIFEST_APPLIED+"-"+manifestId)

	return metav1.ConditionTrue == condition && observedGeneration == cr.Generation
}

func (cs *ConditionService) SetConditionManifestApplied(ctx context.Context, cr *v1alpha1.EntandoBundleInstanceV2,
	manifestId string, manifestPath string) error {

	cs.deleteCondition(ctx, cr, CONDITION_MANIFEST_APPLIED+"-"+manifestId)
	return utility.AppendCondition(ctx, cs.Base.Client, cr,
		CONDITION_MANIFEST_APPLIED+"-"+manifestId,
		metav1.ConditionTrue,
		CONDITION_MANIFEST_APPLIED_REASON,
		CONDITION_MANIFEST_APPLIED_MSG+" "+manifestPath,
		cr.Generation)
}

func (cs *ConditionService) IsPluginCrApplied(ctx context.Context, cr *v1alpha1.EntandoBundleInstanceV2, pluginCode string) bool {

	condition, observedGeneration := cs.getConditionStatus(ctx, cr, CONDITION_PLUGIN_CR_APPLIED+"-"+pluginCode)

	return metav1.ConditionTrue == condition && observedGeneration == cr.Generation
}

func (cs *ConditionService) SetConditionPluginCrApplied(ctx context.Context, cr *v1alpha1.EntandoBundleInstanceV2, pluginCode string) error {

	cs.deleteCondition(ctx, cr, CONDITION_PLUGIN_CR_APPLIED+"-"+pluginCode)
	return utility.AppendCondition(ctx, cs.Base.Client, cr,
		CONDITION_PLUGIN_CR_APPLIED+"-"+pluginCode,
		metav1.ConditionTrue,
		CONDITION_PLUGIN_CR_APPLIED_REASON,
		CONDITION_PLUGIN_CR_APPLIED_MSG,
		cr.Generation)
}

func (cs *ConditionService) IsInstanceCrReady(ctx context.Context, cr *v1alpha1.EntandoBundleInstanceV2) bool {

	condition, observedGeneration := cs.getConditionStatus(ctx, cr, CONDITION_INSTANCE_CR_READY)

	return metav1.ConditionTrue == condition && observedGeneration == cr.Generation
}

func (cs *ConditionService) SetConditionInstanceCrReady(ctx context.Context, cr *v1alpha1.EntandoBundleInstanceV2) error {

	cs.deleteCondition(ctx, cr, CONDITION_INSTANCE_CR_READY)
	return utility.AppendCondition(ctx, cs.Base.Client, cr,
		CONDITION_INSTANCE_CR_READY,
		metav1.ConditionTrue,
		CONDITION_INSTANCE_CR_READY_REASON,
		CONDITION_INSTANCE_CR_READY_MSG,
		cr.Generation)
}

func (cs *ConditionService) IsInstanceCrApplied(ctx context.Context, cr *v1alpha1.EntandoBundleInstanceV2) bool {

	condition, observedGeneration := cs.getConditionStatus(ctx, cr, CONDITION_INSTANCE_CR_APPLIED)

	return metav1.ConditionTrue == condition && observedGeneration == cr.Generation
}

func (cs *ConditionService) SetConditionInstanceCrApplied(ctx context.Context, cr *v1alpha1.EntandoBundleInstanceV2) error {

	cs.deleteCondition(ctx, cr, CONDITION_INSTANCE_CR_APPLIED)
	return utility.AppendCondition(ctx, cs.Base.Client, cr,
		CONDITION_INSTANCE_CR_APPLIED,
		metav1.ConditionTrue,
		CONDITION_INSTANCE_CR_APPLIED_REASON,
		CONDITION_INSTANCE_CR_APPLIED_MSG,
		cr.Generation)
}

func (cs *ConditionService) SetConditionInstanceReadyTrue(ctx context.Context, cr *v1alpha1.EntandoBundleInstanceV2) error {
	return cs.setConditionInstanceReady(ctx, cr, metav1.ConditionTrue)
}

func (cs *ConditionService) SetConditionInstanceReadyUnknow(ctx context.Context, cr *v1alpha1.EntandoBundleInstanceV2) error {
	return cs.setConditionInstanceReady(ctx, cr, metav1.ConditionUnknown)
}

func (cs *ConditionService) SetConditionInstanceReadyFalse(ctx context.Context, cr *v1alpha1.EntandoBundleInstanceV2) error {
	return cs.setConditionInstanceReady(ctx, cr, metav1.ConditionFalse)
}

func (cs *ConditionService) setConditionInstanceReady(ctx context.Context, cr *v1alpha1.EntandoBundleInstanceV2, status metav1.ConditionStatus) error {

	cs.deleteCondition(ctx, cr, CONDITION_INSTANCE_READY)
	return utility.AppendCondition(ctx, cs.Base.Client, cr,
		CONDITION_INSTANCE_READY,
		status,
		CONDITION_INSTANCE_READY_REASON,
		CONDITION_INSTANCE_READY_MSG,
		cr.Generation)
}

func (cs *ConditionService) SetConditionBundleReadyTrue(ctx context.Context, cr *v1alpha1.EntandoBundleV2) error {
	return cs.setConditionBundleReady(ctx, cr, metav1.ConditionTrue)
}

func (cs *ConditionService) SetConditionBundleReadyUnknow(ctx context.Context, cr *v1alpha1.EntandoBundleV2) error {
	return cs.setConditionBundleReady(ctx, cr, metav1.ConditionUnknown)
}

func (cs *ConditionService) SetConditionBundleReadyFalse(ctx context.Context, cr *v1alpha1.EntandoBundleV2) error {
	return cs.setConditionBundleReady(ctx, cr, metav1.ConditionFalse)
}

func (cs *ConditionService) setConditionBundleReady(ctx context.Context, cr *v1alpha1.EntandoBundleV2, status metav1.ConditionStatus) error {

	cs.deleteCondition(ctx, cr, CONDITION_BUNDLE_READY)
	return utility.AppendCondition(ctx, cs.Base.Client, cr,
		CONDITION_BUNDLE_READY,
		status,
		CONDITION_BUNDLE_READY_REASON,
		CONDITION_BUNDLE_READY_MSG,
		cr.Generation)
}

func (cs *ConditionService) getConditionStatus(ctx context.Context, cr *v1alpha1.EntandoBundleInstanceV2, typeName string) (metav1.ConditionStatus, int64) {

	var output metav1.ConditionStatus = metav1.ConditionUnknown
	var observedGeneration int64

	for _, condition := range cr.Status.Conditions {
		if condition.Type == typeName {
			output = condition.Status
			observedGeneration = condition.ObservedGeneration
		}
	}
	return output, observedGeneration
}

func (cs *ConditionService) deleteCondition(ctx context.Context, cr client.Object, typeName string) error {

	log := log.FromContext(ctx)
	var newConditions = make([]metav1.Condition, 0)
	conditionsAware, conversionSuccessful := (cr).(utility.ConditionsAware)
	if conversionSuccessful {
		for _, condition := range conditionsAware.GetConditions() {
			if condition.Type != typeName {
				newConditions = append(newConditions, condition)
			}
		}
		conditionsAware.SetConditions(newConditions)

		err := cs.Base.Client.Status().Update(ctx, cr)
		if err != nil {
			log.Info("Application resource status update failed.")
		}
		return nil

	} else {
		errMessage := "Status cannot be deleted, resource doesn't support conditions"
		log.Info(errMessage)
		return fmt.Errorf(errMessage)
	}
}
