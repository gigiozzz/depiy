package reconcilers

import (
	"context"

	"github.com/gigiozzz/depiy/operators/gateway-operator/api/v1alpha1"

	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	ctrl "sigs.k8s.io/controller-runtime"
)

func (d *IngressManager) isIngressUpgrade(ctx context.Context, cr *v1alpha1.EntandoGatewayV2, ingress *netv1.Ingress) (error, bool) {
	err := d.Base.Client.Get(ctx, types.NamespacedName{Name: cr.Spec.IngressName, Namespace: cr.GetNamespace()}, ingress)
	if errors.IsNotFound(err) {
		return nil, false
	}
	return err, true
}

func (d *IngressManager) buildIngress(cr *v1alpha1.EntandoGatewayV2, scheme *runtime.Scheme) *netv1.Ingress {
	pp := netv1.PathTypePrefix
	ingress := &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Spec.IngressName,
			Namespace: cr.GetNamespace(),
		},
		Spec: netv1.IngressSpec{
			Rules: []netv1.IngressRule{{
				Host: cr.Spec.IngressHost,
				IngressRuleValue: netv1.IngressRuleValue{
					HTTP: &netv1.HTTPIngressRuleValue{
						Paths: []netv1.HTTPIngressPath{{
							Path:     cr.Spec.IngressPath,
							PathType: &pp,
							Backend: netv1.IngressBackend{
								Service: &netv1.IngressServiceBackend{
									Name: cr.Spec.IngressService,
									Port: netv1.ServiceBackendPort{
										Name: cr.Spec.IngressPort,
									},
								},
							},
						}},
					},
				},
			}},
		},
	} // set owner
	ctrl.SetControllerReference(cr, ingress, scheme)
	return ingress
}

func (d *IngressManager) updateIngressSpec(ingress *netv1.Ingress, baseIngress *netv1.Ingress, cr *v1alpha1.EntandoGatewayV2) {
	found := false
	for _, rule := range ingress.Spec.Rules {
		if rule.Host == cr.Spec.IngressHost {
			found = true
			pathFound := false
			for _, path := range rule.HTTP.Paths {
				if path.Path == cr.Spec.IngressPath {
					pathFound = true
					path.Backend = baseIngress.Spec.Rules[0].IngressRuleValue.HTTP.Paths[0].Backend
					break
				}
			}
			if !pathFound {
				rule.HTTP.Paths = append(rule.HTTP.Paths, baseIngress.Spec.Rules[0].IngressRuleValue.HTTP.Paths[0])
			}
			break
		}
	}

	if !found {
		ingress.Spec.Rules = append(ingress.Spec.Rules, baseIngress.Spec.Rules[0])
	}
}

func (d *IngressManager) ApplyKubeIngress(ctx context.Context, cr *v1alpha1.EntandoGatewayV2, scheme *runtime.Scheme) error {
	baseIngress := d.buildIngress(cr, scheme)
	ingress := &netv1.Ingress{}

	err, isUpgrade := d.isIngressUpgrade(ctx, cr, ingress)
	if err != nil {
		return err
	}

	var applyError error
	if isUpgrade {
		d.updateIngressSpec(ingress, baseIngress, cr)
		applyError = d.Base.Client.Update(ctx, ingress)

	} else {
		applyError = d.Base.Client.Create(ctx, baseIngress)
	}

	if applyError != nil {
		return applyError
	}
	return nil
}
