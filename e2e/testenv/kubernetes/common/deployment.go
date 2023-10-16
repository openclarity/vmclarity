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
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/openclarity/vmclarity/pkg/shared/utils"
)

const (
	TestDeploymentName = "test-deployment"
	TestNamespace      = apiv1.NamespaceDefault
)

func CreateTestDeployment(clientSet kubernetes.Interface) error {
	deploymentsClient := clientSet.AppsV1().Deployments(TestNamespace)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: TestDeploymentName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: utils.PointerTo(int32(2)),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"scanconfig": "test",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"scanconfig": "test",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
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

	_, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create test deployment: %w", err)
	}

	return nil
}

func DeleteTestDeployment(clientSet kubernetes.Interface) error {
	deploymentsClient := clientSet.AppsV1().Deployments(TestNamespace)
	err := deploymentsClient.Delete(context.TODO(), TestDeploymentName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to create test deployment: %w", err)
	}

	return nil
}
