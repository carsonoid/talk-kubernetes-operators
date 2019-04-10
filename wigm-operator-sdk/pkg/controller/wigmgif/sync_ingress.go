package wigmgif

import (
	"context"

	wigmv1 "github.com/carsonoid/talk-kubernetes-operators/wigm-operator-sdk/pkg/apis/wigm/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileWigmGif) syncIngress(instance *wigmv1.WigmGif) (reconcile.Result, error) {
	reqLogger := log.WithValues()

	// Define a new ingress for the instance
	ingress := newIngressForWG(instance)

	// If .spec.ingress.enabled explicitly set to false, cleanup
	if instance.Spec.Ingress != nil && instance.Spec.Ingress.Enabled != nil && *instance.Spec.Ingress.Enabled == false {
		// Clean up any already created ingress
		err := r.client.Delete(context.TODO(), ingress)
		if err != nil && !errors.IsNotFound(err) {
			return reconcile.Result{}, err
		}

		// Ensure status of false
		instance.Status.Ingress = wigmv1.IngressStatus{Created: false}
	} else {
		// Otherwise, sync ingress

		// Set WigmGif instance as the owner and controller
		if err := controllerutil.SetControllerReference(instance, ingress, r.scheme); err != nil {
			return reconcile.Result{}, err
		}

		// Check if this Ingress already exists
		foundIngress := &extensionsv1beta1.Ingress{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: ingress.Name, Namespace: ingress.Namespace}, foundIngress)
		if err != nil && errors.IsNotFound(err) {
			reqLogger.Info("Creating a new Ingress", "Ingress.Namespace", ingress.Namespace, "Ingress.Name", ingress.Name)
			err = r.client.Create(context.TODO(), ingress)
			if err != nil {
				return reconcile.Result{}, err
			}

			// Ingress created successfully - don't requeue
			return reconcile.Result{}, nil
		} else if err != nil {
			return reconcile.Result{}, err
		}

		// Ingress already exists - set status
		instance.Status.Ingress = wigmv1.IngressStatus{Created: true}
	}

	reqLogger.Info("Ingress is in desired state", "Ingress.Namespace", ingress.Namespace, "Ingress.Name", ingress.Name)
	return reconcile.Result{}, nil
}

func newIngressForWG(wg *wigmv1.WigmGif) *extensionsv1beta1.Ingress {
	name := wg.GetName()

	labels := map[string]string{
		"app": "wigm",
		"gif": name,
	}

	ingress := &extensionsv1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wigm-host-" + name,
			Namespace: wg.GetNamespace(),
			Labels:    labels,
		},
		Spec: extensionsv1beta1.IngressSpec{
			Rules: []extensionsv1beta1.IngressRule{
				{
					Host: name + ".wigm.carson-anderson.com",
					IngressRuleValue: extensionsv1beta1.IngressRuleValue{
						HTTP: &extensionsv1beta1.HTTPIngressRuleValue{
							Paths: []extensionsv1beta1.HTTPIngressPath{
								{
									Backend: extensionsv1beta1.IngressBackend{
										ServiceName: "wigm-host-" + name,
										ServicePort: intstr.FromInt(80),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return ingress
}
