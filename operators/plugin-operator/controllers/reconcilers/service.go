package reconcilers

import (
	"context"

	utility "github.com/gigiozzz/depiy/common-libs/utilities"
	"github.com/gigiozzz/depiy/operators/plugin-operator/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"

	ctrl "sigs.k8s.io/controller-runtime"
)

func (d *ServiceManager) isServiceUpgrade(ctx context.Context, cr *v1alpha1.EntandoPluginV2, service *corev1.Service) (error, bool) {
	err := d.Base.Client.Get(ctx, types.NamespacedName{Name: MakeServiceName(cr), Namespace: cr.GetNamespace()}, service)
	if errors.IsNotFound(err) {
		return nil, false
	}
	return err, true
}

func (d *ServiceManager) buildService(cr *v1alpha1.EntandoPluginV2, scheme *runtime.Scheme) *corev1.Service {
	serviceName := MakeServiceName(cr)
	servicePort := MakeServicePort(cr)
	labels := map[string]string{labelKey: makeContainerName(cr)}
	port := int32(cr.Spec.Port)

	service := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: cr.GetNamespace(),
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{{
				Name:       servicePort,
				Port:       port,
				Protocol:   corev1.ProtocolTCP,
				TargetPort: intstr.IntOrString{StrVal: serverPortName, Type: intstr.String},
			}},
			Selector: labels,
		},
	}
	// set owner
	ctrl.SetControllerReference(cr, service, scheme)
	return service
}

func MakeServiceName(cr *v1alpha1.EntandoPluginV2) string {
	return utility.TruncateString(cr.GetName(), 208) + "-service"
}

func MakeServicePort(cr *v1alpha1.EntandoPluginV2) string {
	serviceName := MakeServiceName(cr)
	return utility.TruncateString(utility.GenerateSha256(serviceName), 9) + "-port"
}

func (d *ServiceManager) ApplyKubeService(ctx context.Context, cr *v1alpha1.EntandoPluginV2, scheme *runtime.Scheme) error {
	baseService := d.buildService(cr, scheme)
	service := &corev1.Service{}

	err, isUpgrade := d.isServiceUpgrade(ctx, cr, service)
	if err != nil {
		return err
	}

	var applyError error
	if isUpgrade {
		service.Spec = baseService.Spec
		applyError = d.Base.Client.Update(ctx, service)

	} else {
		applyError = d.Base.Client.Create(ctx, baseService)
	}

	if applyError != nil {
		return applyError
	}
	return nil
}
