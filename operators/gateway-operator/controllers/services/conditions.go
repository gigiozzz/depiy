package services

import (
	"context"

	common "github.com/gigiozzz/depiy/common-libs/commons"
	utility "github.com/gigiozzz/depiy/common-libs/utilities"
	"github.com/gigiozzz/depiy/operators/gateway-operator/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	CONDITION_INGRESS_APPLIED        = "IngressApplied"
	CONDITION_INGRESS_APPLIED_REASON = "IngressIsApplied"
	CONDITION_INGRESS_APPLIED_MSG    = "Your ingress was applied"

	CONDITION_INGRESS_READY        = "IngressReady"
	CONDITION_INGRESS_READY_REASON = "IngressIsReady"
	CONDITION_INGRESS_READY_MSG    = "Your ingress is ready"

	CONDITION_GATEWAY_READY        = "Ready"
	CONDITION_GATEWAY_READY_REASON = "GatewayIsReady"
	CONDITION_GATEWAY_READY_MSG    = "Your Gateway ingress is ready"
)

type ConditionService struct {
	Base *common.BaseK8sStructure
}

func NewConditionService(base *common.BaseK8sStructure) *ConditionService {
	return &ConditionService{
		Base: base,
	}
}

func (cs *ConditionService) IsIngressReady(ctx context.Context, cr *v1alpha1.EntandoGatewayV2) bool {

	condition, observedGeneration := cs.getConditionStatus(ctx, cr, CONDITION_INGRESS_READY)

	return metav1.ConditionTrue == condition && observedGeneration == cr.Generation
}

func (cs *ConditionService) SetConditionIngressReady(ctx context.Context, cr *v1alpha1.EntandoGatewayV2) error {

	cs.deleteCondition(ctx, cr, CONDITION_INGRESS_READY)
	return utility.AppendCondition(ctx, cs.Base.Client, cr,
		CONDITION_INGRESS_READY,
		metav1.ConditionTrue,
		CONDITION_INGRESS_READY_REASON,
		CONDITION_INGRESS_READY_MSG,
		cr.Generation)
}

func (cs *ConditionService) IsIngressApplied(ctx context.Context, cr *v1alpha1.EntandoGatewayV2) bool {

	condition, observedGeneration := cs.getConditionStatus(ctx, cr, CONDITION_INGRESS_APPLIED)

	return metav1.ConditionTrue == condition && observedGeneration == cr.Generation
}

func (cs *ConditionService) SetConditionIngressApplied(ctx context.Context, cr *v1alpha1.EntandoGatewayV2) error {

	cs.deleteCondition(ctx, cr, CONDITION_INGRESS_APPLIED)
	return utility.AppendCondition(ctx, cs.Base.Client, cr,
		CONDITION_INGRESS_APPLIED,
		metav1.ConditionTrue,
		CONDITION_INGRESS_APPLIED_REASON,
		CONDITION_INGRESS_APPLIED_MSG,
		cr.Generation)
}

func (cs *ConditionService) SetConditionGatewayReadyTrue(ctx context.Context, cr *v1alpha1.EntandoGatewayV2) error {
	return cs.setConditionGatewayReady(ctx, cr, metav1.ConditionTrue)
}

func (cs *ConditionService) SetConditionGatewayReadyUnknow(ctx context.Context, cr *v1alpha1.EntandoGatewayV2) error {
	return cs.setConditionGatewayReady(ctx, cr, metav1.ConditionUnknown)
}

func (cs *ConditionService) SetConditionGatewayReadyFalse(ctx context.Context, cr *v1alpha1.EntandoGatewayV2) error {
	return cs.setConditionGatewayReady(ctx, cr, metav1.ConditionFalse)
}

func (cs *ConditionService) getConditionStatus(ctx context.Context, cr *v1alpha1.EntandoGatewayV2, typeName string) (metav1.ConditionStatus, int64) {

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

func (cs *ConditionService) setConditionGatewayReady(ctx context.Context, cr *v1alpha1.EntandoGatewayV2, status metav1.ConditionStatus) error {

	cs.deleteCondition(ctx, cr, CONDITION_GATEWAY_READY)
	return utility.AppendCondition(ctx, cs.Base.Client, cr,
		CONDITION_GATEWAY_READY,
		status,
		CONDITION_GATEWAY_READY_REASON,
		CONDITION_GATEWAY_READY_MSG,
		cr.Generation)
}

func (cs *ConditionService) deleteCondition(ctx context.Context, cr *v1alpha1.EntandoGatewayV2, typeName string) error {

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
