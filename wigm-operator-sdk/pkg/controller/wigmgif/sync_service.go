package wigmgif

import (
	"context"

	wigmv1 "github.com/carsonoid/talk-kubernetes-operators/wigm-operator-sdk/pkg/apis/wigm/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileWigmGif) syncService(instance *wigmv1.WigmGif) (reconcile.Result, error) {
	reqLogger := log.WithValues()

	// Define a new service for the instance
	service := newServiceForWG(instance)

	// Set WigmGif instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, service, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Service already exists
	foundService := &corev1.Service{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, foundService)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
		err = r.client.Create(context.TODO(), service)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Service created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Service already exists - set status
	instance.Status.Service = wigmv1.ServiceStatus{
		Created: true,
		Type:    service.Spec.Type,
	}

	reqLogger.Info("Service is in desired state", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
	return reconcile.Result{}, nil
}

func newServiceForWG(wg *wigmv1.WigmGif) *corev1.Service {
	name := wg.GetName()

	labels := map[string]string{
		"app": "wigm",
		"gif": name,
	}

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wigm-host-" + name,
			Namespace: wg.GetNamespace(),
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Port:       80,
					TargetPort: intstr.FromInt(80),
					Protocol:   corev1.ProtocolTCP,
				},
			},
		},
	}

	if wg.Spec.Service != nil && wg.Spec.Service.CreateCloudLB {
		service.Spec.Type = corev1.ServiceTypeLoadBalancer
	}

	return service
}
