package instance

import (
	"context"
	"strings"
	"time"

	common "github.com/gigiozzz/depiy/common-libs/commons"
	utility "github.com/gigiozzz/depiy/common-libs/utilities"
	"github.com/gigiozzz/depiy/operators/bundle-operator/api/v1alpha1"
	"github.com/gigiozzz/depiy/operators/bundle-operator/bundles"
	"github.com/gigiozzz/depiy/operators/bundle-operator/controllers/services"

	pluginapi "github.com/gigiozzz/depiy/operators/plugin-operator/api/v1alpha1"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

type PluginManager struct {
	Base       *common.BaseK8sStructure
	Conditions *services.ConditionService
}

func NewPluginManager(base *common.BaseK8sStructure, conditions *services.ConditionService) *PluginManager {
	return &PluginManager{
		Base:       base,
		Conditions: conditions,
	}
}

func (d *PluginManager) IsPluginApplied(ctx context.Context, cr *v1alpha1.EntandoBundleInstanceV2,
	plugin *bundles.Plugin) bool {

	return d.Conditions.IsPluginCrApplied(ctx, cr, d.GenPluginCode(cr, plugin))
}

func (d *PluginManager) IsPluginReady(ctx context.Context, cr *v1alpha1.EntandoBundleInstanceV2) bool {

	return d.Conditions.IsPluginCrReady(ctx, cr)
}

func (d *PluginManager) ApplyPlugin(ctx context.Context, cr *v1alpha1.EntandoBundleInstanceV2, plugin *bundles.Plugin,
	scheme *runtime.Scheme) error {
	log := d.Base.Log
	basePluginCr := d.buildPluginCr(cr, plugin, scheme)
	log.Info("generated plugin", "pluginCR", basePluginCr)
	pluginCr := &pluginapi.EntandoPluginV2{}

	err, isUpgrade := d.isCrUpgrade(ctx, cr, pluginCr, plugin)
	if err != nil {
		return err
	}

	var applyError error
	if isUpgrade {
		pluginCr.Spec = basePluginCr.Spec
		log.Info("Update plugin cr", "pluginCR", pluginCr)
		applyError = d.Base.Client.Update(ctx, pluginCr)

	} else {
		log.Info("Create plugin cr", "pluginCR", basePluginCr)
		applyError = d.Base.Client.Create(ctx, basePluginCr)
	}

	if applyError != nil {
		return applyError
	}

	return d.Conditions.SetConditionPluginCrApplied(ctx, cr, d.GenPluginCode(cr, plugin))
}

func (d *PluginManager) CheckPluginCr(ctx context.Context, cr *v1alpha1.EntandoBundleInstanceV2) (bool, error) {
	time.Sleep(time.Second * 10)
	ready := true
	// check condition "Available" is "True"
	if ready {
		return ready, d.Conditions.SetConditionPluginCrReady(ctx, cr)
	}

	return ready, nil

}

func (d *PluginManager) isCrUpgrade(ctx context.Context, cr *v1alpha1.EntandoBundleInstanceV2,
	pluginCr *pluginapi.EntandoPluginV2,
	plugin *bundles.Plugin) (error, bool) {
	err := d.Base.Client.Get(ctx, types.NamespacedName{Name: d.GenPluginCode(cr, plugin), Namespace: cr.GetNamespace()}, pluginCr)
	if errors.IsNotFound(err) {
		return nil, false
	}
	return err, true
}

func (d *PluginManager) GenPluginId(instanceCode string, plugin *bundles.Plugin) string {
	pluginFullInfo := plugin.Repository + "@" + plugin.Digest + "-" + instanceCode
	s := utility.GenerateSha256(pluginFullInfo)
	return utility.TruncateString(s, 8)
}

func (d *PluginManager) GenPluginCode(cr *v1alpha1.EntandoBundleInstanceV2, plugin *bundles.Plugin) string {
	splittedName := strings.Split(cr.ObjectMeta.Name, "-")
	instanceCode := splittedName[2]
	bundleName := splittedName[0] + "-" + splittedName[1]
	pluginId := d.GenPluginId(instanceCode, plugin)
	pluginCode := strings.ToLower(bundleName + "-pn-" + pluginId + "-" + utility.TruncateString(cr.Name, 180))
	d.Base.Log.Info("generated pluginCode", "pluginCode", pluginCode)
	return pluginCode
}

func (d *PluginManager) buildPluginCr(cr *v1alpha1.EntandoBundleInstanceV2, plugin *bundles.Plugin, scheme *runtime.Scheme) *pluginapi.EntandoPluginV2 {
	pluginCode := d.GenPluginCode(cr, plugin)
	pluginCr := &pluginapi.EntandoPluginV2{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pluginCode,
			Namespace: cr.GetNamespace(),
		},
		Spec: pluginapi.EntandoPluginV2Spec{
			Replicas:        1,
			Image:           plugin.Repository + "@" + plugin.Digest,
			HealthCheckPath: plugin.HealthCheckPath,
			Port:            int32(plugin.Port),
			IngressName:     plugin.IngressName,
			IngressHost:     plugin.IngressHost,
			IngressPath:     plugin.IngressPath,
		},
	}
	// set owner
	ctrl.SetControllerReference(cr, pluginCr, scheme)
	return pluginCr
}
