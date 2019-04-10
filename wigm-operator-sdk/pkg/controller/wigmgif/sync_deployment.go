package wigmgif

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	wigmv1 "github.com/carsonoid/talk-kubernetes-operators/wigm-operator-sdk/pkg/apis/wigm/v1"
)

func (r *ReconcileWigmGif) syncDeployment(instance *wigmv1.WigmGif) (reconcile.Result, error) {
	reqLogger := log.WithValues()

	// Define a new deployment for the instance
	deployment := newDeploymentForWG(instance)

	// Set WigmGif instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, deployment, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Deployment already exists
	foundDeployment := &appsv1.Deployment{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, foundDeployment)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
		err = r.client.Create(context.TODO(), deployment)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Deployment created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Deployment exists - set status
	instance.Status.Deployment = wigmv1.DeploymentStatus{Created: true}

	// Detect a diff in desired ENV variables and update the deployment if they are different
	if deploymentContainerEnvDiff(deployment, foundDeployment, "gifhost", "GIF_NAME") ||
		deploymentContainerEnvDiff(deployment, foundDeployment, "gifhost", "GIF_LINK") {
		reqLogger.Info("Env variables have changed. Updating Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
		err = r.client.Update(context.TODO(), deployment)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else {
		reqLogger.Info("Deployment is in desired state", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
	}

	return reconcile.Result{}, nil
}

func newDeploymentForWG(wg *wigmv1.WigmGif) *appsv1.Deployment {
	name := wg.GetName()

	giflink := wg.Spec.Gif.Link
	gifname := name
	if wg.Spec.Gif.Name != "" {
		gifname = wg.Spec.Gif.Name
	}

	labels := map[string]string{
		"app": "wigm",
		"gif": name,
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "wigm-host-" + name,
			Namespace: wg.GetNamespace(),
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: getInt32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "gifhost",
							Image: "nginx:1.15-alpine",
							Command: []string{
								"sh",
							},
							Args: []string{
								"-exc",
								startScript,
							},
							Env: []corev1.EnvVar{
								{
									Name:  "GIF_NAME",
									Value: gifname,
								},
								{
									Name:  "GIF_SOURCE_LINK",
									Value: giflink,
								},
							},
						},
					},
				},
			},
		},
	}

	return deployment
}

const startScript = `
# get curl
apk update && apk add curl

# go to data dir
cd /usr/share/nginx/html

# fetch original
curl -Lo wigm.gif "$GIF_SOURCE_LINK"

# write page using heredoc
cat > index.html <<EOF
<head>
  <title>WIGM: $GIF_NAME</title>
</head>
<body>
  <h1>WIGM: $GIF_NAME</h1>
  <img src="wigm.gif" />
</body>
EOF

# start nginx
exec nginx -g "daemon off;"
`

func getInt32Ptr(i int32) *int32 {
	return &i
}

func deploymentContainerEnvDiff(d1 *appsv1.Deployment, d2 *appsv1.Deployment, containerName string, envName string) bool {
	// different until proven the same
	different := true

	// store found values
	var v1, v2 string

	// find named container in d1
	for _, c := range d1.Spec.Template.Spec.Containers {
		if c.Name == containerName {
			v1 = getContainerEnvValue(&c, envName)
		}
	}

	// find named container in d2
	for _, c := range d2.Spec.Template.Spec.Containers {
		if c.Name == containerName {
			v2 = getContainerEnvValue(&c, envName)
		}
	}

	// Check that the values are the same
	if v1 == v2 {
		different = false
	}

	return different
}

func getContainerEnvValue(c *corev1.Container, envName string) string {
	if c.Env != nil {
		for _, envVar := range c.Env {
			if envVar.Name == envName {
				return envVar.Value
			}
		}
	}

	return ""
}
