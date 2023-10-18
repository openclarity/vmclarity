// Copyright © 2023 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package common

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openclarity/vmclarity/pkg/shared/utils"
)

const (
	TestDeploymentName       = "test-deployment"
	TestNamespace            = corev1.NamespaceDefault
	TestReplicaNumber  int32 = 2
)

func CreateTestDeployment(ctx context.Context) error {
	clientSet, err := GetKubernetesClientFromContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get kuberntes clientset ftom context: %w", err)
	}
	deploymentsClient := clientSet.AppsV1().Deployments(TestNamespace)
	testReplicas := TestReplicaNumber
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: TestDeploymentName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: utils.PointerTo(testReplicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"scanconfig": "test",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"scanconfig": "test",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "alpine",
							Image:   "alpine:3.18.2",
							Command: []string{"sleep", "infinity"},
						},
					},
				},
			},
		},
	}

	_, err = deploymentsClient.Create(ctx, deployment, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create test deployment: %w", err)
	}

	return nil
}

func DeleteTestDeployment(ctx context.Context) error {
	clientSet, err := GetKubernetesClientFromContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get kuberntes clientset ftom context: %w", err)
	}
	deploymentsClient := clientSet.AppsV1().Deployments(TestNamespace)
	err = deploymentsClient.Delete(ctx, TestDeploymentName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to create test deployment: %w", err)
	}

	return nil
}
