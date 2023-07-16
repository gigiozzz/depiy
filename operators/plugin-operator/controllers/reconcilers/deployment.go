package reconcilers

import (
	"context"

	utility "github.com/gigiozzz/depiy/common-libs/utilities"
	"github.com/gigiozzz/depiy/operators/plugin-operator/api/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"

	ctrl "sigs.k8s.io/controller-runtime"
)

func (d *DeployManager) isDeploymentUpgrade(ctx context.Context, cr *v1alpha1.EntandoPluginV2, deployment *appsv1.Deployment) (error, bool) {
	err := d.Base.Client.Get(ctx, types.NamespacedName{Name: makeDeploymentName(cr), Namespace: cr.GetNamespace()}, deployment)
	if errors.IsNotFound(err) {
		return nil, false
	}
	return err, true
}

func (d *DeployManager) buildDeployment(cr *v1alpha1.EntandoPluginV2, scheme *runtime.Scheme) *appsv1.Deployment {
	replicas := cr.Spec.Replicas
	deploymentName := makeDeploymentName(cr)
	containerName := makeContainerName(cr)
	labels := map[string]string{labelKey: containerName}
	port := int32(cr.Spec.Port)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: cr.GetNamespace(),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: &intstr.IntOrString{StrVal: "25%", Type: intstr.String},
					MaxSurge:       &intstr.IntOrString{StrVal: "25%", Type: intstr.String},
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           cr.Spec.Image,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Name:            containerName,
						Ports: []corev1.ContainerPort{{
							ContainerPort: port,
							Name:          serverPortName,
						}},
						Env: cr.Spec.EnvironmentVariables,
						ReadinessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{Path: cr.Spec.HealthCheckPath, Port: intstr.IntOrString{
									IntVal: port,
								}},
							},
							InitialDelaySeconds: 10,
						},
						LivenessProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{Path: cr.Spec.HealthCheckPath, Port: intstr.IntOrString{
									IntVal: port,
								}},
							},
							InitialDelaySeconds: 10,
						},
						StartupProbe: &corev1.Probe{
							ProbeHandler: corev1.ProbeHandler{
								HTTPGet: &corev1.HTTPGetAction{Path: cr.Spec.HealthCheckPath, Port: intstr.IntOrString{
									IntVal: port,
								}},
							},
							InitialDelaySeconds: 20,
						},
					}},
				},
			},
		},
	}
	// set owner
	ctrl.SetControllerReference(cr, deployment, scheme)
	return deployment
}

func makeContainerName(cr *v1alpha1.EntandoPluginV2) string {
	return utility.TruncateString(cr.GetName(), 208) + "-container"
}

func makeDeploymentName(cr *v1alpha1.EntandoPluginV2) string {
	return utility.TruncateString(cr.GetName(), 208) + "-deployment"
}

func (d *DeployManager) ApplyKubeDeployment(ctx context.Context, cr *v1alpha1.EntandoPluginV2, scheme *runtime.Scheme) error {
	baseDeployment := d.buildDeployment(cr, scheme)
	deployment := &appsv1.Deployment{}

	err, isUpgrade := d.isDeploymentUpgrade(ctx, cr, deployment)
	if err != nil {
		return err
	}

	var applyError error
	if isUpgrade {
		deployment.Spec = baseDeployment.Spec
		applyError = d.Base.Client.Update(ctx, deployment)

	} else {
		applyError = d.Base.Client.Create(ctx, baseDeployment)
	}

	if applyError != nil {
		return applyError
	}
	return nil
}
