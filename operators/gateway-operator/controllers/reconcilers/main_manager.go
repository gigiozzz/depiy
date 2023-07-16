package reconcilers

import (
	"context"
	"fmt"
	"time"

	common "github.com/gigiozzz/depiy/common-libs/commons"
	"github.com/gigiozzz/depiy/operators/gateway-operator/api/v1alpha1"
	"github.com/gigiozzz/depiy/operators/gateway-operator/controllers/services"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const labelKey = "app"
const serverPortName = "server-port"

type ReconcileManager struct {
	Base      *common.BaseK8sStructure
	Scheme    *runtime.Scheme
	Recorder  record.EventRecorder
	Condition *services.ConditionService
}

func NewReconcileManager(client client.Client, log logr.Logger, scheme *runtime.Scheme, recorder record.EventRecorder) *ReconcileManager {
	base := &common.BaseK8sStructure{Client: client, Log: log}
	return &ReconcileManager{
		Base:      base,
		Scheme:    scheme,
		Recorder:  recorder,
		Condition: services.NewConditionService(base),
	}
}

func (r *ReconcileManager) MainReconcile(ctx context.Context, req ctrl.Request, cr *v1alpha1.EntandoGatewayV2) (ctrl.Result, error) {

	log := r.Base.Log
	ingressManager := NewIngressManager(r.Base, r.Condition)

	if err := r.Condition.SetConditionGatewayReadyUnknow(ctx, cr); err != nil {
		log.Info("error on set Gateway ingress ready unknow")
		return ctrl.Result{}, err
	}

	// deploy done
	applied := ingressManager.IsIngressApplied(ctx, cr)

	if !applied {
		if err := ingressManager.ApplyIngress(ctx, cr, r.Scheme); err != nil {
			log.Info("error ApplyDeploy reschedule reconcile", "error", err)
			r.Condition.SetConditionGatewayReadyFalse(ctx, cr)
			return ctrl.Result{}, err
		}
	}
	r.Recorder.Eventf(cr, "Normal", "Updated", fmt.Sprintf("Updated ingress %s/%s", req.Namespace, req.Name))

	// deploy ready
	var err error
	ready := ingressManager.IsIngressReady(ctx, cr)

	if !ready {
		if ready, err = ingressManager.CheckIngress(ctx, cr); err != nil {
			log.Info("error CheckIngress reschedule reconcile", "error", err)
			r.Condition.SetConditionGatewayReadyFalse(ctx, cr)
			return ctrl.Result{}, err
		}
		if !ready {
			log.Info("Ingress not ready reschedule operator", "seconds", 10)
			r.Recorder.Eventf(cr, "Warning", "NotReady", fmt.Sprintf("Gateway ingress not ready %s/%s", req.Namespace, req.Name))
			r.Condition.SetConditionGatewayReadyFalse(ctx, cr)
			return ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}, nil
		}
	}

	// ingress requested

	// ingress ready

	r.Recorder.Eventf(cr, "Normal", "Done", fmt.Sprintf("Gateway ingress deployed %s/%s", req.Namespace, req.Name))
	r.Condition.SetConditionGatewayReadyTrue(ctx, cr)
	return ctrl.Result{}, nil
}
