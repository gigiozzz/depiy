package reconcilers

import (
	"context"
	"fmt"
	"time"

	common "github.com/gigiozzz/depiy/common-libs/commons"
	"github.com/gigiozzz/depiy/operators/plugin-operator/api/v1alpha1"
	"github.com/gigiozzz/depiy/operators/plugin-operator/controllers/services"
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

func (r *ReconcileManager) MainReconcile(ctx context.Context, req ctrl.Request, cr *v1alpha1.EntandoPluginV2) (ctrl.Result, error) {

	log := r.Base.Log
	deployManager := NewDeployManager(r.Base, r.Condition)
	serviceManager := NewServiceManager(r.Base, r.Condition)
	gatewayManager := NewGatewayManager(r.Base, r.Condition)

	if err := r.Condition.SetConditionPluginReadyUnknow(ctx, cr); err != nil {
		log.Info("error on set plugin ready unknow")
		return ctrl.Result{}, err
	}

	// deploy done
	applied := deployManager.IsDeployApplied(ctx, cr)

	if !applied {
		if err := deployManager.ApplyDeploy(ctx, cr, r.Scheme); err != nil {
			log.Info("error ApplyDeploy reschedule reconcile", "error", err)
			r.Condition.SetConditionPluginReadyFalse(ctx, cr)
			return ctrl.Result{}, err
		}
	}
	r.Recorder.Eventf(cr, "Normal", "Updated", fmt.Sprintf("Updated deployment %s/%s", req.Namespace, req.Name))

	// deploy ready
	var err error
	ready := deployManager.IsDeployReady(ctx, cr)

	if !ready {
		if ready, err = deployManager.CheckDeploy(ctx, cr); err != nil {
			log.Info("error CheckDeploy reschedule reconcile", "error", err)
			r.Condition.SetConditionPluginReadyFalse(ctx, cr)
			return ctrl.Result{}, err
		}
		if !ready {
			log.Info("Deploy not ready reschedule operator", "seconds", 10)
			r.Recorder.Eventf(cr, "Warning", "NotReady", fmt.Sprintf("Plugin deployment not ready %s/%s", req.Namespace, req.Name))
			r.Condition.SetConditionPluginReadyFalse(ctx, cr)
			return ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}, nil
		}
	}

	// service done
	applied = serviceManager.IsServiceApplied(ctx, cr)

	if !applied {
		if err := serviceManager.ApplyService(ctx, cr, r.Scheme); err != nil {
			log.Info("error ApplyService reschedule reconcile", "error", err)
			r.Condition.SetConditionPluginReadyFalse(ctx, cr)
			return ctrl.Result{}, err
		}
	}
	r.Recorder.Eventf(cr, "Normal", "Updated", fmt.Sprintf("Updated service %s/%s", req.Namespace, req.Name))

	// service ready
	ready = serviceManager.IsServiceReady(ctx, cr)

	if !ready {
		if ready, err = serviceManager.CheckService(ctx, cr); err != nil {
			log.Info("error CheckService reschedule reconcile", "error", err)
			r.Condition.SetConditionPluginReadyFalse(ctx, cr)
			return ctrl.Result{}, err
		}
		if !ready {
			log.Info("Service not ready reschedule operator", "seconds", 10)
			r.Recorder.Eventf(cr, "Warning", "NotReady", fmt.Sprintf("Plugin serice not ready %s/%s", req.Namespace, req.Name))
			r.Condition.SetConditionPluginReadyFalse(ctx, cr)
			return ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}, nil
		}
	}

	// ingress requested
	applied = gatewayManager.IsCrApplied(ctx, cr)

	if !applied {
		if err := gatewayManager.ApplyCr(ctx, cr, r.Scheme); err != nil {
			log.Info("error ApplyCr reschedule reconcile", "error", err)
			r.Condition.SetConditionPluginReadyFalse(ctx, cr)
			return ctrl.Result{}, err
		}
	}
	r.Recorder.Eventf(cr, "Normal", "Updated", fmt.Sprintf("Updated gateway %s/%s", req.Namespace, req.Name))

	// ingress ready
	ready = gatewayManager.IsCrReady(ctx, cr)

	if !ready {
		if ready, err = gatewayManager.CheckCr(ctx, cr); err != nil {
			log.Info("error CheckCr reschedule reconcile", "error", err)
			r.Condition.SetConditionPluginReadyFalse(ctx, cr)
			return ctrl.Result{}, err
		}
		if !ready {
			log.Info("GatewayCr not ready reschedule operator", "seconds", 10)
			r.Recorder.Eventf(cr, "Warning", "NotReady", fmt.Sprintf("Gateway cr not ready %s/%s", req.Namespace, req.Name))
			r.Condition.SetConditionPluginReadyFalse(ctx, cr)
			return ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}, nil
		}
	}

	r.Recorder.Eventf(cr, "Normal", "Done", fmt.Sprintf("Plugin deployed %s/%s", req.Namespace, req.Name))
	r.Condition.SetConditionPluginReadyTrue(ctx, cr)
	return ctrl.Result{}, nil
}
