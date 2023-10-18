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

package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/docker/distribution/reference"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	envtypes "github.com/openclarity/vmclarity/e2e/testenv/types"
)

func GetImageRegistryRepositoryTag(image string) (map[string]string, error) {
	named, err := reference.ParseNormalizedNamed(image)
	if err != nil {
		return nil, fmt.Errorf("failed to parse image: %s", image)
	}

	registry := reference.Domain(named)
	repository := reference.Path(named)
	tagged, ok := named.(reference.Tagged)
	if !ok {
		return nil, fmt.Errorf("failed to get image tag from image name: %s", image)
	}
	tag := tagged.Tag()

	return map[string]string{
		"registry":   registry,
		"repository": repository,
		"tag":        tag,
	}, nil
}

func CreateK8sClient(kubeConfig string) (kubernetes.Interface, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create k8s config: %w", err)
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create k8s client: %w", err)
	}
	return clientSet, nil
}

func ListVMClarityPods(ctx context.Context, namespace string) (*corev1.PodList, error) {
	clientSet, err := GetKubernetesClientFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get kuberntes clientset from context: %w", err)
	}
	pods, err := clientSet.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error getting pods: %w", err)
	}

	return pods, nil
}

func GetVMClarityPodByName(ctx context.Context, podName, namespace string) (*corev1.Pod, error) {
	clientSet, err := GetKubernetesClientFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get kuberntes clientset from context: %w", err)
	}
	pod, err := clientSet.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pod %s: %w", podName, err)
	}

	return pod, nil
}

func GetPodLogs(ctx context.Context, pod *corev1.Pod, startTime time.Time) ([]byte, error) {
	podLogOpts := corev1.PodLogOptions{
		SinceTime: &metav1.Time{
			Time: startTime,
		},
	}
	clientSet, err := GetKubernetesClientFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get kuberntes clientset from context: %w", err)
	}
	req := clientSet.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &podLogOpts)
	podLogs, err := req.Stream(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get pod logs: %w", err)
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return nil, fmt.Errorf("error in copy information from podLogs to buf: %w", err)
	}

	return buf.Bytes(), nil
}

func GetKubernetesClientFromContext(ctx context.Context) (kubernetes.Interface, error) {
	clientSet, ok := ctx.Value(envtypes.KubernetesContextKey).(kubernetes.Interface)
	if !ok {
		return nil, fmt.Errorf(
			"context key doesn't exist: %s",
			ctx.Value(envtypes.KubernetesContextKey),
		)
	}

	return clientSet, nil
}
