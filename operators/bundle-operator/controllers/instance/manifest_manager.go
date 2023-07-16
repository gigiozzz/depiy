package instance

import (
	"context"

	utility "github.com/gigiozzz/depiy/common-libs/utilities"
	"github.com/gigiozzz/depiy/operators/bundle-operator/api/v1alpha1"

	common "github.com/gigiozzz/depiy/common-libs/commons"
	"github.com/gigiozzz/depiy/operators/bundle-operator/controllers/services"

	runtime "k8s.io/apimachinery/pkg/runtime"
)

type ManifestManager struct {
	Base       *common.BaseK8sStructure
	Conditions *services.ConditionService
}

func NewManifestManager(base *common.BaseK8sStructure, conditions *services.ConditionService) *ManifestManager {
	return &ManifestManager{
		Base:       base,
		Conditions: conditions,
	}
}

func (d *ManifestManager) IsManifestApplied(ctx context.Context, cr *v1alpha1.EntandoBundleInstanceV2, manifestPath string) bool {
	manifestId := genManifestId(cr, manifestPath)
	return d.Conditions.IsManifestApplied(ctx, cr, manifestId)
}

func (d *ManifestManager) ApplyManifest(ctx context.Context, cr *v1alpha1.EntandoBundleInstanceV2,
	scheme *runtime.Scheme,
	dir string,
	manifestPath string) error {

	manifestService := NewManifest(d.Base)

	err := manifestService.ApplyManifest(ctx, cr, scheme, dir+manifestPath)
	if err != nil {
		return err
	}
	manifestId := genManifestId(cr, manifestPath)

	return d.Conditions.SetConditionManifestApplied(ctx, cr, manifestId, manifestPath)
}

func genManifestId(cr *v1alpha1.EntandoBundleInstanceV2, manifestPath string) string {
	s := utility.GenerateSha256(cr.Spec.Digest + manifestPath)
	return utility.TruncateString(s, 8)
}
