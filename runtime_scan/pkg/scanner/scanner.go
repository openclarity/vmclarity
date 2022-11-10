// Copyright © 2022 Cisco Systems, Inc. and its affiliates.
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

package scanner

import (
	"fmt"
	"go/scanner"
	"sync"

	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"

	_config "github.com/openclarity/vmclarity/runtime_scan/pkg/config"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/provider"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/types"
)

type Scanner struct {
	instanceIDToScanData map[string]*scanData
	providerClient       provider.Client
	scanConfig           *_config.ScanConfig
	killSignal           chan bool
	progress             types.ScanProgress
	logFields            log.Fields

	region   string
	vpcID    string
	subnetID string
	amiID    string

	sync.Mutex
}

type scanData struct {
	instance types.Instance
	//contexts                     []*imagePodContext // All the pods that contain this image hash
	scanUUID              string
	vulnerabilitiesResult vulnerabilitiesScanResult
	//cisDockerBenchmarkResult     cisDockerBenchmarkScanResult
	//shouldScanCISDockerBenchmark bool
	resultChan                   chan bool
	success                      bool
	completed                    bool
	timeout                      bool
	scanErr                      *types.ScanError
}

type vulnerabilitiesScanResult struct {
	result []string
	//layerCommands []*models.ResourceLayerCommand
	success   bool
	completed bool
	error     *scanner.Error
}

func CreateScanner(config *_config.Config, providerClient provider.Client) *Scanner {
	s := &Scanner{
		instanceIDToScanData: nil,
		providerClient:       providerClient,
		scanConfig:           nil,
		killSignal:           nil,
		progress:             types.ScanProgress{},
		logFields:            nil,
		region:               config.Region,
		vpcID:                config.VpcID,
		subnetID:             config.SubnetID,
		amiID:                config.AmiID,
		Mutex:                sync.Mutex{},
	}

	return s
}

// initScan Calculate properties of scan targets
// nolint:cyclop
func (s *Scanner) initScan() error {
	//// Get all target pods
	//for _, namespace := range s.scanConfig.Instances {
	//	podList, err := s.clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	//	if err != nil {
	//		return fmt.Errorf("failed to list pods. namespace=%s: %v", namespace, err)
	//	}
	//	podsToScan = append(podsToScan, podList.Items...)
	//	if namespace == corev1.NamespaceAll {
	//		break
	//	}
	//}

	imageIDToScanData := make(map[string]*scanData)

	// Populate the image to scanData map from all target pods
	for _, instance := range s.scanConfig.Instances {
		//if s.shouldIgnorePod(&podsToScan[i]) {
		//	continue
		//}

		// TODO: (idanf) verify if we need to read pod image pull secrets or just mount it to the scanner job
		//secrets := k8sutils.GetPodImagePullSecrets(s.clientset, pod)

		// Due to scenarios where image name in the `pod.Status.ContainerStatuses` is different
		// from image name in the `pod.Spec.Containers` we will take only image id from `pod.Status.ContainerStatuses`.
		//containerNameToImageID := make(map[string]string)
		//for _, container := range append(pod.Status.ContainerStatuses, pod.Status.InitContainerStatuses...) {
		//	containerNameToImageID[container.Name] = k8sutils.NormalizeImageID(container.ImageID)
		//}

		//containers := append(pod.Spec.Containers, pod.Spec.InitContainers...)

		//for _, container := range containers {
			//imageID, ok := containerNameToImageID[container.Name]
			//if !ok {
			//	log.Warnf("Image id is missing. pod=%v, namepspace=%v, container=%v ,image=%v",
			//		pod.GetName(), pod.GetNamespace(), container.Name, container.Image)
			//	continue
			//}
			//imageHash := k8sutils.ParseImageHash(imageID)
			//if imageHash == "" {
			//	log.WithFields(s.logFields).Warnf("Failed to get image hash - ignoring image. "+
			//		"pod=%v, namepspace=%v, image name=%v", pod.GetName(), pod.GetNamespace(), container.Image)
			//	continue
			//}
			//// Create pod context
			//podContext := &imagePodContext{
			//	containerName:   container.Name,
			//	podName:         pod.GetName(),
			//	namespace:       pod.GetNamespace(),
			//	imagePullSecret: k8sutils.GetMatchingSecretName(secrets, container.Image),
			//	imageName:       container.Image,
			//	podUID:          string(pod.GetUID()),
			//	podLabels:       labels.Set(pod.GetLabels()),
			//}
			//if data, ok := instanceIDToScanData[imageID]; !ok {
				// Image added for the first time, create scan data and append pod context
				imageIDToScanData[instance.ID] = &scanData{
					instance:                     instance,
					scanUUID:                     uuid.NewV4().String(),
					vulnerabilitiesResult:        vulnerabilitiesScanResult{},
					//shouldScanCISDockerBenchmark: s.scanConfig.ShouldScanCISDockerBenchmark,
					resultChan:                   make(chan bool),
					success:                      false,
					completed:                    false,
					timeout:                      false,
					scanErr:                      nil,
				}
			//} else {
				// Image already exist in map, just append the pod context
				//data.contexts = append(data.contexts, podContext)
			//}
		}
	//}

	s.instanceIDToScanData = imageIDToScanData
	s.progress.ImagesToScan = uint32(len(imageIDToScanData))

	log.WithFields(s.logFields).Infof("Total %d unique images to scan", s.progress.ImagesToScan)

	return nil
}

func (s *Scanner) Scan(scanConfig *_config.ScanConfig, scanDone chan struct{}) error {
	s.Lock()
	defer s.Unlock()

	s.scanConfig = scanConfig

	log.WithFields(s.logFields).Infof("Start scanning...")

	s.progress.Status = types.ScanInit

	if err := s.initScan(); err != nil {
		s.progress.SetStatus(types.ScanInitFailure)
		return fmt.Errorf("failed to initiate scan: %v", err)
	}

	if s.progress.ImagesToScan == 0 {
		log.WithFields(s.logFields).Info("Nothing to scan")
		s.progress.SetStatus(types.NothingToScan)
		nonBlockingNotification(scanDone)
		return nil
	}

	s.progress.SetStatus(types.Scanning)
	go func() {
		s.jobBatchManagement(scanDone)

		s.Lock()
		s.progress.SetStatus(types.DoneScanning)
		s.Unlock()
	}()


	return nil
}
