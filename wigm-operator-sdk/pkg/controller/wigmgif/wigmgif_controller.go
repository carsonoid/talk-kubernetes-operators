package wigmgif

import (
	"context"

	wigmv1 "github.com/carsonoid/talk-kubernetes-operators/wigm-operator-sdk/pkg/apis/wigm/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_wigmgif")

// Add creates a new WigmGif Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileWigmGif{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("wigmgif-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource WigmGif
	err = c.Watch(&source.Kind{Type: &wigmv1.WigmGif{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to child deployments and requeue the owner WigmGif
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &wigmv1.WigmGif{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to child services and requeue the owner WigmGif
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &wigmv1.WigmGif{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to child ingress rules and requeue the owner WigmGif
	err = c.Watch(&source.Kind{Type: &extensionsv1beta1.Ingress{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &wigmv1.WigmGif{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileWigmGif{}

// ReconcileWigmGif reconciles a WigmGif object
type ReconcileWigmGif struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a WigmGif object and makes changes based on the state read
// and what is in the WigmGif.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileWigmGif) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling WigmGif")

	// Fetch the WigmGif instance
	instance := &wigmv1.WigmGif{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Sync the deployment, updating status in the passed instance
	if r, err := r.syncDeployment(instance); err != nil || r.Requeue == true {
		return r, err
	}

	// Sync the service, updating status in the passed instance
	if r, err := r.syncService(instance); err != nil || r.Requeue == true {
		return r, err
	}

	// Sync the ingress, updating status in the passed instance
	if r, err := r.syncIngress(instance); err != nil || r.Requeue == true {
		return r, err
	}

	// update the status
	if err := r.client.Status().Update(context.TODO(), instance); err != nil {
		log.Error(err, "Failed to update WigumGif status")
		return reconcile.Result{}, err
	}

	// Success! Return without requeue
	return reconcile.Result{}, nil
}
