package services

import (
	"context"

	common "github.com/gigiozzz/depiy/common-libs/commons"
	utility "github.com/gigiozzz/depiy/common-libs/utilities"
	"github.com/gigiozzz/depiy/operators/plugin-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	CONDITION_GATEWAY_CR_APPLIED        = "GatewayCrApplied"
	CONDITION_GATEWAY_CR_APPLIED_REASON = "GatewayCrIsApplied"
	CONDITION_GATEWAY_CR_APPLIED_MSG    = "Your gateway cr was applied"

	CONDITION_GATEWAY_CR_READY        = "GatewayCrReady"
	CONDITION_GATEWAY_CR_READY_REASON = "GatewayCrIsReady"
	CONDITION_GATEWAY_CR_READY_MSG    = "Your gateway cr is ready"

	CONDITION_SERVICE_APPLIED        = "ServiceApplied"
	CONDITION_SERVICE_APPLIED_REASON = "ServiceIsApplied"
	CONDITION_SERVICE_APPLIED_MSG    = "Your service was applied"

	CONDITION_SERVICE_READY        = "ServiceReady"
	CONDITION_SERVICE_READY_REASON = "ServiceIsReady"
	CONDITION_SERVICE_READY_MSG    = "Your service is ready"

	CONDITION_DEPLOY_APPLIED        = "DeployApplied"
	CONDITION_DEPLOY_APPLIED_REASON = "DeployIsApplied"
	CONDITION_DEPLOY_APPLIED_MSG    = "Your deploy was applied"

	CONDITION_DEPLOY_READY        = "DeployReady"
	CONDITION_DEPLOY_READY_REASON = "DeployIsReady"
	CONDITION_DEPLOY_READY_MSG    = "Your deploy is ready"

	CONDITION_PLUGIN_READY        = "Ready"
	CONDITION_PLUGIN_READY_REASON = "PluginIsReady"
	CONDITION_PLUGIN_READY_MSG    = "Your plugin is ready"
)

type ConditionService struct {
	Base *common.BaseK8sStructure
}

func NewConditionService(base *common.BaseK8sStructure) *ConditionService {
	return &ConditionService{
		Base: base,
	}
}

func (cs *ConditionService) IsGatewayCrReady(ctx context.Context, cr *v1alpha1.EntandoPluginV2) bool {

	condition, observedGeneration := cs.getConditionStatus(ctx, cr, CONDITION_GATEWAY_CR_READY)

	return metav1.ConditionTrue == condition && observedGeneration == cr.Generation
}

func (cs *ConditionService) SetConditionGatewayCrReady(ctx context.Context, cr *v1alpha1.EntandoPluginV2) error {

	cs.deleteCondition(ctx, cr, CONDITION_GATEWAY_CR_READY)
	return utility.AppendCondition(ctx, cs.Base.Client, cr,
		CONDITION_GATEWAY_CR_READY,
		metav1.ConditionTrue,
		CONDITION_GATEWAY_CR_READY_REASON,
		CONDITION_GATEWAY_CR_READY_MSG,
		cr.Generation)
}

func (cs *ConditionService) IsGatewayCrApplied(ctx context.Context, cr *v1alpha1.EntandoPluginV2) bool {

	condition, observedGeneration := cs.getConditionStatus(ctx, cr, CONDITION_GATEWAY_CR_APPLIED)

	return metav1.ConditionTrue == condition && observedGeneration == cr.Generation
}

func (cs *ConditionService) SetConditionGatewayCrApplied(ctx context.Context, cr *v1alpha1.EntandoPluginV2) error {

	cs.deleteCondition(ctx, cr, CONDITION_GATEWAY_CR_APPLIED)
	return utility.AppendCondition(ctx, cs.Base.Client, cr,
		CONDITION_GATEWAY_CR_APPLIED,
		metav1.ConditionTrue,
		CONDITION_GATEWAY_CR_APPLIED_REASON,
		CONDITION_GATEWAY_CR_APPLIED_MSG,
		cr.Generation)
}

func (cs *ConditionService) IsServiceReady(ctx context.Context, cr *v1alpha1.EntandoPluginV2) bool {

	condition, observedGeneration := cs.getConditionStatus(ctx, cr, CONDITION_SERVICE_READY)

	return metav1.ConditionTrue == condition && observedGeneration == cr.Generation
}

func (cs *ConditionService) SetConditionServiceReady(ctx context.Context, cr *v1alpha1.EntandoPluginV2) error {

	cs.deleteCondition(ctx, cr, CONDITION_SERVICE_READY)
	return utility.AppendCondition(ctx, cs.Base.Client, cr,
		CONDITION_SERVICE_READY,
		metav1.ConditionTrue,
		CONDITION_SERVICE_READY_REASON,
		CONDITION_SERVICE_READY_MSG,
		cr.Generation)
}

func (cs *ConditionService) IsServiceApplied(ctx context.Context, cr *v1alpha1.EntandoPluginV2) bool {

	condition, observedGeneration := cs.getConditionStatus(ctx, cr, CONDITION_SERVICE_APPLIED)

	return metav1.ConditionTrue == condition && observedGeneration == cr.Generation
}

func (cs *ConditionService) SetConditionServiceApplied(ctx context.Context, cr *v1alpha1.EntandoPluginV2) error {

	cs.deleteCondition(ctx, cr, CONDITION_SERVICE_APPLIED)
	return utility.AppendCondition(ctx, cs.Base.Client, cr,
		CONDITION_SERVICE_APPLIED,
		metav1.ConditionTrue,
		CONDITION_SERVICE_APPLIED_REASON,
		CONDITION_SERVICE_APPLIED_MSG,
		cr.Generation)
}

func (cs *ConditionService) IsDeployReady(ctx context.Context, cr *v1alpha1.EntandoPluginV2) bool {

	condition, observedGeneration := cs.getConditionStatus(ctx, cr, CONDITION_DEPLOY_READY)

	return metav1.ConditionTrue == condition && observedGeneration == cr.Generation
}

func (cs *ConditionService) SetConditionDeployReady(ctx context.Context, cr *v1alpha1.EntandoPluginV2) error {

	cs.deleteCondition(ctx, cr, CONDITION_DEPLOY_READY)
	return utility.AppendCondition(ctx, cs.Base.Client, cr,
		CONDITION_DEPLOY_READY,
		metav1.ConditionTrue,
		CONDITION_DEPLOY_READY_REASON,
		CONDITION_DEPLOY_READY_MSG,
		cr.Generation)
}

func (cs *ConditionService) IsDeployApplied(ctx context.Context, cr *v1alpha1.EntandoPluginV2) bool {

	condition, observedGeneration := cs.getConditionStatus(ctx, cr, CONDITION_DEPLOY_APPLIED)

	return metav1.ConditionTrue == condition && observedGeneration == cr.Generation
}

func (cs *ConditionService) SetConditionDeployApplied(ctx context.Context, cr *v1alpha1.EntandoPluginV2) error {

	cs.deleteCondition(ctx, cr, CONDITION_DEPLOY_APPLIED)
	return utility.AppendCondition(ctx, cs.Base.Client, cr,
		CONDITION_DEPLOY_APPLIED,
		metav1.ConditionTrue,
		CONDITION_DEPLOY_APPLIED_REASON,
		CONDITION_DEPLOY_APPLIED_MSG,
		cr.Generation)
}

func (cs *ConditionService) SetConditionPluginReadyTrue(ctx context.Context, cr *v1alpha1.EntandoPluginV2) error {
	return cs.setConditionPluginReady(ctx, cr, metav1.ConditionTrue)
}

func (cs *ConditionService) SetConditionPluginReadyUnknow(ctx context.Context, cr *v1alpha1.EntandoPluginV2) error {
	return cs.setConditionPluginReady(ctx, cr, metav1.ConditionUnknown)
}

func (cs *ConditionService) SetConditionPluginReadyFalse(ctx context.Context, cr *v1alpha1.EntandoPluginV2) error {
	return cs.setConditionPluginReady(ctx, cr, metav1.ConditionFalse)
}

func (cs *ConditionService) getConditionStatus(ctx context.Context, cr *v1alpha1.EntandoPluginV2, typeName string) (metav1.ConditionStatus, int64) {

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

func (cs *ConditionService) setConditionPluginReady(ctx context.Context, cr *v1alpha1.EntandoPluginV2, status metav1.ConditionStatus) error {

	cs.deleteCondition(ctx, cr, CONDITION_PLUGIN_READY)
	return utility.AppendCondition(ctx, cs.Base.Client, cr,
		CONDITION_PLUGIN_READY,
		status,
		CONDITION_PLUGIN_READY_REASON,
		CONDITION_PLUGIN_READY_MSG,
		cr.Generation)
}

func (cs *ConditionService) deleteCondition(ctx context.Context, cr *v1alpha1.EntandoPluginV2, typeName string) error {

	log := log.FromContext(ctx)
	var newConditions = make([]metav1.Condition, 0)
	for _, condition := range cr.Status.Conditions {
		if condition.Type != typeName {
			newConditions = append(newConditions, condition)
		}
	}
	cr.Status.Conditions = newConditions

	err := cs.Base.Client.Status().Update(ctx, cr)
	if err != nil {
		log.Info("Application resource status update failed.")
	}
	return nil
}
