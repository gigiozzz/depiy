package reconcilers

import (
	"context"
	"time"

	"github.com/gigiozzz/depiy/operators/plugin-operator/api/v1alpha1"

	common "github.com/gigiozzz/depiy/common-libs/commons"
	"github.com/gigiozzz/depiy/operators/plugin-operator/controllers/services"

	"k8s.io/apimachinery/pkg/runtime"
)

type DeployManager struct {
	Base       *common.BaseK8sStructure
	Conditions *services.ConditionService
}

func NewDeployManager(base *common.BaseK8sStructure, conditions *services.ConditionService) *DeployManager {
	return &DeployManager{
		Base:       base,
		Conditions: conditions,
	}
}

func (d *DeployManager) IsDeployApplied(ctx context.Context, cr *v1alpha1.EntandoPluginV2) bool {

	return d.Conditions.IsDeployApplied(ctx, cr)
}

func (d *DeployManager) IsDeployReady(ctx context.Context, cr *v1alpha1.EntandoPluginV2) bool {

	return d.Conditions.IsDeployReady(ctx, cr)
}

func (d *DeployManager) ApplyDeploy(ctx context.Context, cr *v1alpha1.EntandoPluginV2, scheme *runtime.Scheme) error {
	applyError := d.ApplyKubeDeployment(ctx, cr, scheme)
	if applyError != nil {
		return applyError
	}

	return d.Conditions.SetConditionDeployApplied(ctx, cr)
}

func (d *DeployManager) CheckDeploy(ctx context.Context, cr *v1alpha1.EntandoPluginV2) (bool, error) {
	time.Sleep(time.Second * 10)
	ready := true
	// check condition "Available" is "True"
	if ready {
		return ready, d.Conditions.SetConditionDeployReady(ctx, cr)
	}

	return ready, nil

}
