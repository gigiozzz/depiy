package reconcilers

import (
	"context"
	"time"

	"github.com/gigiozzz/depiy/operators/plugin-operator/api/v1alpha1"

	common "github.com/gigiozzz/depiy/common-libs/commons"
	"github.com/gigiozzz/depiy/operators/plugin-operator/controllers/services"

	"k8s.io/apimachinery/pkg/runtime"
)

type ServiceManager struct {
	Base       *common.BaseK8sStructure
	Conditions *services.ConditionService
}

func NewServiceManager(base *common.BaseK8sStructure, conditions *services.ConditionService) *ServiceManager {
	return &ServiceManager{
		Base:       base,
		Conditions: conditions,
	}
}

func (d *ServiceManager) IsServiceApplied(ctx context.Context, cr *v1alpha1.EntandoPluginV2) bool {

	return d.Conditions.IsServiceApplied(ctx, cr)
}

func (d *ServiceManager) IsServiceReady(ctx context.Context, cr *v1alpha1.EntandoPluginV2) bool {

	return d.Conditions.IsServiceReady(ctx, cr)
}

func (d *ServiceManager) ApplyService(ctx context.Context, cr *v1alpha1.EntandoPluginV2, scheme *runtime.Scheme) error {
	applyError := d.ApplyKubeService(ctx, cr, scheme)
	if applyError != nil {
		return applyError
	}

	return d.Conditions.SetConditionServiceApplied(ctx, cr)
}

func (d *ServiceManager) CheckService(ctx context.Context, cr *v1alpha1.EntandoPluginV2) (bool, error) {
	time.Sleep(time.Second * 5)
	ready := true

	if ready {
		return ready, d.Conditions.SetConditionServiceReady(ctx, cr)
	}

	return ready, nil

}
