package wigmgif

import (
	"context"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	wigmv1 "github.com/carsonoid/talk-kubernetes-operators/wigm-operator-sdk/pkg/apis/wigm/v1"
)

// Example of e2e testing with the fake client

// TestWigmGifBasicCreate tests that a valid wigmgif object results in a the expected resources
func TestWigmGifBasicCreate(t *testing.T) {
	// Define a minimal cluster which matches one of the cells above
	ingressEnabled := true
	wg := &wigmv1.WigmGif{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "test",
		},
		Spec: wigmv1.WigmGifSpec{
			Gif: wigmv1.GifProperties{
				Link: "http://testlink",
			},
			// Enable ingress
			Ingress: &wigmv1.IngressProperties{
				Enabled: &ingressEnabled,
			},
		},
	}

	// Populate the client with initial data
	objs := []runtime.Object{
		wg,
	}

	// Register operator types with the runtime scheme.
	s := scheme.Scheme
	s.AddKnownTypes(wigmv1.SchemeGroupVersion, &wigmv1.WigmGif{})
	s.AddKnownTypes(appsv1.SchemeGroupVersion, &appsv1.Deployment{})
	s.AddKnownTypes(corev1.SchemeGroupVersion, &corev1.Service{})
	s.AddKnownTypes(extensionsv1beta1.SchemeGroupVersion, &extensionsv1beta1.Ingress{})

	// Create a fake client to mock API calls.
	cl := fake.NewFakeClient(objs...)

	// Create a ReconcileVitessCluster object with the scheme and fake client.
	r := &ReconcileWigmGif{client: cl, scheme: s}

	// Mock request to simulate Reconcile() being called on an event for a
	// watched resource .
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      wg.GetName(),
			Namespace: wg.GetNamespace(),
		},
	}
	res, err := r.Reconcile(req)
	if err != nil {
		t.Errorf("Unexpected error in reconcile: %s", err)
	}

	// Check the result of reconciliation to make sure it has the desired state.
	if res.Requeue {
		t.Error("reconcile requeued request and should not have")
	}

	// Check for expected deployment
	expectedDeployment := &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: "wigm-host-test", Namespace: "test"}}
	foundDeployment := &appsv1.Deployment{}

	err = cl.Get(context.TODO(), types.NamespacedName{Name: expectedDeployment.GetName(), Namespace: expectedDeployment.GetNamespace()}, foundDeployment)
	if err != nil {
		t.Errorf("Error getting expected deployment: %s", err)
	}

	// NOTE: Here is where you could do deep checks of found vs expected

	// Check for expected service
	expectedSvc := &corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: "wigm-host-test", Namespace: "test"}}
	foundSvc := &corev1.Service{}

	err = cl.Get(context.TODO(), types.NamespacedName{Name: expectedSvc.GetName(), Namespace: expectedSvc.GetNamespace()}, foundSvc)
	if err != nil {
		t.Errorf("Error getting expected service: %s", err)
	}

	// NOTE: Here is where you could do deep checks of found vs expected

	// Check for expected ingress
	expectedIngress := &extensionsv1beta1.Ingress{ObjectMeta: metav1.ObjectMeta{Name: "wigm-host-test", Namespace: "test"}}
	foundIngress := &extensionsv1beta1.Ingress{}

	err = cl.Get(context.TODO(), types.NamespacedName{Name: expectedIngress.GetName(), Namespace: expectedIngress.GetNamespace()}, foundIngress)
	if err != nil {
		t.Errorf("Error getting expected ingress: %s", err)
	}

	// NOTE: Here is where you could do deep checks of found vs expected
}
