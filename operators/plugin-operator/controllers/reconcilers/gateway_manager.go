package reconcilers

import (
	"context"

	common "github.com/gigiozzz/depiy/common-libs/commons"
	utility "github.com/gigiozzz/depiy/common-libs/utilities"
	gwapi "github.com/gigiozzz/depiy/operators/gateway-operator/api/v1alpha1"
	gwservice "github.com/gigiozzz/depiy/operators/gateway-operator/controllers/services"
	"github.com/gigiozzz/depiy/operators/plugin-operator/api/v1alpha1"
	"github.com/gigiozzz/depiy/operators/plugin-operator/controllers/services"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
)

type GatewayManager struct {
	Base       *common.BaseK8sStructure
	Conditions *services.ConditionService
}

func NewGatewayManager(base *common.BaseK8sStructure, conditions *services.ConditionService) *GatewayManager {
	return &GatewayManager{
		Base:       base,
		Conditions: conditions,
	}
}

func (d *GatewayManager) IsCrApplied(ctx context.Context, cr *v1alpha1.EntandoPluginV2) bool {

	return d.Conditions.IsGatewayCrApplied(ctx, cr)
}

func (d *GatewayManager) IsCrReady(ctx context.Context, cr *v1alpha1.EntandoPluginV2) bool {

	return d.Conditions.IsGatewayCrReady(ctx, cr)
}

func (d *GatewayManager) buildCr(cr *v1alpha1.EntandoPluginV2, scheme *runtime.Scheme) *gwapi.EntandoGatewayV2 {
	crName := makeCrName(cr)
	gatewayCR := &gwapi.EntandoGatewayV2{
		ObjectMeta: metav1.ObjectMeta{
			Name:      crName,
			Namespace: cr.GetNamespace(),
		},
		Spec: gwapi.EntandoGatewayV2Spec{
			IngressName:    getIngressNameOrDefault(cr),
			IngressHost:    cr.Spec.IngressHost,
			IngressPath:    cr.Spec.IngressPath,
			IngressPort:    MakeServicePort(cr),
			IngressService: MakeServiceName(cr),
		},
	}
	// set owner
	ctrl.SetControllerReference(cr, gatewayCR, scheme)
	return gatewayCR
}

func makeCrName(cr *v1alpha1.EntandoPluginV2) string {
	return utility.TruncateString(cr.GetName(), 208) + "-gateway"
}

func getIngressNameOrDefault(cr *v1alpha1.EntandoPluginV2) string {
	var name string = cr.Spec.IngressName
	if len(name) <= 0 {
		name = utility.TruncateString(cr.GetName(), 208) + "-ingress"
	}
	return name
}

func (d *GatewayManager) isCrUpgrade(ctx context.Context, cr *v1alpha1.EntandoPluginV2, gatewayCr *gwapi.EntandoGatewayV2) (error, bool) {
	err := d.Base.Client.Get(ctx, types.NamespacedName{Name: makeCrName(cr), Namespace: cr.GetNamespace()}, gatewayCr)
	if errors.IsNotFound(err) {
		return nil, false
	}
	return err, true
}

func (d *GatewayManager) ApplyCr(ctx context.Context, cr *v1alpha1.EntandoPluginV2, scheme *runtime.Scheme) error {

	baseGatewayCr := d.buildCr(cr, scheme)
	gatewayCr := &gwapi.EntandoGatewayV2{}

	err, isUpgrade := d.isCrUpgrade(ctx, cr, gatewayCr)
	if err != nil {
		return err
	}

	var applyError error
	if isUpgrade {
		gatewayCr.Spec = baseGatewayCr.Spec
		applyError = d.Base.Client.Update(ctx, gatewayCr)

	} else {
		applyError = d.Base.Client.Create(ctx, baseGatewayCr)
	}

	if applyError != nil {
		return applyError
	}

	return d.Conditions.SetConditionGatewayCrApplied(ctx, cr)
}

func (d *GatewayManager) CheckCr(ctx context.Context, cr *v1alpha1.EntandoPluginV2) (bool, error) {

	gatewayCr := &gwapi.EntandoGatewayV2{}
	err := d.Base.Client.Get(ctx, types.NamespacedName{Name: makeCrName(cr), Namespace: cr.GetNamespace()}, gatewayCr)

	if err != nil {
		return false, err
	}

	ready := d.checkCrCondition(gatewayCr)

	if ready {
		return ready, d.Conditions.SetConditionGatewayCrReady(ctx, cr)
	}

	return ready, nil

}

func (d *GatewayManager) checkCrCondition(cr *gwapi.EntandoGatewayV2) bool {
	var output metav1.ConditionStatus = metav1.ConditionUnknown
	var observedGeneration int64

	for _, condition := range cr.Status.Conditions {
		if condition.Type == gwservice.CONDITION_GATEWAY_READY {
			output = condition.Status
			observedGeneration = condition.ObservedGeneration
		}
	}
	return metav1.ConditionTrue == output && observedGeneration == cr.Generation

}
