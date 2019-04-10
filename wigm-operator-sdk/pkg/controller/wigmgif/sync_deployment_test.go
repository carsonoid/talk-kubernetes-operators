package wigmgif

import (
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Example of unit testing with go
func TestDeploymentContainerEnvDiff(t *testing.T) {

	tests := []struct {
		vars1  []corev1.EnvVar
		vars2  []corev1.EnvVar
		expect bool
	}{
		// same
		{
			[]corev1.EnvVar{{Name: "test", Value: "testval"}},
			[]corev1.EnvVar{{Name: "test", Value: "testval"}},
			false,
		},
		// diff keys
		{
			[]corev1.EnvVar{{Name: "test", Value: "testval"}},
			[]corev1.EnvVar{{Name: "test", Value: "different"}},
			true,
		},
		// missing key1 on d1
		{
			[]corev1.EnvVar{},
			[]corev1.EnvVar{{Name: "test", Value: "different"}},
			true,
		},
		// missing key2 on d2
		{
			[]corev1.EnvVar{{Name: "test", Value: "testval"}},
			[]corev1.EnvVar{},
			true,
		},
	}

	// test env variable sets
	for _, testCase := range tests {
		d1 := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name: "d1",
			},
			Spec: appsv1.DeploymentSpec{
				Template: corev1.PodTemplateSpec{
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name: "test",
								Env:  testCase.vars1,
							},
						},
					},
				},
			},
		}

		d2 := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name: "d1",
			},
			Spec: appsv1.DeploymentSpec{
				Template: corev1.PodTemplateSpec{
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name: "test",
								Env:  testCase.vars2,
							},
						},
					},
				},
			},
		}

		// get diff result
		result := deploymentContainerEnvDiff(d1, d2, "test", "test")

		if result != testCase.expect {
			t.Errorf("Unexpected result testing env variables %#v, %#v: got %v; expected %v", testCase.vars1, testCase.vars2, result, testCase.expect)
		}
	}
}
